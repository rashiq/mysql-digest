package digest

// Character classification utilities for the MySQL lexer.
// These match MySQL's character type tables in strings/sql_chars.cc

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
