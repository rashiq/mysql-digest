package digest

// LexState represents the state of the lexer state machine.
// These match MySQL's my_lex_states enum in strings/sql_chars.h
type LexState int

const (
	MY_LEX_START LexState = iota
	MY_LEX_CHAR
	MY_LEX_IDENT
	MY_LEX_IDENT_SEP
	MY_LEX_IDENT_START
	MY_LEX_REAL
	MY_LEX_HEX_NUMBER
	MY_LEX_BIN_NUMBER
	MY_LEX_CMP_OP
	MY_LEX_LONG_CMP_OP
	MY_LEX_STRING
	MY_LEX_COMMENT
	MY_LEX_END
	MY_LEX_NUMBER_IDENT
	MY_LEX_INT_OR_REAL
	MY_LEX_REAL_OR_POINT
	MY_LEX_BOOL
	MY_LEX_EOL
	MY_LEX_LONG_COMMENT
	MY_LEX_END_LONG_COMMENT
	MY_LEX_SEMICOLON
	MY_LEX_SET_VAR
	MY_LEX_USER_END
	MY_LEX_HOSTNAME
	MY_LEX_SKIP
	MY_LEX_USER_VARIABLE_DELIMITER
	MY_LEX_SYSTEM_VAR
	MY_LEX_IDENT_OR_KEYWORD
	MY_LEX_IDENT_OR_HEX
	MY_LEX_IDENT_OR_BIN
	MY_LEX_IDENT_OR_NCHAR
	MY_LEX_IDENT_OR_DOLLAR_QUOTED_TEXT
	MY_LEX_STRING_OR_DELIMITER
)

// stateMap maps ASCII characters to their initial lexer states.
// This replicates MySQL's init_state_maps() from strings/sql_chars.cc
var stateMap [256]LexState

func init() {
	// Default: all characters start as MY_LEX_CHAR
	for i := 0; i < 256; i++ {
		stateMap[i] = MY_LEX_CHAR
	}

	// Alphabetic characters -> MY_LEX_IDENT
	for c := 'a'; c <= 'z'; c++ {
		stateMap[c] = MY_LEX_IDENT
	}
	for c := 'A'; c <= 'Z'; c++ {
		stateMap[c] = MY_LEX_IDENT
	}

	// High-bit bytes (0x80-0xFF) -> MY_LEX_IDENT
	// In MySQL's UTF-8 ctype table (ctype_utf8mb4), all bytes 0x80-0xFF have
	// ctype value 3 (MY_CHAR_U | MY_CHAR_L), meaning they are treated as alpha.
	// This is how MySQL handles multi-byte UTF-8 characters in identifiers.
	// See: mysql-server/strings/ctype-utf8.cc ctype_utf8mb4[] and
	//      mysql-server/strings/sql_chars.cc init_state_maps()
	for i := 0x80; i <= 0xFF; i++ {
		stateMap[i] = MY_LEX_IDENT
	}

	// Digits -> MY_LEX_NUMBER_IDENT
	for c := '0'; c <= '9'; c++ {
		stateMap[c] = MY_LEX_NUMBER_IDENT
	}

	// Whitespace -> MY_LEX_SKIP
	stateMap[' '] = MY_LEX_SKIP
	stateMap['\t'] = MY_LEX_SKIP
	stateMap['\n'] = MY_LEX_SKIP
	stateMap['\r'] = MY_LEX_SKIP
	stateMap['\v'] = MY_LEX_SKIP // vertical tab
	stateMap['\f'] = MY_LEX_SKIP // form feed

	// Special identifier characters
	stateMap['_'] = MY_LEX_IDENT

	// String delimiters
	stateMap['\''] = MY_LEX_STRING

	// Decimal point
	stateMap['.'] = MY_LEX_REAL_OR_POINT

	// Comparison operators
	stateMap['>'] = MY_LEX_CMP_OP
	stateMap['='] = MY_LEX_CMP_OP
	stateMap['!'] = MY_LEX_CMP_OP
	stateMap['<'] = MY_LEX_LONG_CMP_OP

	// Boolean operators
	stateMap['&'] = MY_LEX_BOOL
	stateMap['|'] = MY_LEX_BOOL

	// Comment starters
	stateMap['#'] = MY_LEX_COMMENT
	stateMap['/'] = MY_LEX_LONG_COMMENT
	stateMap['*'] = MY_LEX_END_LONG_COMMENT

	// Other special characters
	stateMap[';'] = MY_LEX_SEMICOLON
	stateMap[':'] = MY_LEX_SET_VAR
	stateMap['@'] = MY_LEX_USER_END
	stateMap['`'] = MY_LEX_USER_VARIABLE_DELIMITER
	stateMap['"'] = MY_LEX_STRING_OR_DELIMITER

	// End of input
	stateMap[0] = MY_LEX_EOL

	// Special handling for hex/bin/nchar string prefixes
	// These override the default MY_LEX_IDENT for these letters
	stateMap['x'] = MY_LEX_IDENT_OR_HEX
	stateMap['X'] = MY_LEX_IDENT_OR_HEX
	stateMap['b'] = MY_LEX_IDENT_OR_BIN
	stateMap['B'] = MY_LEX_IDENT_OR_BIN
	stateMap['n'] = MY_LEX_IDENT_OR_NCHAR
	stateMap['N'] = MY_LEX_IDENT_OR_NCHAR

	// Dollar for dollar-quoted strings
	stateMap['$'] = MY_LEX_IDENT_OR_DOLLAR_QUOTED_TEXT
}

// getStateMap returns the initial lexer state for a given byte.
func getStateMap(c byte) LexState {
	return stateMap[c]
}

// isIdentChar returns true if the character can be part of an identifier.
// This matches MySQL's ident_map initialization in sql_chars.cc
func isIdentChar(c byte) bool {
	state := stateMap[c]
	return state == MY_LEX_IDENT || state == MY_LEX_NUMBER_IDENT ||
		state == MY_LEX_IDENT_OR_HEX || state == MY_LEX_IDENT_OR_BIN ||
		state == MY_LEX_IDENT_OR_NCHAR || state == MY_LEX_IDENT_OR_DOLLAR_QUOTED_TEXT
}

// isIdentStart returns true if the character can start an identifier.
// Unlike isIdentChar, this excludes digits.
// High-bit bytes (0x80-0xFF) can also start identifiers for UTF-8 characters.
func isIdentStart(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' || c >= 0x80
}

// isSpace returns true if the character is a whitespace character.
func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\v' || c == '\f'
}

// isCntrl returns true if the character is a control character (0x00-0x1F).
func isCntrl(c byte) bool {
	return c < 0x20
}

// isHexDigit returns true if the character is a hex digit (0-9, a-f, A-F).
func isHexDigit(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

// isDigit returns true if the character is a decimal digit.
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}
