package digest

// Optimizer hint comment lexer for MySQL /*+ ... */ syntax.
// This module handles tokenization inside hint comments, which have their
// own mini-grammar with hint keywords, identifiers, strings, and operators.

// lexHintToken tokenizes inside an optimizer hint comment /*+ ... */
// Returns hint keywords, identifiers, numbers, operators, and TOK_HINT_COMMENT_CLOSE
func (l *Lexer) lexHintToken() Token {
	l.startToken()

	// Skip whitespace
	l.skipHintWhitespace()
	l.restartToken()

	// Check for end of hint comment */
	if l.peek() == '*' && l.peekN(1) == '/' {
		return l.lexHintClose()
	}

	// Check for EOF (unclosed hint)
	if l.eof() {
		return l.lexHintEOF()
	}

	c := l.advance()

	// Dispatch based on first character
	switch {
	case isIdentStart(c):
		return l.lexHintIdentOrKeyword()
	case isDigit(c):
		return l.lexHintNumber()
	case c == '\'':
		return l.lexHintString()
	case c == '`':
		return l.lexHintQuotedIdent()
	default:
		return l.lexHintChar(c)
	}
}

// skipHintWhitespace skips whitespace characters inside a hint comment.
func (l *Lexer) skipHintWhitespace() {
	for isSpace(l.peek()) {
		l.skip()
	}
}

// lexHintClose handles the closing */ of a hint comment.
func (l *Lexer) lexHintClose() Token {
	l.skip() // *
	l.skip() // /
	l.inHintComment = false
	return l.returnToken(Token{Type: TOK_HINT_COMMENT_CLOSE, Start: l.tokStart, End: l.pos})
}

// lexHintEOF handles EOF inside an unclosed hint comment.
func (l *Lexer) lexHintEOF() Token {
	l.inHintComment = false
	return l.returnToken(Token{
		Type:  ABORT_SYM,
		Start: l.tokStart,
		End:   l.pos,
		Err:   NewLexError(l.tokStart, ErrUnterminatedHint, ""),
	})
}

// lexHintIdentOrKeyword lexes an identifier or hint keyword.
func (l *Lexer) lexHintIdentOrKeyword() Token {
	for isIdentChar(l.peek()) {
		l.skip()
	}
	length := l.tokenLen()

	// Check if it's a hint keyword using injected resolver
	text := l.input[l.tokStart : l.tokStart+length]
	if tok, ok := l.keywordResolver.ResolveHint(text); ok {
		return l.returnToken(Token{Type: tok, Start: l.tokStart, End: l.pos})
	}

	// Return as IDENT
	return l.returnToken(Token{Type: IDENT, Start: l.tokStart, End: l.pos})
}

// lexHintNumber lexes a numeric literal inside a hint.
func (l *Lexer) lexHintNumber() Token {
	for isDigit(l.peek()) {
		l.skip()
	}
	return l.returnToken(Token{Type: NUM, Start: l.tokStart, End: l.pos})
}

// lexHintString lexes a single-quoted string inside a hint.
func (l *Lexer) lexHintString() Token {
	for {
		ch := l.peek()
		if l.eof() {
			return l.returnToken(Token{
				Type:  ABORT_SYM,
				Start: l.tokStart,
				End:   l.pos,
				Err:   NewLexError(l.tokStart, ErrUnterminatedString, ""),
			})
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

// lexHintQuotedIdent lexes a backtick-quoted identifier inside a hint.
func (l *Lexer) lexHintQuotedIdent() Token {
	for {
		ch := l.peek()
		if l.eof() {
			return l.returnToken(Token{
				Type:  ABORT_SYM,
				Start: l.tokStart,
				End:   l.pos,
				Err:   NewLexError(l.tokStart, ErrUnterminatedIdent, ""),
			})
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

// lexHintChar returns a single character token inside a hint.
func (l *Lexer) lexHintChar(c byte) Token {
	return l.returnToken(Token{Type: int(c), Start: l.tokStart, End: l.pos})
}
