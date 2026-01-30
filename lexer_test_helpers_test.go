package digest

func (l *Lexer) MustTokenText(t Token) string {
	text, err := l.TokenText(t)
	if err != nil {
		panic(err)
	}
	return text
}
