package internal

type StateHandler func(l *Lexer) lexResult

const maxLexState = int(MY_LEX_STRING_OR_DELIMITER) + 1

var stateHandlers [maxLexState]StateHandler

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
	stateHandlers[MY_LEX_END_LONG_COMMENT] = (*Lexer).handleAsterisks

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
	stateHandlers[MY_LEX_USER_END] = (*Lexer).handleUserVariable
	stateHandlers[MY_LEX_HOSTNAME] = (*Lexer).handleHostname
	stateHandlers[MY_LEX_SYSTEM_VAR] = (*Lexer).handleSystemVar
	stateHandlers[MY_LEX_IDENT_OR_KEYWORD] = (*Lexer).handleIdentOrKeyword

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
