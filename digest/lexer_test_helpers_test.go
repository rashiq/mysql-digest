package digest

// Test helper functions for the lexer.
// These are only used in tests and are kept separate from production code.

// MustTokenText returns the text for a token, panicking if bounds are invalid.
// This is a test helper - use TokenText() in production code.
func (l *Lexer) MustTokenText(t Token) string {
	text, err := l.TokenText(t)
	if err != nil {
		panic(err)
	}
	return text
}
