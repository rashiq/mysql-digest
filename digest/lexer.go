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
	Type  int // Token type ID (from tokens.go, or ASCII value for single chars)
	Start int // Start position in input (inclusive)
	End   int // End position in input (exclusive)
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
func (l *Lexer) TokenText(t Token) string {
	if t.Start < 0 || t.End > len(l.input) || t.Start > t.End {
		return ""
	}
	return l.input[t.Start:t.End]
}

// Input returns the original input string.
func (l *Lexer) Input() string {
	return l.input
}

// ---- Internal helper methods (matching MySQL's Lex_input_stream) ----

// yyPeek returns the current character without advancing.
func (l *Lexer) yyPeek() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

// yyPeekn returns the character at offset n from current position.
func (l *Lexer) yyPeekn(n int) byte {
	if l.pos+n >= len(l.input) {
		return 0
	}
	return l.input[l.pos+n]
}

// yyGet returns the current character and advances position.
func (l *Lexer) yyGet() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	c := l.input[l.pos]
	l.pos++
	return c
}

// yySkip advances position by one.
func (l *Lexer) yySkip() {
	if l.pos < len(l.input) {
		l.pos++
	}
}

// yySkipn advances position by n.
func (l *Lexer) yySkipn(n int) {
	l.pos += n
	if l.pos > len(l.input) {
		l.pos = len(l.input)
	}
}

// yyUnget moves position back by one.
func (l *Lexer) yyUnget() {
	if l.pos > 0 {
		l.pos--
	}
}

// yyLength returns the length of the current token.
func (l *Lexer) yyLength() int {
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

// toUpper converts a string to uppercase (ASCII only)
func toUpper(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			b[i] = c - 32
		} else {
			b[i] = c
		}
	}
	return string(b)
}

