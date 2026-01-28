package digest

// Comment handler methods for the MySQL lexer.
// These handle C-style comments (/* */), optimizer hints (/*+ */),
// and version comments (/*! */).

// handleLongComment handles MY_LEX_LONG_COMMENT state - C-style comments, version comments, and hints.
// This is the coordinator function that delegates to specific handlers.
func (l *Lexer) handleLongComment() lexResult {
	c := l.input[l.tokStart]
	if l.peek() != '*' {
		// Not a comment, just a '/' character (division operator)
		return l.handleDivisionOp(c)
	}

	// Skip the '*'
	l.skip()

	// Check for optimizer hint /*+
	if l.peek() == '+' {
		return l.handleOptimizerHint()
	}

	// Check for version comment /*!
	if l.peek() == '!' {
		return l.handleVersionComment()
	}

	// Regular block comment /* ... */
	return l.handleBlockComment()
}

// handleDivisionOp handles the case where '/' is not followed by '*'.
// This is the division operator, not a comment.
func (l *Lexer) handleDivisionOp(c byte) lexResult {
	return done(l.returnToken(Token{Type: int(c), Start: l.tokStart, End: l.pos}))
}

// handleOptimizerHint handles /*+ optimizer hints.
// Returns a TOK_HINT_COMMENT_OPEN token if after a hintable keyword,
// otherwise consumes as a regular comment.
func (l *Lexer) handleOptimizerHint() lexResult {
	l.skip() // Skip '+'

	// Check if last token was a hintable keyword (SELECT, INSERT, etc.)
	if TokenIsHintable(l.lastToken) {
		// Enter hint mode - subsequent Lex() calls will use lexHintToken()
		l.inHintComment = true
		return done(l.returnToken(Token{Type: TOK_HINT_COMMENT_OPEN, Start: l.tokStart, End: l.pos}))
	}

	// Not after hintable keyword - treat as regular comment
	return l.consumeBlockComment()
}

// handleVersionComment handles /*! version comments.
// Executes content if version is <= current MySQL version or no version specified.
func (l *Lexer) handleVersionComment() lexResult {
	l.skip() // Skip '!'

	// Check for version number (5 or 6 digits)
	version, digitCount := l.scanVersionNumber()

	if digitCount >= 5 {
		// Skip the version digits
		l.skipN(digitCount)

		// Check if version is <= configured MySQL version
		if version <= l.mysqlVersion {
			// Execute the content as code - restart lexing
			return cont(MY_LEX_START)
		}
		// Version is too high - skip as comment
		return l.consumeBlockComment()
	}

	if digitCount == 0 {
		// /*! without version - always execute
		return cont(MY_LEX_START)
	}

	// Invalid version format (1-4 digits) - skip as comment
	return l.consumeBlockComment()
}

// handleBlockComment handles regular block comments /* ... */.
func (l *Lexer) handleBlockComment() lexResult {
	return l.consumeBlockComment()
}

// consumeBlockComment consumes a block comment and returns the appropriate result.
// Returns ABORT_SYM if comment is unclosed, otherwise continues to MY_LEX_START.
func (l *Lexer) consumeBlockComment() lexResult {
	if !l.consumeComment() {
		return done(Token{
			Type:  ABORT_SYM,
			Start: l.tokStart,
			End:   l.pos,
			Err:   NewLexError(l.tokStart, ErrUnterminatedComment, l.input),
		})
	}
	return cont(MY_LEX_START)
}
