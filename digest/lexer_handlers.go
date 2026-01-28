package digest

// State handler methods for the MySQL lexer.
// These extract complex state handling logic from the main Lex() switch statement.

// lexResult represents the result of a state handler.
// It wraps a Token with information about whether to return or continue.
type lexResult struct {
	token        Token
	done         bool     // If true, return the token. If false, continue with nextState.
	nextState    LexState // State to transition to if not done (current Lex() loop)
	setNextLex   bool     // If true, set l.nextState to nextLexState
	nextLexState LexState // State for the NEXT Lex() call (persists across calls)
}

// done creates a lexResult that returns a token.
func done(t Token) lexResult {
	return lexResult{token: t, done: true}
}

// doneWithNext creates a lexResult that returns a token and sets the next Lex() state.
func doneWithNext(t Token, nextLex LexState) lexResult {
	return lexResult{token: t, done: true, setNextLex: true, nextLexState: nextLex}
}

// cont creates a lexResult that continues to another state.
func cont(state LexState) lexResult {
	return lexResult{done: false, nextState: state}
}

// handleChar handles MY_LEX_CHAR state - single character tokens and special sequences.
func (l *Lexer) handleChar(c byte) lexResult {
	// Check for special two-char sequences with '-'
	if c == '-' && l.peek() == '-' {
		// Check for "-- " comment (-- followed by space or control char)
		nextChar := l.peekN(1)
		if isSpace(nextChar) || isCntrl(nextChar) {
			return cont(MY_LEX_COMMENT)
		}
	}

	// Check for JSON arrow operators
	if c == '-' && l.peek() == '>' {
		l.skip() // consume '>'
		if l.peek() == '>' {
			l.skip() // consume second '>'
			return doneWithNext(Token{Type: JSON_UNQUOTED_SEPARATOR_SYM, Start: l.tokStart, End: l.pos}, MY_LEX_START)
		}
		return doneWithNext(Token{Type: JSON_SEPARATOR_SYM, Start: l.tokStart, End: l.pos}, MY_LEX_START)
	}

	// Check for placeholder '?' in prepare mode
	if c == '?' && l.stmtPrepareMode && !isIdentChar(l.peek()) {
		return doneWithNext(Token{Type: PARAM_MARKER, Start: l.tokStart, End: l.pos}, MY_LEX_START)
	}

	// Close paren does NOT allow signed numbers after (don't set nextState)
	if c == ')' {
		return done(Token{Type: int(c), Start: l.tokStart, End: l.pos})
	}

	// All other chars set nextState = MY_LEX_START to allow signed numbers
	return doneWithNext(Token{Type: int(c), Start: l.tokStart, End: l.pos}, MY_LEX_START)
}

// handleIdent handles MY_LEX_IDENT state - identifier scanning with keyword lookup.
func (l *Lexer) handleIdent() lexResult {
	// Scan identifier - first char was already consumed
	for isIdentChar(l.peek()) {
		l.skip()
	}

	length := l.tokenLen()

	// Check if followed by '.' and identifier char
	if l.peek() == '.' && isIdentChar(l.peekN(1)) {
		// Still do keyword lookup for system variable scopes
		if tokval := l.findKeyword(length, false); tokval != 0 {
			return doneWithNext(l.returnToken(Token{Type: tokval, Start: l.tokStart, End: l.tokStart + length}), MY_LEX_IDENT_SEP)
		}
		return doneWithNext(l.returnToken(Token{Type: IDENT, Start: l.tokStart, End: l.tokStart + length}), MY_LEX_IDENT_SEP)
	}

	l.backup() // Unget the non-ident char

	// Check if it's a keyword
	nextChar := l.peekN(1)
	if tokval := l.findKeyword(length, nextChar == '('); tokval != 0 {
		l.skip() // Re-skip the character we ungot
		return doneWithNext(l.returnToken(Token{Type: tokval, Start: l.tokStart, End: l.tokStart + length}), MY_LEX_START)
	}
	l.skip() // Re-skip

	// Return as IDENT
	return done(l.returnToken(Token{Type: IDENT, Start: l.tokStart, End: l.tokStart + length}))
}

// handleIdentSep handles MY_LEX_IDENT_SEP state - identifier separator (dot between parts).
func (l *Lexer) handleIdentSep() lexResult {
	c := l.advance()
	if isIdentChar(l.peek()) {
		return doneWithNext(Token{Type: int(c), Start: l.tokStart, End: l.pos}, MY_LEX_IDENT_START)
	}
	return doneWithNext(Token{Type: int(c), Start: l.tokStart, End: l.pos}, MY_LEX_START)
}

// handleIdentStart handles MY_LEX_IDENT_START state - scanning identifier after separator.
func (l *Lexer) handleIdentStart() lexResult {
	for isIdentChar(l.peek()) {
		l.skip()
	}
	length := l.tokenLen()
	if l.peek() == '.' && isIdentChar(l.peekN(1)) {
		return doneWithNext(Token{Type: IDENT, Start: l.tokStart, End: l.tokStart + length}, MY_LEX_IDENT_SEP)
	}
	return done(Token{Type: IDENT, Start: l.tokStart, End: l.tokStart + length})
}

// handleCmpOp handles MY_LEX_CMP_OP state - comparison operators >, >=, =, !=.
func (l *Lexer) handleCmpOp() lexResult {
	nextState := getStateMap(l.peek())
	if nextState == MY_LEX_CMP_OP || nextState == MY_LEX_LONG_CMP_OP {
		l.skip()
	}
	length := l.tokenLen()
	if tokval := l.findKeyword(length, false); tokval != 0 {
		return doneWithNext(Token{Type: tokval, Start: l.tokStart, End: l.pos}, MY_LEX_START)
	}
	return cont(MY_LEX_CHAR)
}

// handleLongCmpOp handles MY_LEX_LONG_CMP_OP state - operators like <, <=, <>, <=>.
func (l *Lexer) handleLongCmpOp() lexResult {
	nextState := getStateMap(l.peek())
	if nextState == MY_LEX_CMP_OP || nextState == MY_LEX_LONG_CMP_OP {
		l.skip()
		if getStateMap(l.peek()) == MY_LEX_CMP_OP {
			l.skip()
		}
	}
	length := l.tokenLen()
	if tokval := l.findKeyword(length, false); tokval != 0 {
		return doneWithNext(Token{Type: tokval, Start: l.tokStart, End: l.pos}, MY_LEX_START)
	}
	return cont(MY_LEX_CHAR)
}

// handleBool handles MY_LEX_BOOL state - && and || operators.
func (l *Lexer) handleBool(c byte) lexResult {
	if l.peek() != c {
		return done(Token{Type: int(c), Start: l.tokStart, End: l.pos})
	}
	l.skip()
	if tokval := l.findKeyword(2, false); tokval != 0 {
		return doneWithNext(Token{Type: tokval, Start: l.tokStart, End: l.pos}, MY_LEX_START)
	}
	return done(Token{Type: int(c), Start: l.tokStart, End: l.pos})
}

// handleSetVar handles MY_LEX_SET_VAR state - := operator.
func (l *Lexer) handleSetVar(c byte) lexResult {
	if l.peek() != '=' {
		return done(Token{Type: int(c), Start: l.tokStart, End: l.pos})
	}
	l.skip()
	return done(Token{Type: SET_VAR, Start: l.tokStart, End: l.pos})
}
