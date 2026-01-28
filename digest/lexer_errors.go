package digest

import "fmt"

// LexError represents an error encountered during lexical analysis.
type LexError struct {
	Position int
	Message  string
	Input    string
}

func (e *LexError) Error() string {
	if e.Input != "" {
		return fmt.Sprintf("lex error at position %d: %s (near %q)", e.Position, e.Message, e.Input)
	}
	return fmt.Sprintf("lex error at position %d: %s", e.Position, e.Message)
}

// Error message constants
const (
	ErrUnterminatedString   = "unterminated string literal"
	ErrUnterminatedComment  = "unterminated block comment"
	ErrUnterminatedIdent    = "unterminated quoted identifier"
	ErrUnterminatedHint     = "unterminated optimizer hint"
	ErrInvalidHexLiteral    = "invalid hex literal"
	ErrInvalidBinaryLiteral = "invalid binary literal"
	ErrInvalidTokenBounds   = "invalid token bounds"
	ErrUnterminatedDollar   = "unterminated dollar-quoted string"
)

func NewLexError(position int, message string, input string) *LexError {
	return &LexError{Position: position, Message: message, Input: input}
}
