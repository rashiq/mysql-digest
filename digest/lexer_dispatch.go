package digest

// State dispatch table for the MySQL lexer.
// Handlers are registered in init() and dispatched by Lex().

// StateHandler is the signature for all state handler functions.
type StateHandler func(l *Lexer) lexResult

// stateHandlers maps lexer states to their handler functions.
var stateHandlers = make(map[LexState]StateHandler)

// dispatchState looks up and executes a handler for the given state.
func (l *Lexer) dispatchState(state LexState) (lexResult, bool) {
	if handler, ok := stateHandlers[state]; ok {
		return handler(l), true
	}
	return lexResult{}, false
}

func init() {
	// Core state machine
	stateHandlers[MY_LEX_START] = (*Lexer).handleStart
	stateHandlers[MY_LEX_SKIP] = (*Lexer).handleSkip
	stateHandlers[MY_LEX_EOL] = (*Lexer).handleEOL
	stateHandlers[MY_LEX_COMMENT] = (*Lexer).handleLineComment
	stateHandlers[MY_LEX_SEMICOLON] = (*Lexer).handleCharToken
	stateHandlers[MY_LEX_END_LONG_COMMENT] = (*Lexer).handleCharToken

	// Branching states
	stateHandlers[MY_LEX_IDENT_OR_HEX] = (*Lexer).handleIdentOrHex
	stateHandlers[MY_LEX_IDENT_OR_BIN] = (*Lexer).handleIdentOrBin
	stateHandlers[MY_LEX_INT_OR_REAL] = (*Lexer).handleIntOrReal
	stateHandlers[MY_LEX_REAL_OR_POINT] = (*Lexer).handleRealOrPoint
	stateHandlers[MY_LEX_STRING_OR_DELIMITER] = (*Lexer).handleStringOrDelimiter

	// Character and operators
	stateHandlers[MY_LEX_CHAR] = (*Lexer).handleChar
	stateHandlers[MY_LEX_CMP_OP] = (*Lexer).handleCmpOp
	stateHandlers[MY_LEX_LONG_CMP_OP] = (*Lexer).handleLongCmpOp
	stateHandlers[MY_LEX_BOOL] = (*Lexer).handleBool
	stateHandlers[MY_LEX_SET_VAR] = (*Lexer).handleSetVar

	// Identifiers
	stateHandlers[MY_LEX_IDENT] = (*Lexer).handleIdent
	stateHandlers[MY_LEX_IDENT_SEP] = (*Lexer).handleIdentSep
	stateHandlers[MY_LEX_IDENT_START] = (*Lexer).handleIdentStart
	stateHandlers[MY_LEX_IDENT_OR_NCHAR] = (*Lexer).handleNChar
	stateHandlers[MY_LEX_IDENT_OR_DOLLAR_QUOTED_TEXT] = (*Lexer).handleDollarQuoted

	// Numbers
	stateHandlers[MY_LEX_NUMBER_IDENT] = (*Lexer).handleNumberIdent
	stateHandlers[MY_LEX_HEX_NUMBER] = (*Lexer).handleHexNumber
	stateHandlers[MY_LEX_BIN_NUMBER] = (*Lexer).handleBinNumber
	stateHandlers[MY_LEX_REAL] = (*Lexer).handleReal

	// Strings
	stateHandlers[MY_LEX_STRING] = (*Lexer).handleString
	stateHandlers[MY_LEX_USER_VARIABLE_DELIMITER] = (*Lexer).handleQuotedIdent

	// Comments
	stateHandlers[MY_LEX_LONG_COMMENT] = (*Lexer).handleLongComment
}

// ---- Core state handlers ----

func (l *Lexer) handleStart() lexResult {
	for getStateMap(l.peek()) == MY_LEX_SKIP {
		l.skip()
	}
	l.startToken()
	c := l.advance()
	return cont(getStateMap(c))
}

func (l *Lexer) handleSkip() lexResult {
	l.skip()
	return cont(MY_LEX_START)
}

func (l *Lexer) handleEOL() lexResult {
	return done(Token{Type: END_OF_INPUT, Start: l.tokStart, End: l.pos})
}

func (l *Lexer) handleLineComment() lexResult {
	l.scanLineComment()
	return cont(MY_LEX_START)
}

func (l *Lexer) handleCharToken() lexResult {
	c := l.input[l.tokStart]
	return done(Token{Type: int(c), Start: l.tokStart, End: l.pos})
}

// ---- Branching handlers ----

func (l *Lexer) handleIdentOrHex() lexResult {
	if l.peek() == '\'' {
		return cont(MY_LEX_HEX_NUMBER)
	}
	return cont(MY_LEX_IDENT)
}

func (l *Lexer) handleIdentOrBin() lexResult {
	if l.peek() == '\'' {
		return cont(MY_LEX_BIN_NUMBER)
	}
	return cont(MY_LEX_IDENT)
}

func (l *Lexer) handleIntOrReal() lexResult {
	if l.peek() != '.' {
		return done(Token{Type: l.intToken(l.tokenLen()), Start: l.tokStart, End: l.pos})
	}
	l.skip()
	return cont(MY_LEX_REAL)
}

func (l *Lexer) handleRealOrPoint() lexResult {
	if isDigit(l.peek()) {
		return cont(MY_LEX_REAL)
	}
	c := l.input[l.tokStart]
	return done(Token{Type: int(c), Start: l.tokStart, End: l.pos})
}

func (l *Lexer) handleStringOrDelimiter() lexResult {
	if (l.sqlMode & MODE_ANSI_QUOTES) != 0 {
		return cont(MY_LEX_USER_VARIABLE_DELIMITER)
	}
	return cont(MY_LEX_STRING)
}
