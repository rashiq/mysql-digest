package digest

// String and quoted identifier handler methods for the MySQL lexer.
// These handle single-quoted strings, double-quoted identifiers,
// backtick-quoted identifiers, N'string' literals, and dollar-quoted strings.

// QuoteScanMode determines how the quote scanner handles escape sequences.
type QuoteScanMode int

const (
	// QuoteModeString allows backslash escapes (unless MODE_NO_BACKSLASH_ESCAPES is set)
	QuoteModeString QuoteScanMode = iota
	// QuoteModeIdentifier does not allow backslash escapes
	QuoteModeIdentifier
)

// scanQuoted scans a quoted string or identifier.
// The opening quote has already been consumed by the caller.
// sep is the quote character (', ", or `).
// mode determines whether backslash escapes are processed.
// tokenType is the token type to return on success.
// Returns a lexResult with the appropriate token.
func (l *Lexer) scanQuoted(sep byte, mode QuoteScanMode, tokenType int) lexResult {
	allowBackslashEscape := mode == QuoteModeString && (l.sqlMode&MODE_NO_BACKSLASH_ESCAPES) == 0

	for {
		c := l.advance()
		if c == 0 {
			// Unterminated quoted literal
			return done(Token{
				Type:  ABORT_SYM,
				Start: l.tokStart,
				End:   l.pos,
				Err:   NewLexError(l.tokStart, ErrUnterminatedString, l.input),
			})
		}

		// Handle backslash escapes in string mode
		if allowBackslashEscape && c == '\\' {
			if l.peek() != 0 {
				l.skip() // Skip the escaped character
			}
			continue
		}

		// Handle quote character
		if c == sep {
			if l.peek() == sep {
				// Doubled quote = escape
				l.skip()
				continue
			}
			// End of quoted literal
			break
		}
	}

	return done(Token{Type: tokenType, Start: l.tokStart, End: l.pos})
}

// handleString handles MY_LEX_STRING state - single-quoted strings.
// Uses the unified scanQuoted with string mode for backslash escape handling.
func (l *Lexer) handleString() lexResult {
	sep := l.input[l.tokStart]
	return l.scanQuoted(sep, QuoteModeString, TEXT_STRING)
}

// handleQuotedIdent handles MY_LEX_USER_VARIABLE_DELIMITER state - backtick/double-quoted identifiers.
// Uses the unified scanQuoted with identifier mode (no backslash escapes).
func (l *Lexer) handleQuotedIdent() lexResult {
	sep := l.input[l.tokStart]
	return l.scanQuoted(sep, QuoteModeIdentifier, IDENT_QUOTED)
}

// handleNChar handles MY_LEX_IDENT_OR_NCHAR state - N'string' or identifier.
// If followed by a single quote, parses as NCHAR_STRING.
// Otherwise, falls through to identifier handling.
func (l *Lexer) handleNChar() lexResult {
	if l.peek() != '\'' {
		return cont(MY_LEX_IDENT)
	}
	// Found N'string' - parse as NCHAR_STRING
	l.skip() // Skip the opening '
	return l.scanQuoted('\'', QuoteModeIdentifier, NCHAR_STRING)
}

// handleDollarQuoted handles MY_LEX_IDENT_OR_DOLLAR_QUOTED_TEXT state.
// Handles $$...$$ anonymous and $tag$...$tag$ tagged dollar-quoted strings.
func (l *Lexer) handleDollarQuoted() lexResult {
	// c is '$' (already consumed)
	if l.peek() == '$' {
		// $$...$$ anonymous dollar-quoted string
		l.skip() // consume second $
		return done(l.scanDollarQuotedString(""))
	}

	// Check for $tag$...$tag$ (tag is identifier chars between two $)
	tagStart := l.pos
	for isIdentChar(l.peek()) && l.peek() != '$' {
		l.skip()
	}

	if l.peek() == '$' && l.pos > tagStart {
		// We have $tag$ - this is a tagged dollar-quoted string
		tag := l.input[tagStart:l.pos]
		l.skip() // consume the closing $ of the tag
		return done(l.scanDollarQuotedString(tag))
	}

	// Not a dollar-quoted string - treat as identifier
	for isIdentChar(l.peek()) {
		l.skip()
	}

	length := l.tokenLen()
	return done(l.returnToken(Token{Type: IDENT, Start: l.tokStart, End: l.tokStart + length}))
}
