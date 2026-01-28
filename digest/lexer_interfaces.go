package digest

// Interfaces for the MySQL lexer.
// These enable dependency injection for testing and extensibility.

// KeywordResolver resolves identifiers to keyword token types.
// This abstraction allows for custom keyword sets or testing with mock resolvers.
type KeywordResolver interface {
	// Resolve looks up a keyword and returns its token type.
	// Returns the token type and true if found, or 0 and false if not a keyword.
	Resolve(text string) (tokenType int, found bool)

	// ResolveHint looks up an optimizer hint keyword.
	// Returns the token type and true if found, or 0 and false if not a hint keyword.
	ResolveHint(text string) (tokenType int, found bool)
}

// StateMapper maps characters to initial lexer states.
// This abstraction allows for custom state mappings or testing.
type StateMapper interface {
	// GetState returns the initial lexer state for a given byte.
	GetState(c byte) LexState
}

// DefaultKeywordResolver implements KeywordResolver using the standard keyword maps.
type DefaultKeywordResolver struct{}

// Resolve looks up a keyword in TokenKeywords.
func (r DefaultKeywordResolver) Resolve(text string) (int, bool) {
	upper := toUpper(text)
	tok, ok := TokenKeywords[upper]
	return tok, ok
}

// ResolveHint looks up a hint keyword in HintKeywords.
func (r DefaultKeywordResolver) ResolveHint(text string) (int, bool) {
	upper := toUpper(text)
	tok, ok := HintKeywords[upper]
	return tok, ok
}

// DefaultStateMapper implements StateMapper using the standard state map.
type DefaultStateMapper struct{}

// GetState returns the initial lexer state for a given byte.
func (m DefaultStateMapper) GetState(c byte) LexState {
	return getStateMap(c)
}

// Compile-time interface compliance checks
var _ KeywordResolver = DefaultKeywordResolver{}
var _ StateMapper = DefaultStateMapper{}
