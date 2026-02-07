package internal

type tokenHandler struct {
	lexer   *Lexer
	store   *tokenStore
	reducer *reducer
}

// NewTokenHandler creates a new token handler.
func NewTokenHandler(lexer *Lexer, store *tokenStore, reducer *reducer) *tokenHandler {
	return &tokenHandler{
		lexer:   lexer,
		store:   store,
		reducer: reducer,
	}
}

func (h *tokenHandler) ProcessAll() error {
	for {
		tok := h.lexer.Lex()

		if tok.Type == END_OF_INPUT {
			h.store.removeTrailingSemicolon()
			return nil
		}
		if tok.Type == ABORT_SYM {
			return tok.Err
		}

		if err := h.handleToken(tok); err != nil {
			return err
		}
	}
}

func (h *tokenHandler) handleToken(tok Token) error {
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
		return h.handleIdentifier(tok)

	default:
		h.store.push(tok.Type)
		h.reducer.reduceAll()
	}
	return nil
}

// Absorbs any preceding unary +/- signs before normalizing.
func (h *tokenHandler) handleNumericLiteral() {
	h.reducer.reduceUnarySign()
	h.store.push(TOK_GENERIC_VALUE)
	h.reducer.reduceAfterValue()
}

func (h *tokenHandler) handleLiteral() {
	h.store.push(TOK_GENERIC_VALUE)
	h.reducer.reduceAfterValue()
}

// NULL is kept as a keyword after IS/IS NOT, otherwise normalized to a value.
func (h *tokenHandler) handleNull() {
	if h.isNullKeywordContext() {
		h.store.push(NULL_SYM)
	} else {
		h.store.push(TOK_GENERIC_VALUE)
		h.reducer.reduceAfterValue()
	}
}

func (h *tokenHandler) handleCloseParen() {
	h.store.push(')')
	h.reducer.reduceAll()
}

func (h *tokenHandler) handleIdentifier(tok Token) error {
	text, err := h.lexer.TokenText(tok)
	if err != nil {
		return err
	}
	if tok.Type == IDENT_QUOTED {
		text = stripIdentifierQuotes(text)
	}
	h.store.pushIdent(text)
	return nil
}

// isNullKeywordContext checks if NULL should be kept as a keyword.
// Returns true for IS NULL or IS NOT NULL.
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
		prev, _ := h.store.peek2()
		if prev == IS {
			return true
		}
	}

	return false
}

func isNumericLiteral(t int) bool {
	switch t {
	case NUM, LONG_NUM, ULONGLONG_NUM, DECIMAL_NUM, FLOAT_NUM, BIN_NUM, HEX_NUM:
		return true
	}
	return false
}

func isStringLiteral(t int) bool {
	switch t {
	case LEX_HOSTNAME, TEXT_STRING, NCHAR_STRING, PARAM_MARKER:
		return true
	}
	return false
}
