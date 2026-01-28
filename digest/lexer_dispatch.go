package digest

// State dispatch table for the MySQL lexer.
// This implements the Open/Closed Principle - new states can be added
// by registering handlers without modifying the Lex() loop.

// StateHandler is the signature for state handler functions.
// c is the first character of the token (may be unused by some handlers).
// Returns a lexResult indicating whether to return a token or continue.
type StateHandler func(l *Lexer, c byte) lexResult

// stateHandlers maps lexer states to their handler functions.
// Handlers registered here are dispatched automatically by Lex().
var stateHandlers = make(map[LexState]StateHandler)

// RegisterStateHandler registers a handler for a lexer state.
// This allows extending the lexer without modifying the Lex() loop.
func RegisterStateHandler(state LexState, handler StateHandler) {
	stateHandlers[state] = handler
}

// dispatchState looks up and executes a handler for the given state.
// Returns the lexResult and true if a handler was found, or false if not.
func (l *Lexer) dispatchState(state LexState, c byte) (lexResult, bool) {
	if handler, ok := stateHandlers[state]; ok {
		return handler(l, c), true
	}
	return lexResult{}, false
}

// init registers all standard state handlers.
func init() {
	// Core state machine handlers
	RegisterStateHandler(MY_LEX_START, (*Lexer).dispatchStart)
	RegisterStateHandler(MY_LEX_SKIP, (*Lexer).dispatchSkip)
	RegisterStateHandler(MY_LEX_EOL, (*Lexer).dispatchEOL)
	RegisterStateHandler(MY_LEX_COMMENT, (*Lexer).dispatchComment)
	RegisterStateHandler(MY_LEX_SEMICOLON, (*Lexer).dispatchSemicolon)
	RegisterStateHandler(MY_LEX_END_LONG_COMMENT, (*Lexer).dispatchEndLongComment)

	// Branching state handlers
	RegisterStateHandler(MY_LEX_IDENT_OR_HEX, (*Lexer).dispatchIdentOrHex)
	RegisterStateHandler(MY_LEX_IDENT_OR_BIN, (*Lexer).dispatchIdentOrBin)
	RegisterStateHandler(MY_LEX_INT_OR_REAL, (*Lexer).dispatchIntOrReal)
	RegisterStateHandler(MY_LEX_REAL_OR_POINT, (*Lexer).dispatchRealOrPoint)
	RegisterStateHandler(MY_LEX_STRING_OR_DELIMITER, (*Lexer).dispatchStringOrDelimiter)

	// Character and operator handlers
	RegisterStateHandler(MY_LEX_CHAR, (*Lexer).dispatchChar)
	RegisterStateHandler(MY_LEX_CMP_OP, (*Lexer).dispatchCmpOp)
	RegisterStateHandler(MY_LEX_LONG_CMP_OP, (*Lexer).dispatchLongCmpOp)
	RegisterStateHandler(MY_LEX_BOOL, (*Lexer).dispatchBool)
	RegisterStateHandler(MY_LEX_SET_VAR, (*Lexer).dispatchSetVar)

	// Identifier handlers
	RegisterStateHandler(MY_LEX_IDENT, (*Lexer).dispatchIdent)
	RegisterStateHandler(MY_LEX_IDENT_SEP, (*Lexer).dispatchIdentSep)
	RegisterStateHandler(MY_LEX_IDENT_START, (*Lexer).dispatchIdentStart)
	RegisterStateHandler(MY_LEX_IDENT_OR_NCHAR, (*Lexer).dispatchNChar)
	RegisterStateHandler(MY_LEX_IDENT_OR_DOLLAR_QUOTED_TEXT, (*Lexer).dispatchDollarQuoted)

	// Numeric handlers
	RegisterStateHandler(MY_LEX_NUMBER_IDENT, (*Lexer).dispatchNumberIdent)
	RegisterStateHandler(MY_LEX_HEX_NUMBER, (*Lexer).dispatchHexNumber)
	RegisterStateHandler(MY_LEX_BIN_NUMBER, (*Lexer).dispatchBinNumber)
	RegisterStateHandler(MY_LEX_REAL, (*Lexer).dispatchReal)

	// String handlers
	RegisterStateHandler(MY_LEX_STRING, (*Lexer).dispatchString)
	RegisterStateHandler(MY_LEX_USER_VARIABLE_DELIMITER, (*Lexer).dispatchQuotedIdent)

	// Comment handlers
	RegisterStateHandler(MY_LEX_LONG_COMMENT, (*Lexer).dispatchLongComment)
}

// ---- Dispatch wrapper methods ----
// These adapt the existing handlers to the StateHandler signature.

func (l *Lexer) dispatchChar(c byte) lexResult {
	return l.handleChar(c)
}

func (l *Lexer) dispatchCmpOp(_ byte) lexResult {
	return l.handleCmpOp()
}

func (l *Lexer) dispatchLongCmpOp(_ byte) lexResult {
	return l.handleLongCmpOp()
}

func (l *Lexer) dispatchBool(c byte) lexResult {
	return l.handleBool(c)
}

