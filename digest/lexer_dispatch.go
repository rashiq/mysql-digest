package digest

// State dispatch table for the MySQL lexer.
// Handlers are registered in init() and dispatched by Lex().

// StateHandler is the signature for all state handler functions.
type StateHandler func(l *Lexer) lexResult

// maxLexState is the number of lexer states (highest state value + 1).
// MY_LEX_STRING_OR_DELIMITER = 32, so we need 33 slots.
const maxLexState = int(MY_LEX_STRING_OR_DELIMITER) + 1

// stateHandlers is a fixed-size array mapping states to handlers.
// Array indexing is O(1) with no hash overhead, unlike map lookups.
var stateHandlers [maxLexState]StateHandler

// dispatchState executes the handler for the given state.
func (l *Lexer) dispatchState(state LexState) (lexResult, bool) {
	idx := int(state)
	if idx >= 0 && idx < maxLexState {
		if handler := stateHandlers[idx]; handler != nil {
			return handler(l), true
		}
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
	stateHandlers[MY_LEX_END_LONG_COMMENT] = (*Lexer).handleAsteriks

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