// Lex returns the next token from the input.
// This is a stub that will be implemented state by state.
func (l *Lexer) Lex() Token {
	l.startToken()
	state := l.nextState
	l.nextState = MY_LEX_START

	var c byte

	for {
		switch state {
		case MY_LEX_START:
			// Skip leading whitespace
			for getStateMap(l.yyPeek()) == MY_LEX_SKIP {
				l.yySkip()
			}

			// Start of real token
			l.restartToken()
			c = l.yyGet()
			state = getStateMap(c)

		case MY_LEX_SKIP:
			// Should not normally reach here, but handle it
			l.yySkip()
			state = MY_LEX_START

		case MY_LEX_EOL:
			// End of input
			return Token{Type: END_OF_INPUT, Start: l.tokStart, End: l.pos}

		case MY_LEX_CHAR:
			// Unknown or single char token
			// Check for special two-char sequences with '-'
			if c == '-' && l.yyPeek() == '-' {
				// Check for "-- " comment (-- followed by space or control char)
				nextChar := l.yyPeekn(1)
				if isSpace(nextChar) || isCntrl(nextChar) {
					state = MY_LEX_COMMENT
					continue
				}
			}

			// Check for JSON arrow operators
			if c == '-' && l.yyPeek() == '>' {
				l.yySkip() // consume '>'
				l.nextState = MY_LEX_START
				if l.yyPeek() == '>' {
					l.yySkip() // consume second '>'
					return Token{Type: JSON_UNQUOTED_SEPARATOR_SYM, Start: l.tokStart, End: l.pos}
				}
				return Token{Type: JSON_SEPARATOR_SYM, Start: l.tokStart, End: l.pos}
			}

			// Close paren does NOT allow signed numbers after
			// (other chars allow signed numbers by setting next_state = MY_LEX_START)
			if c != ')' {
				l.nextState = MY_LEX_START
			}

			// Check for placeholder '?'
			if c == '?' && l.stmtPrepareMode && !isIdentChar(l.yyPeek()) {
				return Token{Type: PARAM_MARKER, Start: l.tokStart, End: l.pos}
			}

			// Return the character as its ASCII value
			return Token{Type: int(c), Start: l.tokStart, End: l.pos}

		case MY_LEX_COMMENT:
			// Single-line comment (-- or #)
			// Skip until end of line
			for {
				c = l.yyGet()
				if c == 0 || c == '\n' {
					break
				}
			}
			// Continue lexing from start
			state = MY_LEX_START

		case MY_LEX_IDENT_OR_NCHAR:
			// Check for N'string'
			if l.yyPeek() != '\'' {
				state = MY_LEX_IDENT
				continue
			}
			// Found N'string' - parse as NCHAR_STRING
			l.yySkip() // Skip the opening '
			// Find closing quote (handling escaped quotes)
			for {
				c = l.yyGet()
				if c == 0 {
					// Unclosed string - error
					return Token{Type: ABORT_SYM, Start: l.tokStart, End: l.pos}
				}
				if c == '\'' {
					if l.yyPeek() == '\'' {
						l.yySkip() // Skip escaped quote
					} else {
						break // End of string
					}
				}
			}
			return Token{Type: NCHAR_STRING, Start: l.tokStart, End: l.pos}

		case MY_LEX_IDENT_OR_HEX:
			// Check for X'hex'
			if l.yyPeek() == '\'' {
				state = MY_LEX_HEX_NUMBER
				continue
			}
			// Fall through to IDENT
			state = MY_LEX_IDENT
			continue

		case MY_LEX_IDENT_OR_BIN:
			// Check for B'bin'
			if l.yyPeek() == '\'' {
				state = MY_LEX_BIN_NUMBER
				continue
			}
			// Fall through to IDENT
			state = MY_LEX_IDENT
			continue

		case MY_LEX_HEX_NUMBER:
			// X'hex' - skip the opening quote and consume hex digits
			l.yySkip() // Skip the '
			for {
				c = l.yyGet()
				if c == '\'' {
					break
				}
				if c == 0 {
					return Token{Type: ABORT_SYM, Start: l.tokStart, End: l.pos}
				}
				// Validate hex digit
				if !isHexDigit(c) {
					return Token{Type: ABORT_SYM, Start: l.tokStart, End: l.pos}
				}
			}
			return Token{Type: HEX_NUM, Start: l.tokStart, End: l.pos}

		case MY_LEX_BIN_NUMBER:
			// B'bin' - skip the opening quote and consume binary digits
			l.yySkip() // Skip the '
			for {
				c = l.yyGet()
				if c == '\'' {
					break
				}
				if c == 0 {
					return Token{Type: ABORT_SYM, Start: l.tokStart, End: l.pos}
				}
				// Validate binary digit
				if c != '0' && c != '1' {
					return Token{Type: ABORT_SYM, Start: l.tokStart, End: l.pos}
				}
			}
			return Token{Type: BIN_NUM, Start: l.tokStart, End: l.pos}

		case MY_LEX_IDENT, MY_LEX_IDENT_OR_DOLLAR_QUOTED_TEXT:
			// Scan identifier
			// The first character (c) was already consumed in MY_LEX_START
			// Continue consuming identifier characters
			for isIdentChar(l.yyPeek()) {
				l.yySkip()
			}

			length := l.yyLength()

			// Check if followed by '.' and identifier char
			if l.yyPeek() == '.' && isIdentChar(l.yyPeekn(1)) {
				l.nextState = MY_LEX_IDENT_SEP
			} else {
				l.yyUnget() // Unget the non-ident char

				// Check if it's a keyword
				// The '(' check is for function keywords
				nextChar := l.yyPeekn(1)
				if tokval := l.findKeyword(length, nextChar == '('); tokval != 0 {
					l.yySkip() // Re-skip the character we ungot
					l.nextState = MY_LEX_START
					return Token{Type: tokval, Start: l.tokStart, End: l.tokStart + length}
				}
				l.yySkip() // Re-skip
			}

			// Return as IDENT
			return Token{Type: IDENT, Start: l.tokStart, End: l.tokStart + length}

		case MY_LEX_IDENT_SEP:
			// Found ident and now '.'
			// Return the '.' and set next state
			c = l.yyGet() // Should be '.'
			if isIdentChar(l.yyPeek()) {
				l.nextState = MY_LEX_IDENT_START
			} else {
				l.nextState = MY_LEX_START
			}
			return Token{Type: int(c), Start: l.tokStart, End: l.pos}

		case MY_LEX_IDENT_START:
			// Identifier after separator (like after '.' or in other contexts)
			// Consume the identifier
			for isIdentChar(l.yyPeek()) {
				l.yySkip()
			}
			length := l.yyLength()

			// Check if followed by another '.' and identifier char
			if l.yyPeek() == '.' && isIdentChar(l.yyPeekn(1)) {
				l.nextState = MY_LEX_IDENT_SEP
			}

			// After separator, we don't do keyword lookup - it's always an identifier
			return Token{Type: IDENT, Start: l.tokStart, End: l.tokStart + length}

		default:
			// For now, return the character as a single-char token
			// This will be expanded in subsequent phases
			return Token{Type: int(c), Start: l.tokStart, End: l.pos}
		}
	}
}
