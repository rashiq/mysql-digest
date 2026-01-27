package digest

import "strings"

// stripIdentifierQuotes removes surrounding quotes from a quoted identifier.
// Handles both backtick-quoted (`ident`) and double-quoted ("ident") identifiers,
// unescaping any doubled quote characters within.
func stripIdentifierQuotes(s string) string {
	if len(s) < 2 {
		return s
	}

	switch {
	case s[0] == '`' && s[len(s)-1] == '`':
		return strings.ReplaceAll(s[1:len(s)-1], "``", "`")
	case s[0] == '"' && s[len(s)-1] == '"':
		return strings.ReplaceAll(s[1:len(s)-1], `""`, `"`)
	default:
		return s
	}
}

// escapeBackticks escapes backticks in an identifier by doubling them.
func escapeBackticks(s string) string {
	return strings.ReplaceAll(s, "`", "``")
}
