package digest

// Numeric literal handlers for the MySQL lexer.
// This module handles hex (0x...), binary (0b...), integer, and float literals.

// handleNumberIdent handles MY_LEX_NUMBER_IDENT state - numbers or identifiers starting with digit.
// This is a coordinator that dispatches to specialized handlers based on the prefix.
func (l *Lexer) handleNumberIdent() lexResult {
	c := l.input[l.tokStart]
	// Check for 0x (hex) or 0b (binary) prefix
	if c == '0' {
		nextC := l.advance()
		switch {
		case nextC == 'x' || nextC == 'X':
			return l.handleHexLiteral0x()
		case nextC == 'b' || nextC == 'B':
			return l.handleBinLiteral0b()
		default:
			l.backup() // Put back the char after '0'
		}
	}

	return l.handleDigitSequence()
}

// handleHexLiteral0x handles 0x... hex literals.
// Called after '0x' or '0X' prefix has been consumed.
func (l *Lexer) handleHexLiteral0x() lexResult {
	// Consume hex digits
	for isHexDigit(l.peek()) {
		l.skip()
	}

	// Valid hex if length >= 3 (0x + at least one digit) and not followed by ident char
	if l.tokenLen() >= 3 && !isIdentChar(l.peek()) {
		return done(Token{Type: HEX_NUM, Start: l.tokStart, End: l.pos})
	}

	// Not valid hex - treat as identifier
	l.backup()
	return cont(MY_LEX_IDENT_START)
}

// handleBinLiteral0b handles 0b... binary literals.
// Called after '0b' or '0B' prefix has been consumed.
func (l *Lexer) handleBinLiteral0b() lexResult {
	// Consume binary digits
	for {
		peek := l.peek()
		if peek != '0' && peek != '1' {
			break
		}
		l.skip()
	}

	// Valid binary if length >= 3 (0b + at least one digit) and not followed by ident char
	if l.tokenLen() >= 3 && !isIdentChar(l.peek()) {
		return done(Token{Type: BIN_NUM, Start: l.tokStart, End: l.pos})
	}

	// Not valid binary - treat as identifier
	l.backup()
	return cont(MY_LEX_IDENT_START)
}

// handleDigitSequence handles numeric sequences that may be integers, floats, or identifiers.
// Called after initial digit(s) have been processed.
func (l *Lexer) handleDigitSequence() lexResult {
	// Consume remaining digits
	for isDigit(l.peek()) {
		l.skip()
	}

	// Check what follows the digits
	nextC := l.peek()
	if !isIdentChar(nextC) {
		// Pure number, check for decimal or stay as int
		return cont(MY_LEX_INT_OR_REAL)
	}

	// Check for exponent (e/E)
	if nextC == 'e' || nextC == 'E' {
		if result, ok := l.tryParseExponent(); ok {
			return result
		}
	}

	// Number followed by identifier chars - becomes identifier
	return cont(MY_LEX_IDENT_START)
}

// tryParseExponent attempts to parse an exponent suffix (e.g., e10, e+10, e-10).
// Returns the lexResult and true if successful, or false if the exponent is invalid.
func (l *Lexer) tryParseExponent() (lexResult, bool) {
	// Save position before consuming exponent in case it's invalid
	savedPos := l.pos
	l.skip() // consume e/E

	peek := l.peek()
	if isDigit(peek) {
		// 1e10 format
		l.skip()
		for isDigit(l.peek()) {
			l.skip()
		}
		return done(Token{Type: FLOAT_NUM, Start: l.tokStart, End: l.pos}), true
	}

	if peek == '+' || peek == '-' {
		l.skip() // consume sign
		if isDigit(l.peek()) {
			l.skip()
			for isDigit(l.peek()) {
				l.skip()
			}
			return done(Token{Type: FLOAT_NUM, Start: l.tokStart, End: l.pos}), true
		}
	}

	// Not a valid exponent - restore position
	l.pos = savedPos
	return lexResult{}, false
}

// handleHexNumber handles MY_LEX_HEX_NUMBER state - X'hex' literals.
func (l *Lexer) handleHexNumber() lexResult {
	l.skip() // Skip the opening '
	for {
		c := l.advance()
		if c == '\'' {
			break
		}
		if c == 0 {
			return done(Token{
				Type:  ABORT_SYM,
				Start: l.tokStart,
				End:   l.pos,
				Err:   NewLexError(l.tokStart, ErrInvalidHexLiteral, ""),
			})
		}
		if !isHexDigit(c) {
			return done(Token{
				Type:  ABORT_SYM,
				Start: l.tokStart,
				End:   l.pos,
				Err:   NewLexError(l.tokStart, ErrInvalidHexLiteral, ""),
			})
		}
	}
	// MySQL checks: length includes x' (2) and closing ' (1), so total = hex_digits + 3
	// For valid hex, need even number of hex digits, so (length % 2) should be 1 (odd)
	length := l.tokenLen()
	if (length % 2) == 0 {
		return done(Token{
			Type:  ABORT_SYM,
			Start: l.tokStart,
			End:   l.pos,
			Err:   NewLexError(l.tokStart, ErrInvalidHexLiteral, ""),
		})
	}
	return done(Token{Type: HEX_NUM, Start: l.tokStart, End: l.pos})
}

// handleBinNumber handles MY_LEX_BIN_NUMBER state - B'bin' literals.
func (l *Lexer) handleBinNumber() lexResult {
	l.skip() // Skip the '
	for {
		c := l.advance()
		if c == '\'' {
			break
		}
		if c == 0 {
			return done(Token{
				Type:  ABORT_SYM,
				Start: l.tokStart,
				End:   l.pos,
				Err:   NewLexError(l.tokStart, ErrInvalidBinaryLiteral, ""),
			})
		}
		if c != '0' && c != '1' {
			return done(Token{
				Type:  ABORT_SYM,
				Start: l.tokStart,
				End:   l.pos,
				Err:   NewLexError(l.tokStart, ErrInvalidBinaryLiteral, ""),
			})
		}
	}
	return done(Token{Type: BIN_NUM, Start: l.tokStart, End: l.pos})
}

// handleReal handles MY_LEX_REAL state - fractional part of decimal numbers.
func (l *Lexer) handleReal() lexResult {
	for isDigit(l.peek()) {
		l.skip()
	}

	nextC := l.peek()
	if nextC == 'e' || nextC == 'E' {
		// Save position before consuming exponent in case it's invalid
		savedPos := l.pos
		l.skip() // consume e/E
		peek := l.peek()
		if peek == '+' || peek == '-' {
			l.skip() // consume sign
		}
		if !isDigit(l.peek()) {
			// Not a valid exponent - restore position and return as DECIMAL_NUM
			l.pos = savedPos
			return done(Token{Type: DECIMAL_NUM, Start: l.tokStart, End: l.pos})
		}
		for isDigit(l.peek()) {
			l.skip()
		}
		return done(Token{Type: FLOAT_NUM, Start: l.tokStart, End: l.pos})
	}

	return done(Token{Type: DECIMAL_NUM, Start: l.tokStart, End: l.pos})
}