func (l *Lexer) dispatchSetVar(c byte) lexResult {
	return l.handleSetVar(c)
}

func (l *Lexer) dispatchIdent(_ byte) lexResult {
	return l.handleIdent()
}

func (l *Lexer) dispatchIdentSep(_ byte) lexResult {
	return l.handleIdentSep()
}

func (l *Lexer) dispatchIdentStart(_ byte) lexResult {
	return l.handleIdentStart()
}

func (l *Lexer) dispatchNChar(_ byte) lexResult {
	return l.handleNChar()
}

func (l *Lexer) dispatchDollarQuoted(_ byte) lexResult {
	return l.handleDollarQuoted()
}

func (l *Lexer) dispatchNumberIdent(c byte) lexResult {
	return l.handleNumberIdent(c)
}

func (l *Lexer) dispatchHexNumber(_ byte) lexResult {
	return l.handleHexNumber()
}

func (l *Lexer) dispatchBinNumber(_ byte) lexResult {
	return l.handleBinNumber()
}

func (l *Lexer) dispatchReal(_ byte) lexResult {
	return l.handleReal()
}

func (l *Lexer) dispatchString(c byte) lexResult {
	return l.handleString(c)
}

func (l *Lexer) dispatchQuotedIdent(c byte) lexResult {
	return l.handleQuotedIdent(c)
}

func (l *Lexer) dispatchLongComment(c byte) lexResult {
	return l.handleLongComment(c)
}

// ---- Core state machine handlers ----

// dispatchStart handles MY_LEX_START: skip whitespace and determine next state.
// This is special because it sets 'c' which subsequent handlers may need.
// We use a special lexResult field to communicate the character.
func (l *Lexer) dispatchStart(_ byte) lexResult {
	// Skip leading whitespace
	for l.stateMapper.GetState(l.peek()) == MY_LEX_SKIP {
		l.skip()
	}
	// Start of real token
	l.restartToken()
	c := l.advance()
	nextState := l.stateMapper.GetState(c)
	// Continue to the next state with c available
	return lexResult{nextState: nextState, startChar: c}
}

// dispatchSkip handles MY_LEX_SKIP: skip a character and restart.
func (l *Lexer) dispatchSkip(_ byte) lexResult {
	l.skip()
	return cont(MY_LEX_START)
}

// dispatchEOL handles MY_LEX_EOL: return end of input token.
func (l *Lexer) dispatchEOL(_ byte) lexResult {
	return done(Token{Type: END_OF_INPUT, Start: l.tokStart, End: l.pos})
}

// dispatchComment handles MY_LEX_COMMENT: scan line comment and restart.
func (l *Lexer) dispatchComment(_ byte) lexResult {
	l.scanLineComment()
	return cont(MY_LEX_START)
}

// dispatchSemicolon handles MY_LEX_SEMICOLON: return semicolon token.
func (l *Lexer) dispatchSemicolon(c byte) lexResult {
	return done(Token{Type: int(c), Start: l.tokStart, End: l.pos})
}

// dispatchEndLongComment handles MY_LEX_END_LONG_COMMENT: return character token.
func (l *Lexer) dispatchEndLongComment(c byte) lexResult {
	return done(Token{Type: int(c), Start: l.tokStart, End: l.pos})
}

// ---- Branching state handlers ----

// dispatchIdentOrHex handles MY_LEX_IDENT_OR_HEX: branch to hex or ident.
func (l *Lexer) dispatchIdentOrHex(_ byte) lexResult {
	if l.peek() == '\'' {
		return cont(MY_LEX_HEX_NUMBER)
	}
	return cont(MY_LEX_IDENT)
}

// dispatchIdentOrBin handles MY_LEX_IDENT_OR_BIN: branch to bin or ident.
func (l *Lexer) dispatchIdentOrBin(_ byte) lexResult {
	if l.peek() == '\'' {
		return cont(MY_LEX_BIN_NUMBER)
	}
	return cont(MY_LEX_IDENT)
}

// dispatchIntOrReal handles MY_LEX_INT_OR_REAL: return int or continue to real.
func (l *Lexer) dispatchIntOrReal(_ byte) lexResult {
	if l.peek() != '.' {
		length := l.tokenLen()
		return done(Token{Type: l.intToken(length), Start: l.tokStart, End: l.pos})
	}
	l.skip()
	return cont(MY_LEX_REAL)
}

// dispatchRealOrPoint handles MY_LEX_REAL_OR_POINT: continue to real or return point.
func (l *Lexer) dispatchRealOrPoint(c byte) lexResult {
	if isDigit(l.peek()) {
		return cont(MY_LEX_REAL)
	}
	return done(Token{Type: int(c), Start: l.tokStart, End: l.pos})
}

// dispatchStringOrDelimiter handles MY_LEX_STRING_OR_DELIMITER: branch based on SQL mode.
func (l *Lexer) dispatchStringOrDelimiter(_ byte) lexResult {
	if (l.sqlMode & MODE_ANSI_QUOTES) != 0 {
		return cont(MY_LEX_USER_VARIABLE_DELIMITER)
	}
	return cont(MY_LEX_STRING)
}
