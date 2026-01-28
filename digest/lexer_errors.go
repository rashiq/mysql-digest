package digest

import "fmt"

// LexError represents an error encountered during lexical analysis.
// It captures the position, type, and context of the error for precise diagnostics.
type LexError struct {
	Position int    // Position in input where error occurred
	Message  string // Human-readable error message
	Input    string // The problematic input fragment (for context)
}

// Error implements the error interface.
func (e *LexError) Error() string {
	if e.Input != "" {
		return fmt.Sprintf("lex error at position %d: %s (near %q)", e.Position, e.Message, e.Input)
	}
	return fmt.Sprintf("lex error at position %d: %s", e.Position, e.Message)
}

// Error message constants for consistent error reporting.
// These match the error conditions that can occur during MySQL lexical analysis.
const (
	ErrUnterminatedString   = "unterminated string literal"
	ErrUnterminatedComment  = "unterminated block comment"
	ErrUnterminatedIdent    = "unterminated quoted identifier"
	ErrUnterminatedHint     = "unterminated optimizer hint"
	ErrInvalidHexLiteral    = "invalid hex literal"
	ErrInvalidBinaryLiteral = "invalid binary literal"
	ErrInvalidTokenBounds   = "invalid token bounds"
	ErrInvalidEscapeSeq     = "invalid escape sequence"
	ErrUnterminatedDollar   = "unterminated dollar-quoted string"
)

// NewLexError creates a new LexError with the given parameters.
func NewLexError(position int, message string, input string) *LexError {
	return &LexError{
		Position: position,
		Message:  message,
		Input:    input,
	}
}

// NewLexErrorFromToken creates a LexError from a token and lexer context.
// This is useful for creating errors after lexing when we have a token
// but need error details.
func NewLexErrorFromToken(l *Lexer, t Token, message string) *LexError {
	input := ""
	if t.Start >= 0 && t.End <= len(l.input) && t.Start <= t.End {
		// Limit context to avoid huge error messages
		maxLen := 50
		if t.End-t.Start > maxLen {
			input = l.input[t.Start:t.Start+maxLen] + "..."
		} else {
			input = l.input[t.Start:t.End]
		}
	}
	return &LexError{
		Position: t.Start,
		Message:  message,
		Input:    input,
	}
}

// IsAbortToken returns true if the token type represents an error/abort condition.
// This provides a clean way to check for error tokens without magic number comparisons.
func IsAbortToken(tokenType int) bool {
	return tokenType == ABORT_SYM
}
