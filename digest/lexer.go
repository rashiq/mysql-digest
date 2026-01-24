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
func (l *Lexer) TokenText(t Token) string {
	if t.Start < 0 || t.End > len(l.input) || t.Start > t.End {
		return ""
	}
	return l.input[t.Start:t.End]
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

// scanDollarQuotedString scans a dollar-quoted string.
// The opening delimiter (either $$ or $tag$) has already been consumed.
// For anonymous: tag is empty, we look for $$
// For tagged: we look for $tag$
// Returns DOLLAR_QUOTED_STRING_SYM on success, ABORT_SYM on unterminated.
func (l *Lexer) scanDollarQuotedString(tag string) Token {
	// Build the closing delimiter
	closingDelim := "$" + tag + "$"
	closingLen := len(closingDelim)

	// Scan until we find the closing delimiter
	for !l.eof() {
		// Check if current position starts with closing delimiter
		if l.pos+closingLen <= len(l.input) {
			if l.input[l.pos:l.pos+closingLen] == closingDelim {
				// Found closing delimiter
				l.pos += closingLen
				return l.returnToken(Token{Type: DOLLAR_QUOTED_STRING_SYM, Start: l.tokStart, End: l.pos})
			}
		}
		l.pos++
	}

	// Unterminated dollar-quoted string
	return l.returnToken(Token{Type: ABORT_SYM, Start: l.tokStart, End: l.pos})
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

// Constants for int_token comparison (matching MySQL's sql_lex.cc)
const (
	longStr             = "2147483647"
	longLen             = 10
	signedLongStr       = "-2147483648"
	longlongStr         = "9223372036854775807"
	longlongLen         = 19
	signedLonglongStr   = "-9223372036854775808"
	signedLonglongLen   = 19
	unsignedLonglongStr = "18446744073709551615"
	unsignedLonglongLen = 20
)

// intToken determines the token type for an integer based on its length/value.
// Returns NUM, LONG_NUM, ULONGLONG_NUM, or DECIMAL_NUM.
// This matches MySQL's int_token() function in sql_lex.cc.
func (l *Lexer) intToken(length int) int {
	str := l.input[l.tokStart : l.tokStart+length]

	// Quick normal case - short numbers are always NUM
	if length < longLen {
		return NUM
	}

	neg := false
	offset := 0

	// Remove sign and pre-zeros
	if len(str) > 0 && str[0] == '+' {
		offset++
		length--
	} else if len(str) > 0 && str[0] == '-' {
		offset++
		length--
		neg = true
	}
	str = str[offset:]

	// Skip leading zeros
	for length > 0 && len(str) > 0 && str[0] == '0' {
		str = str[1:]
		length--
	}

	if length < longLen {
		return NUM
	}

	var cmp string
	var smaller, bigger int

	if neg {
		if length == longLen {
			cmp = signedLongStr[1:] // Skip the '-'
			smaller = NUM
			bigger = LONG_NUM
		} else if length < signedLonglongLen {
			return LONG_NUM
		} else if length > signedLonglongLen {
			return DECIMAL_NUM
		} else {
			cmp = signedLonglongStr[1:] // Skip the '-'
			smaller = LONG_NUM
			bigger = DECIMAL_NUM
		}
	} else {
		if length == longLen {
			cmp = longStr
			smaller = NUM
			bigger = LONG_NUM
		} else if length < longlongLen {
			return LONG_NUM
		} else if length > longlongLen {
			if length > unsignedLonglongLen {
				return DECIMAL_NUM
			}
			cmp = unsignedLonglongStr
			smaller = ULONGLONG_NUM
			bigger = DECIMAL_NUM
		} else {
			cmp = longlongStr
			smaller = LONG_NUM
			bigger = ULONGLONG_NUM
		}
	}

	// Compare digit by digit
	for i := 0; i < len(str) && i < len(cmp); i++ {
		if str[i] < cmp[i] {
			return smaller
		}
		if str[i] > cmp[i] {
			return bigger
		}
	}
	return smaller // Equal means it fits
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

// lexHintToken tokenizes inside an optimizer hint comment /*+ ... */
// Returns hint keywords, identifiers, numbers, operators, and TOK_HINT_COMMENT_CLOSE
func (l *Lexer) lexHintToken() Token {
	l.startToken()

	// Skip whitespace
	for {
		c := l.peek()
		if isSpace(c) {
			l.skip()
		} else {
			break
		}
	}
	l.restartToken()

	// Check for end of hint comment */
	if l.peek() == '*' && l.peekN(1) == '/' {
		l.skip() // *
		l.skip() // /
		l.inHintComment = false
		return l.returnToken(Token{Type: TOK_HINT_COMMENT_CLOSE, Start: l.tokStart, End: l.pos})
	}

	// Check for EOF (unclosed hint)
	if l.eof() {
		l.inHintComment = false
		return l.returnToken(Token{Type: ABORT_SYM, Start: l.tokStart, End: l.pos})
	}

	c := l.advance()

	// Identifier or hint keyword
	if isIdentStart(c) {
		for isIdentChar(l.peek()) {
			l.skip()
		}
		length := l.tokenLen()
		// Check if it's a hint keyword
		text := l.input[l.tokStart : l.tokStart+length]
		upper := toUpper(text)
		if tok, ok := HintKeywords[upper]; ok {
			return l.returnToken(Token{Type: tok, Start: l.tokStart, End: l.pos})
		}
		// Return as IDENT
		return l.returnToken(Token{Type: IDENT, Start: l.tokStart, End: l.pos})
	}

	// Number
	if isDigit(c) {
		for isDigit(l.peek()) {
			l.skip()
		}
		return l.returnToken(Token{Type: NUM, Start: l.tokStart, End: l.pos})
	}

	// Single-quoted string literal
	if c == '\'' {
		for {
			ch := l.peek()
			if l.eof() {
				return l.returnToken(Token{Type: ABORT_SYM, Start: l.tokStart, End: l.pos})
			}
			l.skip()
			if ch == '\'' {
				// Check for escaped quote ''
				if l.peek() == '\'' {
					l.skip()
					continue
				}
				break
			}
		}
		return l.returnToken(Token{Type: TEXT_STRING, Start: l.tokStart, End: l.pos})
	}

	// Backtick-quoted identifier
	if c == '`' {
		for {
			ch := l.peek()
			if l.eof() {
				return l.returnToken(Token{Type: ABORT_SYM, Start: l.tokStart, End: l.pos})
			}
			l.skip()
			if ch == '`' {
				// Check for escaped backtick ``
				if l.peek() == '`' {
					l.skip()
					continue
				}
				break
			}
		}
		return l.returnToken(Token{Type: IDENT, Start: l.tokStart, End: l.pos})
	}

	// Single-char tokens (parens, comma, etc.)
	// Return as their ASCII value
	return l.returnToken(Token{Type: int(c), Start: l.tokStart, End: l.pos})
}

// Lex returns the next token from the input.
// This is a stub that will be implemented state by state.
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
			// Should not normally reach here, but handle it
			l.skip()
			state = MY_LEX_START

		case MY_LEX_EOL:
			// End of input
			return Token{Type: END_OF_INPUT, Start: l.tokStart, End: l.pos}

		case MY_LEX_CHAR:
			// Unknown or single char token
			// Check for special two-char sequences with '-'
			if c == '-' && l.peek() == '-' {
				// Check for "-- " comment (-- followed by space or control char)
				nextChar := l.peekN(1)
				if isSpace(nextChar) || isCntrl(nextChar) {
					state = MY_LEX_COMMENT
					continue
				}
			}

			// Check for JSON arrow operators
			if c == '-' && l.peek() == '>' {
				l.skip() // consume '>'
				l.nextState = MY_LEX_START
				if l.peek() == '>' {
					l.skip() // consume second '>'
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
			if c == '?' && l.stmtPrepareMode && !isIdentChar(l.peek()) {
				return Token{Type: PARAM_MARKER, Start: l.tokStart, End: l.pos}
			}

			// Return the character as its ASCII value
			return Token{Type: int(c), Start: l.tokStart, End: l.pos}

		case MY_LEX_COMMENT:
			// Single-line comment (-- or #)
			// Skip until end of line
			for {
				c = l.advance()
				if c == 0 || c == '\n' {
					break
				}
			}
			// Continue lexing from start
			state = MY_LEX_START

		case MY_LEX_IDENT_OR_NCHAR:
			// Check for N'string'
			if l.peek() != '\'' {
				state = MY_LEX_IDENT
				continue
			}
			// Found N'string' - parse as NCHAR_STRING
			l.skip() // Skip the opening '
			// Find closing quote (handling escaped quotes)
			for {
				c = l.advance()
				if c == 0 {
					// Unclosed string - error
					return Token{Type: ABORT_SYM, Start: l.tokStart, End: l.pos}
				}
				if c == '\'' {
					if l.peek() == '\'' {
						l.skip() // Skip escaped quote
					} else {
						break // End of string
					}
				}
			}
			return Token{Type: NCHAR_STRING, Start: l.tokStart, End: l.pos}

		case MY_LEX_IDENT_OR_HEX:
			// Check for X'hex'
			if l.peek() == '\'' {
				state = MY_LEX_HEX_NUMBER
				continue
			}
			// Fall through to IDENT
			state = MY_LEX_IDENT
			continue

		case MY_LEX_IDENT_OR_BIN:
			// Check for B'bin'
			if l.peek() == '\'' {
				state = MY_LEX_BIN_NUMBER
				continue
			}
			// Fall through to IDENT
			state = MY_LEX_IDENT
			continue

		case MY_LEX_HEX_NUMBER:
			// X'hex' or x'hex' - skip the opening quote and consume hex digits
			l.skip() // Skip the opening '
			for {
				c = l.advance()
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
			// MySQL checks: length includes x' (2) and closing ' (1), so total = hex_digits + 3
			// For valid hex, need even number of hex digits, so (length % 2) should be 1 (odd)
			// If length is even, we have odd hex digits → ABORT_SYM
			length := l.tokenLen()
			if (length % 2) == 0 {
				return Token{Type: ABORT_SYM, Start: l.tokStart, End: l.pos}
			}
			return Token{Type: HEX_NUM, Start: l.tokStart, End: l.pos}

		case MY_LEX_BIN_NUMBER:
			// B'bin' - skip the opening quote and consume binary digits
			l.skip() // Skip the '
			for {
				c = l.advance()
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

		case MY_LEX_IDENT_OR_DOLLAR_QUOTED_TEXT:
			// Handle $ - could be:
			// 1. $$ ... $$ (anonymous dollar-quoted string)
			// 2. $tag$ ... $tag$ (tagged dollar-quoted string)
			// 3. $ident (identifier starting with $)
			// 4. $ alone (identifier)
			//
			// c is '$' (already consumed)
			if l.peek() == '$' {
				// $$...$$ anonymous dollar-quoted string
				l.skip() // consume second $
				return l.scanDollarQuotedString("")
			}

			// Check for $tag$...$tag$ (tag is identifier chars between two $)
			tagStart := l.pos
			for isIdentChar(l.peek()) && l.peek() != '$' {
				l.skip()
			}

			if l.peek() == '$' && l.pos > tagStart {
				// We have $tag$ - this is a tagged dollar-quoted string
				tag := l.input[tagStart:l.pos]
				l.skip() // consume the closing $ of the tag
				return l.scanDollarQuotedString(tag)
			}

			// Not a dollar-quoted string - reset and treat as identifier
			// Continue scanning as identifier ($ followed by ident chars)
			for isIdentChar(l.peek()) {
				l.skip()
			}

			length := l.tokenLen()
			// Return as IDENT (no keyword check for $ identifiers)
			return l.returnToken(Token{Type: IDENT, Start: l.tokStart, End: l.tokStart + length})

		case MY_LEX_IDENT:
			// Scan identifier
			// The first character (c) was already consumed in MY_LEX_START
			// Continue consuming identifier characters
			for isIdentChar(l.peek()) {
				l.skip()
			}

			length := l.tokenLen()

			// Check if followed by '.' and identifier char
			if l.peek() == '.' && isIdentChar(l.peekN(1)) {
				l.nextState = MY_LEX_IDENT_SEP
				// Still do keyword lookup for system variable scopes
				// (global, session, etc. should be recognized as keywords)
				if tokval := l.findKeyword(length, false); tokval != 0 {
					return l.returnToken(Token{Type: tokval, Start: l.tokStart, End: l.tokStart + length})
				}
			} else {
				l.backup() // Unget the non-ident char

				// Check if it's a keyword
				// The '(' check is for function keywords
				nextChar := l.peekN(1)
				if tokval := l.findKeyword(length, nextChar == '('); tokval != 0 {
					l.skip() // Re-skip the character we ungot
					l.nextState = MY_LEX_START
					return l.returnToken(Token{Type: tokval, Start: l.tokStart, End: l.tokStart + length})
				}
				l.skip() // Re-skip
			}

			// Return as IDENT
			return l.returnToken(Token{Type: IDENT, Start: l.tokStart, End: l.tokStart + length})

		case MY_LEX_IDENT_SEP:
			// Found ident and now '.'
			// Return the '.' and set next state
			c = l.advance() // Should be '.'
			if isIdentChar(l.peek()) {
				l.nextState = MY_LEX_IDENT_START
			} else {
				l.nextState = MY_LEX_START
			}
			return Token{Type: int(c), Start: l.tokStart, End: l.pos}

		case MY_LEX_IDENT_START:
			// Identifier after separator (like after '.' or in other contexts)
			// Consume the identifier
			for isIdentChar(l.peek()) {
				l.skip()
			}
			length := l.tokenLen()

			// Check if followed by another '.' and identifier char
			if l.peek() == '.' && isIdentChar(l.peekN(1)) {
				l.nextState = MY_LEX_IDENT_SEP
			}

			// After separator, we don't do keyword lookup - it's always an identifier
			return Token{Type: IDENT, Start: l.tokStart, End: l.tokStart + length}

		case MY_LEX_NUMBER_IDENT:
			// Number or identifier starting with digit
			// c contains the first digit (already consumed)

			// Check for 0x (hex) or 0b (binary) prefix
			if c == '0' {
				nextC := l.advance()
				if nextC == 'x' || nextC == 'X' {
					// Potential hex literal 0x...
					for isHexDigit(l.peek()) {
						l.skip()
					}
					// Valid hex if length >= 3 (0x + at least one digit) and not followed by ident char
					if l.tokenLen() >= 3 && !isIdentChar(l.peek()) {
						return Token{Type: HEX_NUM, Start: l.tokStart, End: l.pos}
					}
					// Not valid hex - treat as identifier
					l.backup()
					state = MY_LEX_IDENT_START
					continue
				} else if nextC == 'b' || nextC == 'B' {
					// Potential binary literal 0b...
					for {
						peek := l.peek()
						if peek != '0' && peek != '1' {
							break
						}
						l.skip()
					}
					// Valid binary if length >= 3 (0b + at least one digit) and not followed by ident char
					if l.tokenLen() >= 3 && !isIdentChar(l.peek()) {
						return Token{Type: BIN_NUM, Start: l.tokStart, End: l.pos}
					}
					// Not valid binary - treat as identifier
					l.backup()
					state = MY_LEX_IDENT_START
					continue
				}
				l.backup() // Put back the char after '0'
			}

			// Consume remaining digits
			for isDigit(l.peek()) {
				l.skip()
			}

			// Check what follows the digits
			nextC := l.peek()
			if !isIdentChar(nextC) {
				// Pure number, check for decimal or stay as int
				state = MY_LEX_INT_OR_REAL
				continue
			}

			// Check for exponent (e/E)
			if nextC == 'e' || nextC == 'E' {
				l.skip() // consume e/E
				peek := l.peek()
				if isDigit(peek) {
					// 1e10 format
					l.skip()
					for isDigit(l.peek()) {
						l.skip()
					}
					return Token{Type: FLOAT_NUM, Start: l.tokStart, End: l.pos}
				}
				if peek == '+' || peek == '-' {
					l.skip() // consume sign
					if isDigit(l.peek()) {
						l.skip()
						for isDigit(l.peek()) {
							l.skip()
						}
						return Token{Type: FLOAT_NUM, Start: l.tokStart, End: l.pos}
					}
				}
				// Not a valid float - unget and continue as identifier
				l.backup()
			}

			// Number followed by identifier chars - becomes identifier
			// Fall through to IDENT_START to consume rest
			state = MY_LEX_IDENT_START
			continue

		case MY_LEX_INT_OR_REAL:
			// Complete int or start of real (after decimal point)
			// c was the last char we read
			nextC := l.peek()
			if nextC != '.' {
				// Complete integer
				length := l.tokenLen()
				return Token{Type: l.intToken(length), Start: l.tokStart, End: l.pos}
			}
			// Has decimal point - continue to REAL
			l.skip() // consume '.'
			state = MY_LEX_REAL
			continue

		case MY_LEX_REAL:
			// Incomplete real number - consume fractional part
			for isDigit(l.peek()) {
				l.skip()
			}

			// Check for exponent
			nextC := l.peek()
			if nextC == 'e' || nextC == 'E' {
				l.skip() // consume e/E
				peek := l.peek()
				if peek == '+' || peek == '-' {
					l.skip() // consume sign
				}
				if !isDigit(l.peek()) {
					// No digit after sign - error, return as char
					state = MY_LEX_CHAR
					continue
				}
				for isDigit(l.peek()) {
					l.skip()
				}
				return Token{Type: FLOAT_NUM, Start: l.tokStart, End: l.pos}
			}

			// Decimal number without exponent
			return Token{Type: DECIMAL_NUM, Start: l.tokStart, End: l.pos}

		case MY_LEX_REAL_OR_POINT:
			// '.' - could be decimal number or just a dot
			if isDigit(l.peek()) {
				// .5 format - decimal number
				state = MY_LEX_REAL
				continue
			}
			// Just a dot
			return Token{Type: int(c), Start: l.tokStart, End: l.pos}

		case MY_LEX_STRING:
			// Single-quoted string 'text'
			// c is the opening quote (already consumed)
			sep := c
			for {
				c = l.advance()
				if c == 0 {
					// Unclosed string
					return Token{Type: ABORT_SYM, Start: l.tokStart, End: l.pos}
				}
				if c == '\\' && (l.sqlMode&MODE_NO_BACKSLASH_ESCAPES) == 0 {
					// Backslash escape - skip the next character
					if l.peek() != 0 {
						l.skip()
					}
					continue
				}
				if c == sep {
					// Check for doubled quote (escape)
					if l.peek() == sep {
						l.skip() // Skip the second quote
						continue
					}
					// End of string
					break
				}
			}
			return Token{Type: TEXT_STRING, Start: l.tokStart, End: l.pos}

		case MY_LEX_STRING_OR_DELIMITER:
			// Double-quoted string or identifier (depends on ANSI_QUOTES mode)
			// c is the opening quote (already consumed)
			if (l.sqlMode & MODE_ANSI_QUOTES) != 0 {
				// ANSI_QUOTES mode: " is an identifier delimiter
				state = MY_LEX_USER_VARIABLE_DELIMITER
				continue
			}
			// Default: " is a string delimiter (same as single quote)
			state = MY_LEX_STRING
			continue

		case MY_LEX_USER_VARIABLE_DELIMITER:
			// Backtick-quoted identifier `ident` or double-quoted ident in ANSI mode
			// c is the opening quote (already consumed)
			sep := c
			for {
				c = l.advance()
				if c == 0 {
					// Unclosed identifier
					return Token{Type: ABORT_SYM, Start: l.tokStart, End: l.pos}
				}
				if c == sep {
					// Check for doubled delimiter (escape)
					if l.peek() == sep {
						l.skip() // Skip the second delimiter
						continue
					}
					// End of identifier
					break
				}
			}
			return Token{Type: IDENT_QUOTED, Start: l.tokStart, End: l.pos}

		case MY_LEX_LONG_COMMENT:
			// Long C-style comment /* ... */ or version comment /*!50000 ... */
			// or optimizer hint /*+ ... */
			// c is '/' (already consumed)
			if l.peek() != '*' {
				// Not a comment, just a '/' character (probably division)
				return l.returnToken(Token{Type: int(c), Start: l.tokStart, End: l.pos})
			}

			// Skip the '*'
			l.skip()

			// Check for optimizer hint /*+
			if l.peek() == '+' {
				l.skip() // Skip '+'
				// Check if last token was a hintable keyword
				if TokenIsHintable(l.lastToken) {
					// Enter hint mode
					l.inHintComment = true
					return l.returnToken(Token{Type: TOK_HINT_COMMENT_OPEN, Start: l.tokStart, End: l.pos})
				}
				// Not after hintable keyword - treat as regular comment
				// Need to go back and consume the comment
				if !l.consumeComment() {
					return l.returnToken(Token{Type: ABORT_SYM, Start: l.tokStart, End: l.pos})
				}
				state = MY_LEX_START
				continue
			}

			// Check for version comment /*!
			if l.peek() == '!' {
				l.skip() // Skip '!'

				// Check for version number (5 or 6 digits)
				// Format: /*!50000 code */ or /*!32302 code */
				version := 0
				digitCount := 0
				for i := 0; i < 6; i++ {
					ch := l.peekN(i)
					if isDigit(ch) {
						version = version*10 + int(ch-'0')
						digitCount++
					} else {
						break
					}
				}

				if digitCount >= 5 {
					// Skip the version digits
					l.skipN(digitCount)

					// Check if version is <= current MySQL version (8.0.0 = 80000)
					// We'll use 80400 as a reasonable current version
					const currentVersion = 80400
					if version <= currentVersion {
						// Execute the content as code - restart lexing
						state = MY_LEX_START
						continue
					}
				} else if digitCount == 0 {
					// /*! without version - always execute
					state = MY_LEX_START
					continue
				}

				// Version is too high or invalid - skip as comment
				// Fall through to consume the comment
			}

			// Regular comment or version comment to skip - consume until */
			if !l.consumeComment() {
				return l.returnToken(Token{Type: ABORT_SYM, Start: l.tokStart, End: l.pos})
			}
			state = MY_LEX_START
			continue

		case MY_LEX_END_LONG_COMMENT:
			// '*' character - could be end of comment or just asterisk
			// In normal parsing (not inside comment), this is just '*'
			// The comment ending is handled inside consumeComment()
			return Token{Type: int(c), Start: l.tokStart, End: l.pos}

		case MY_LEX_CMP_OP:
			// Comparison operators: >, >=, =, !=
			// c is the first character (already consumed)
			// Check if next char is also a comparison operator
			nextState := getStateMap(l.peek())
			if nextState == MY_LEX_CMP_OP || nextState == MY_LEX_LONG_CMP_OP {
				l.skip()
			}
			// Look up the operator in keywords
			length := l.tokenLen()
			if tokval := l.findKeyword(length, false); tokval != 0 {
				l.nextState = MY_LEX_START // Allow signed numbers after
				return Token{Type: tokval, Start: l.tokStart, End: l.pos}
			}
			// Not found - return as single char
			state = MY_LEX_CHAR
			continue

		case MY_LEX_LONG_CMP_OP:
			// Long comparison operators: <, <=, <>, <=>, <<
			// c is '<' (already consumed)
			// Can have up to 3 characters: <=>
			nextState := getStateMap(l.peek())
			if nextState == MY_LEX_CMP_OP || nextState == MY_LEX_LONG_CMP_OP {
				l.skip()
				// Check for third char (for <=>)
				if getStateMap(l.peek()) == MY_LEX_CMP_OP {
					l.skip()
				}
			}
			// Look up the operator in keywords
			length := l.tokenLen()
			if tokval := l.findKeyword(length, false); tokval != 0 {
				l.nextState = MY_LEX_START // Allow signed numbers after
				return Token{Type: tokval, Start: l.tokStart, End: l.pos}
			}
			// Not found - return as single char
			state = MY_LEX_CHAR
			continue

		case MY_LEX_BOOL:
			// Boolean operators: && and ||
			// c is & or | (already consumed)
			// Need the same character again for &&/||
			if l.peek() != c {
				// Single & or | - return as char
				return Token{Type: int(c), Start: l.tokStart, End: l.pos}
			}
			l.skip()
			// Look up && or ||
			if tokval := l.findKeyword(2, false); tokval != 0 {
				l.nextState = MY_LEX_START // Allow signed numbers after
				return Token{Type: tokval, Start: l.tokStart, End: l.pos}
			}
			// Fallback (shouldn't happen)
			return Token{Type: int(c), Start: l.tokStart, End: l.pos}

		case MY_LEX_SET_VAR:
			// := operator or just :
			// c is ':' (already consumed)
			if l.peek() != '=' {
				// Just ':'
				return Token{Type: int(c), Start: l.tokStart, End: l.pos}
			}
			l.skip()
			return Token{Type: SET_VAR, Start: l.tokStart, End: l.pos}

		case MY_LEX_SEMICOLON:
			// Semicolon - just return as char
			return Token{Type: int(c), Start: l.tokStart, End: l.pos}

		default:
			// For now, return the character as a single-char token
			// This will be expanded in subsequent phases
			return Token{Type: int(c), Start: l.tokStart, End: l.pos}
		}
	}
}
