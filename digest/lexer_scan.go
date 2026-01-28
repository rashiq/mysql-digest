package digest

// Reusable scanning functions for the MySQL lexer.
// These handle common scanning patterns that are used across multiple lexer states.

// scanIdentifier consumes identifier characters from the current position.
// The first character has already been consumed. Returns the total length.
func (l *Lexer) scanIdentifier() int {
	for isIdentChar(l.peek()) {
		l.skip()
	}
	return l.tokenLen()
}

// scanQuotedString scans a quoted string starting after the opening quote.
// sep is the quote character (', ", or `).
// Returns true if properly closed, false on EOF.
func (l *Lexer) scanQuotedString(sep byte) bool {
	for {
		c := l.advance()
		if c == 0 {
			return false // Unclosed string
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
			return true // End of string
		}
	}
}

// scanQuotedIdentifier scans a quoted identifier (backtick or double-quote in ANSI mode).
// sep is the quote character (` or ").
// Returns true if properly closed, false on EOF.
func (l *Lexer) scanQuotedIdentifier(sep byte) bool {
	for {
		c := l.advance()
		if c == 0 {
			return false // Unclosed identifier
		}
		if c == sep {
			// Check for doubled delimiter (escape)
			if l.peek() == sep {
				l.skip() // Skip the second delimiter
				continue
			}
			return true // End of identifier
		}
	}
}

// scanHexDigits consumes hex digits from the current position.
// Returns true if at least one hex digit was consumed.
func (l *Lexer) scanHexDigits() bool {
	if !isHexDigit(l.peek()) {
		return false
	}
	for isHexDigit(l.peek()) {
		l.skip()
	}
	return true
}

// scanBinaryDigits consumes binary digits (0 or 1) from the current position.
// Returns true if at least one binary digit was consumed.
func (l *Lexer) scanBinaryDigits() bool {
	peek := l.peek()
	if peek != '0' && peek != '1' {
		return false
	}
	for {
		peek = l.peek()
		if peek != '0' && peek != '1' {
			break
		}
		l.skip()
	}
	return true
}

// scanDigits consumes decimal digits from the current position.
// Returns true if at least one digit was consumed.
func (l *Lexer) scanDigits() bool {
	if !isDigit(l.peek()) {
		return false
	}
	for isDigit(l.peek()) {
		l.skip()
	}
	return true
}

// scanExponent consumes an exponent part (e/E followed by optional sign and digits).
// Returns true if a valid exponent was consumed.
// If the exponent is invalid (e.g., "e+" without digits), restores position exactly.
func (l *Lexer) scanExponent() bool {
	peek := l.peek()
	if peek != 'e' && peek != 'E' {
		return false
	}

	// Save position before consuming anything so we can restore on failure
	savedPos := l.pos

	l.skip() // consume e/E

	peek = l.peek()
	if peek == '+' || peek == '-' {
		l.skip() // consume sign
	}

	if !isDigit(l.peek()) {
		// Invalid exponent - restore to exact saved position
		l.pos = savedPos
		return false
	}

	for isDigit(l.peek()) {
		l.skip()
	}
	return true
}

// scanLineComment consumes a single-line comment (# or --).
// Consumes until end of line or EOF.
func (l *Lexer) scanLineComment() {
	for {
		c := l.advance()
		if c == 0 || c == '\n' {
			break
		}
	}
}

// scanVersionNumber scans a version number (5 or 6 digits) at the current position.
// Returns the version number and digit count without consuming any characters.
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
