package digest

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

// Lexer tokenizes MySQL SQL statements.
// It replicates the behavior of MySQL's lex_one_token() from sql/sql_lex.cc.
type Lexer struct {
	input           string   // Original SQL input
	pos             int      // Current position in input
	tokStart        int      // Start position of current token
	nextState       LexState // State for next Lex() call
	sqlMode         SQLMode  // SQL mode flags
	stmtPrepareMode bool     // Whether we're in prepared statement mode
	inHintComment   bool     // Whether we're parsing inside a hint comment /*+ ... */
	lastToken       int      // Last token returned (for hint detection)
}

// NewLexer creates a new lexer for the given input.
func NewLexer(input string) *Lexer {
	return &Lexer{
		input:           input,
		pos:             0,
		tokStart:        0,
		nextState:       MY_LEX_START,
		sqlMode:         0,
		stmtPrepareMode: false,
	}
}

// SetSQLMode configures SQL mode flags.
func (l *Lexer) SetSQLMode(mode SQLMode) {
	l.sqlMode = mode
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

// restartToken resets token start to current position (after skipping whitespace).
func (l *Lexer) restartToken() {
	l.tokStart = l.pos
}

// eof returns true if we've reached end of input.
func (l *Lexer) eof() bool {
	return l.pos >= len(l.input)
}

// findKeyword looks up a keyword in the keyword map.
// If the identifier is a keyword, returns the token type.
// If function is true, also checks function keywords.
// Returns 0 if not a keyword.
func (l *Lexer) findKeyword(length int, isFunction bool) int {
	if length == 0 {
		return 0
	}
	// Get the token text and convert to uppercase for lookup
	text := l.input[l.tokStart : l.tokStart+length]
	upper := toUpper(text)

	if tok, ok := TokenKeywords[upper]; ok {
		return tok
	}
	return 0
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

// consumeComment consumes a C-style comment until closing */.
// Returns true if comment was properly closed, false if EOF reached.
func (l *Lexer) consumeComment() bool {
	for !l.eof() {
		c := l.advance()
		if c == '*' && l.peek() == '/' {
			l.skip() // Skip the '/'
			return true
		}
	}
	return false // Unclosed comment
}

// applyResult processes a lexResult and returns the token if done, otherwise returns nil.
// This centralizes the logic for handling state mutations from handler results.
func (l *Lexer) applyResult(result lexResult) *Token {
	if result.setNextLex {
		l.nextState = result.nextLexState
	}
	if result.done {
		return &result.token
	}
	return nil
}

// Lex returns the next token from the input.
// This is the main lexer loop that dispatches to state handlers.
func (l *Lexer) Lex() Token {
	// Handle hint mode - parse optimizer hint content
	if l.inHintComment {
		return l.lexHintToken()
	}

	l.startToken()
	state := l.nextState
	l.nextState = MY_LEX_START

	var c byte

	for {
		switch state {
		case MY_LEX_START:
			// Skip leading whitespace
			for getStateMap(l.peek()) == MY_LEX_SKIP {
				l.skip()
			}
			// Start of real token
			l.restartToken()
			c = l.advance()
			state = getStateMap(c)

		case MY_LEX_SKIP:
			l.skip()
			state = MY_LEX_START

		case MY_LEX_EOL:
			return Token{Type: END_OF_INPUT, Start: l.tokStart, End: l.pos}

		case MY_LEX_CHAR:
			result := l.handleChar(c)
			if tok := l.applyResult(result); tok != nil {
				return *tok
			}
			state = result.nextState

		case MY_LEX_COMMENT:
			l.scanLineComment()
			state = MY_LEX_START

		case MY_LEX_IDENT_OR_NCHAR:
			result := l.handleNChar()
			if tok := l.applyResult(result); tok != nil {
				return *tok
			}
			state = result.nextState

		case MY_LEX_IDENT_OR_HEX:
			if l.peek() == '\'' {
				state = MY_LEX_HEX_NUMBER
			} else {
				state = MY_LEX_IDENT
			}

		case MY_LEX_IDENT_OR_BIN:
			if l.peek() == '\'' {
				state = MY_LEX_BIN_NUMBER
			} else {
				state = MY_LEX_IDENT
			}

		case MY_LEX_HEX_NUMBER:
			result := l.handleHexNumber()
			if tok := l.applyResult(result); tok != nil {
				return *tok
			}
			state = result.nextState

		case MY_LEX_BIN_NUMBER:
			result := l.handleBinNumber()
			if tok := l.applyResult(result); tok != nil {
				return *tok
			}
			state = result.nextState

		case MY_LEX_IDENT_OR_DOLLAR_QUOTED_TEXT:
			result := l.handleDollarQuoted()
			if tok := l.applyResult(result); tok != nil {
				return *tok
			}
			state = result.nextState

		case MY_LEX_IDENT:
			result := l.handleIdent()
			if tok := l.applyResult(result); tok != nil {
				return *tok
			}
			state = result.nextState

		case MY_LEX_IDENT_SEP:
			result := l.handleIdentSep()
			if tok := l.applyResult(result); tok != nil {
				return *tok
			}
			state = result.nextState

		case MY_LEX_IDENT_START:
			result := l.handleIdentStart()
			if tok := l.applyResult(result); tok != nil {
				return *tok
			}
			state = result.nextState

		case MY_LEX_NUMBER_IDENT:
			result := l.handleNumberIdent(c)
			if tok := l.applyResult(result); tok != nil {
				return *tok
			}
			state = result.nextState

		case MY_LEX_INT_OR_REAL:
			if l.peek() != '.' {
				length := l.tokenLen()
				return Token{Type: l.intToken(length), Start: l.tokStart, End: l.pos}
			}
			l.skip()
			state = MY_LEX_REAL

		case MY_LEX_REAL:
			result := l.handleReal()
			if tok := l.applyResult(result); tok != nil {
				return *tok
			}
			state = result.nextState

		case MY_LEX_REAL_OR_POINT:
			if isDigit(l.peek()) {
				state = MY_LEX_REAL
			} else {
				return Token{Type: int(c), Start: l.tokStart, End: l.pos}
			}

		case MY_LEX_STRING:
			result := l.handleString(c)
			if tok := l.applyResult(result); tok != nil {
				return *tok
			}
			return result.token

		case MY_LEX_STRING_OR_DELIMITER:
			if (l.sqlMode & MODE_ANSI_QUOTES) != 0 {
				state = MY_LEX_USER_VARIABLE_DELIMITER
			} else {
				state = MY_LEX_STRING
			}

		case MY_LEX_USER_VARIABLE_DELIMITER:
			result := l.handleQuotedIdent(c)
			if tok := l.applyResult(result); tok != nil {
				return *tok
			}
			return result.token

		case MY_LEX_LONG_COMMENT:
			result := l.handleLongComment(c)
			if tok := l.applyResult(result); tok != nil {
				return *tok
			}
			state = result.nextState

		case MY_LEX_END_LONG_COMMENT:
			return Token{Type: int(c), Start: l.tokStart, End: l.pos}

		case MY_LEX_CMP_OP:
			result := l.handleCmpOp()
			if tok := l.applyResult(result); tok != nil {
				return *tok
			}
			state = result.nextState
			c = l.input[l.tokStart] // Restore c for MY_LEX_CHAR fallback

		case MY_LEX_LONG_CMP_OP:
			result := l.handleLongCmpOp()
			if tok := l.applyResult(result); tok != nil {
				return *tok
			}
			state = result.nextState
			c = l.input[l.tokStart]

		case MY_LEX_BOOL:
			result := l.handleBool(c)
			if tok := l.applyResult(result); tok != nil {
				return *tok
			}
			return result.token

		case MY_LEX_SET_VAR:
			result := l.handleSetVar(c)
			if tok := l.applyResult(result); tok != nil {
				return *tok
			}
			return result.token

		case MY_LEX_SEMICOLON:
			return Token{Type: int(c), Start: l.tokStart, End: l.pos}

		default:
			return Token{Type: int(c), Start: l.tokStart, End: l.pos}
		}
	}
}
