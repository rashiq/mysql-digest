package digest

// Scanning utility functions for the MySQL lexer.

// scanLineComment consumes a single-line comment (# or --) until EOL.
func (l *Lexer) scanLineComment() {
	for {
		c := l.advance()
		if c == 0 || c == '\n' {
			break
		}
	}
}

// scanVersionNumber scans a version number (5 or 6 digits) at current position.
// Returns the version number and digit count without consuming characters.
func (l *Lexer) scanVersionNumber() (version int, digitCount int) {
	for i := 0; i < 6; i++ {
		ch := l.peekN(i)
		if isDigit(ch) {
			version = version*10 + int(ch-'0')
			digitCount++
		} else {
			break
		}
	}
	return version, digitCount
}

// scanDollarQuotedString scans a dollar-quoted string.
// The opening delimiter (either $$ or $tag$) has already been consumed.
func (l *Lexer) scanDollarQuotedString(tag string) Token {
	closingDelim := "$" + tag + "$"
	closingLen := len(closingDelim)

	for !l.eof() {
		if l.pos+closingLen <= len(l.input) && l.input[l.pos:l.pos+closingLen] == closingDelim {
			l.pos += closingLen
			return l.returnToken(Token{Type: DOLLAR_QUOTED_STRING_SYM, Start: l.tokStart, End: l.pos})
		}
		l.pos++
	}
	return l.returnToken(Token{Type: ABORT_SYM, Start: l.tokStart, End: l.pos})
}
