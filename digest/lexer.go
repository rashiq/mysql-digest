package digest

import "fmt"

// ---- Lexer Errors ----

// LexError represents an error encountered during lexical analysis.
type LexError struct {
	Position int
	Message  string
	Input    string
}

func (e *LexError) Error() string {
	if e.Input != "" {
		return fmt.Sprintf("lex error at position %d: %s (near %q)", e.Position, e.Message, e.Input)
	}
	return fmt.Sprintf("lex error at position %d: %s", e.Position, e.Message)
}

// Error message constants
const (
	ErrUnterminatedString   = "unterminated string literal"
	ErrUnterminatedComment  = "unterminated block comment"
	ErrUnterminatedIdent    = "unterminated quoted identifier"
	ErrUnterminatedHint     = "unterminated optimizer hint"
	ErrInvalidHexLiteral    = "invalid hex literal"
	ErrInvalidBinaryLiteral = "invalid binary literal"
	ErrInvalidTokenBounds   = "invalid token bounds"
	ErrUnterminatedDollar   = "unterminated dollar-quoted string"
)

func NewLexError(position int, message string, input string) *LexError {
	return &LexError{Position: position, Message: message, Input: input}
}

// ---- SQL Mode ----

// SQLMode flags that affect lexer behavior
type SQLMode uint64

const (
	// MODE_NO_BACKSLASH_ESCAPES disables backslash as escape character in strings
	MODE_NO_BACKSLASH_ESCAPES SQLMode = 1 << 0
	// MODE_ANSI_QUOTES treats " as identifier delimiter instead of string delimiter
	MODE_ANSI_QUOTES SQLMode = 1 << 1
)

// Token represents a lexed token with position information.
// The actual text can be retrieved via Lexer.TokenText(t).
type Token struct {
	Type  int   // Token type ID (from tokens.go, or ASCII value for single chars)
	Start int   // Start position in input (inclusive)
	End   int   // End position in input (exclusive)
	Err   error // Non-nil if token represents an error state
}

// IsError returns true if this token represents an error condition.
// A token is an error if it has type ABORT_SYM or has a non-nil Err field.
func (t Token) IsError() bool {
	return t.Type == ABORT_SYM || t.Err != nil
}

type Lexer struct {
	input            string
	pos              int
	tokStart         int
	nextState        LexState
	sqlMode          SQLMode
	stmtPrepareMode  bool
	inHintComment    bool
	inVersionComment bool
	lastToken        int
	digestVersion    MySQLVersion
	tokenConfig      *TokenConfig
}

var mysqlVersionMap = map[MySQLVersion]int{
	MySQL84: 80400, // MySQL 8.4.0
	MySQL80: 80000, // MySQL 8.0.0
	MySQL57: 50700, // MySQL 5.7.0
}

func NewLexer(input string) *Lexer {
	return &Lexer{
		input:         input,
		nextState:     MY_LEX_START,
		digestVersion: MySQL84,
		tokenConfig:   GetTokenConfig(MySQL84),
	}
}

// SetSQLMode configures SQL mode flags.
func (l *Lexer) SetSQLMode(mode SQLMode) {
	l.sqlMode = mode
}

func (l *Lexer) SetDigestVersion(version MySQLVersion) {
	l.digestVersion = version
	l.tokenConfig = GetTokenConfig(version)
}

// mysqlVersionInt returns the integer MySQL version for version comments.
func (l *Lexer) mysqlVersionInt() int {
	return mysqlVersionMap[l.digestVersion]
}

// SetPrepareMode sets whether the lexer is in prepared statement mode.
// In prepare mode, '?' is returned as PARAM_MARKER when not followed by identifier chars.
func (l *Lexer) SetPrepareMode(enabled bool) {
	l.stmtPrepareMode = enabled
}

// TokenText returns the text for a token (slice of original input).
// Returns an error if the token bounds are invalid.
func (l *Lexer) TokenText(t Token) (string, error) {
	if t.Start < 0 || t.End > len(l.input) || t.Start > t.End {
		return "", NewLexError(t.Start, ErrInvalidTokenBounds, "")
	}
	return l.input[t.Start:t.End], nil
}

// ---- Internal helper methods ----

// peek returns the current character without advancing.
func (l *Lexer) peek() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

// peekN returns the character at offset n from current position.
func (l *Lexer) peekN(n int) byte {
	if l.pos+n >= len(l.input) {
		return 0
	}
	return l.input[l.pos+n]
}

// advance returns the current character and advances position.
func (l *Lexer) advance() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	c := l.input[l.pos]
	l.pos++
	return c
}

// skip advances position by one.
func (l *Lexer) skip() {
	if l.pos < len(l.input) {
		l.pos++
	}
}

// skipN advances position by n.
func (l *Lexer) skipN(n int) {
	l.pos += n
	if l.pos > len(l.input) {
		l.pos = len(l.input)
	}
}

// backup moves position back by one.
func (l *Lexer) backup() {
	if l.pos > 0 {
		l.pos--
	}
}

// tokenLen returns the length of the current token.
func (l *Lexer) tokenLen() int {
	return l.pos - l.tokStart
}

// startToken marks the beginning of a new token.
func (l *Lexer) startToken() {
	l.tokStart = l.pos
}

// eof returns true if we've reached end of input.
func (l *Lexer) eof() bool {
	return l.pos >= len(l.input)
}

func (l *Lexer) findKeyword(length int) int {
	if length == 0 {
		return 0
	}
	text := l.input[l.tokStart : l.tokStart+length]
	upper := toUpper(text)
	return l.tokenConfig.LookupKeyword(upper)
}

// returnToken wraps token return to track lastToken for hint detection
func (l *Lexer) returnToken(t Token) Token {
	l.lastToken = t.Type
	return t
}

// intToken determines the token type for an integer based on its length/value.
// Returns NUM, LONG_NUM, ULONGLONG_NUM, or DECIMAL_NUM.
// This matches MySQL's int_token() function in sql_lex.cc.
func (l *Lexer) intToken(length int) int {
	str := l.input[l.tokStart : l.tokStart+length]
	return ClassifyInteger(str)
}

// Lex returns the next token from the input.
func (l *Lexer) Lex() Token {
	if l.inHintComment {
		return l.lexHintToken()
	}

	l.startToken()
	state := l.nextState
	l.nextState = MY_LEX_START

	for {
		result, handled := l.dispatchState(state)
		// Fallback for unregistered states. Shouldn't happen.
		if !handled {
			c := byte(0)
			if l.tokStart < len(l.input) {
				c = l.input[l.tokStart]
			}
			return Token{Type: int(c), Start: l.tokStart, End: l.pos}
		}

		if result.setNextLex {
			l.nextState = result.nextLexState
		}
		if result.done {
			return result.token
		}
		state = result.nextState
	}
}
