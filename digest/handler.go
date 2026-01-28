package digest

// tokenHandler processes lexer tokens and coordinates storage and reduction.
// It classifies each token and routes it to the appropriate handling logic.
type tokenHandler struct {
	lexer   *Lexer
	store   *tokenStore
	reducer *reducer
}

// newTokenHandler creates a new token handler.
func newTokenHandler(lexer *Lexer, store *tokenStore, reducer *reducer) *tokenHandler {
	return &tokenHandler{
		lexer:   lexer,
		store:   store,
		reducer: reducer,
	}
}

// processAll reads all tokens from the lexer and normalizes them.
func (h *tokenHandler) processAll() {
	for {
		tok := h.lexer.Lex()

		if tok.Type == END_OF_INPUT {
			h.store.removeTrailingSemicolon()
			break
		}
		if tok.Type == ABORT_SYM {
			break
		}

		h.handleToken(tok)
	}
}

// handleToken processes a single token based on its type.
func (h *tokenHandler) handleToken(tok Token) {
	switch {
	case isNumericLiteral(tok.Type):
		h.handleNumericLiteral()

	case isStringLiteral(tok.Type):
		h.handleLiteral()

	case tok.Type == NULL_SYM:
		h.handleNull()

	case tok.Type == ')':
		h.handleCloseParen()

	case tok.Type == IDENT || tok.Type == IDENT_QUOTED:
		h.handleIdentifier(tok)

	default:
		h.store.push(tok.Type)
		h.reducer.reduceAll()
	}
}

// handleNumericLiteral processes numeric literal tokens.
// Absorbs any preceding unary +/- signs before normalizing.
func (h *tokenHandler) handleNumericLiteral() {
	h.reducer.reduceUnarySign()
	h.store.push(TOK_GENERIC_VALUE)
	h.reducer.reduceAfterValue()
}

// handleLiteral processes string and parameter marker tokens.
func (h *tokenHandler) handleLiteral() {
	h.store.push(TOK_GENERIC_VALUE)
	h.reducer.reduceAfterValue()
}

// handleNull processes NULL tokens.
// NULL is kept as a keyword after IS/IS NOT, otherwise normalized to a value.
func (h *tokenHandler) handleNull() {
	if h.isNullKeywordContext() {
		h.store.push(NULL_SYM)
	} else {
		h.store.push(TOK_GENERIC_VALUE)
		h.reducer.reduceAfterValue()
	}
}

// handleCloseParen stores ')' and triggers reductions.
func (h *tokenHandler) handleCloseParen() {
	h.store.push(')')
	h.reducer.reduceAll()
}

// handleIdentifier stores an identifier with its text.
func (h *tokenHandler) handleIdentifier(tok Token) {
	text, err := h.lexer.TokenText(tok)
	if err != nil {
		// Invalid token bounds should not happen for valid tokens from the lexer.
		// If it does, skip the identifier - this is a safety fallback.
		return
	}
	if tok.Type == IDENT_QUOTED {
		text = stripIdentifierQuotes(text)
	}
	h.store.pushIdent(text)
}

// isNullKeywordContext checks if NULL should be kept as a keyword.
// Returns true when NULL follows IS or IS NOT.
func (h *tokenHandler) isNullKeywordContext() bool {
	if h.store.len() == 0 {
		return false
	}

	last := h.store.last()
	if last == IS {
		return true
	}

	// Check for IS NOT pattern
	if last == NOT_SYM && h.store.len() >= 2 {
		p := h.store.peek(2)
		if p[0] == IS {
			return true
		}
	}

	return false
}

// isNumericLiteral returns true if the token type is a numeric literal.
func isNumericLiteral(t int) bool {
	switch t {
	case NUM, LONG_NUM, ULONGLONG_NUM, DECIMAL_NUM, FLOAT_NUM, BIN_NUM, HEX_NUM:
		return true
	}
	return false
}

// isStringLiteral returns true if the token type is a string-like literal.
func isStringLiteral(t int) bool {
	switch t {
	case LEX_HOSTNAME, TEXT_STRING, NCHAR_STRING, PARAM_MARKER:
		return true
	}
	return false
}
