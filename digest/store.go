package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// storedToken represents a token in the normalized token stack.
type storedToken struct {
	tokType int
	text    string // Only set for identifier tokens (TOK_IDENT)
}

// tokenStore manages the dual storage of normalized tokens:
// - tokens: slice of storedToken for building human-readable output
// - tokenArray: binary byte array for computing the hash (MySQL format)
type tokenStore struct {
	tokens     []storedToken
	tokenArray []byte
}

// newTokenStore creates a new token store with pre-allocated capacity.
func newTokenStore() *tokenStore {
	return &tokenStore{
		tokens:     make([]storedToken, 0, 256),
		tokenArray: make([]byte, 0, 1024),
	}
}

// push adds a simple token to both stores.
// Binary format: 2 bytes (little-endian token type).
func (s *tokenStore) push(tokType int) {
	s.tokens = append(s.tokens, storedToken{tokType: tokType})
	s.tokenArray = append(s.tokenArray,
		byte(tokType&0xff),
		byte((tokType>>8)&0xff))
}

// pushIdent adds an identifier token with its text to both stores.
// Binary format: 2 bytes (token) + 2 bytes (length) + N bytes (text).
func (s *tokenStore) pushIdent(text string) {
	s.tokens = append(s.tokens, storedToken{tokType: TOK_IDENT, text: text})

	s.tokenArray = append(s.tokenArray,
		byte(TOK_IDENT&0xff),
		byte((TOK_IDENT>>8)&0xff),
		byte(len(text)&0xff),
		byte((len(text)>>8)&0xff))
	s.tokenArray = append(s.tokenArray, text...)
}

// pop removes the last n tokens from both stores.
// Each token removal removes 2 bytes from the binary array.
func (s *tokenStore) pop(n int) {
	if n <= 0 || n > len(s.tokens) {
		return
	}
	s.tokens = s.tokens[:len(s.tokens)-n]

	bytesToRemove := n * 2
	if bytesToRemove > len(s.tokenArray) {
		bytesToRemove = len(s.tokenArray)
	}
	s.tokenArray = s.tokenArray[:len(s.tokenArray)-bytesToRemove]
}

// peek returns the token types of the last n tokens.
// Index 0 is the oldest of the peeked tokens, index n-1 is the most recent.
// Missing positions are filled with TOK_UNUSED.
func (s *tokenStore) peek(n int) []int {
	result := make([]int, n)
	for i := 0; i < n; i++ {
		result[i] = TOK_UNUSED
	}

	storeLen := len(s.tokens)
	for i := 0; i < n && i < storeLen; i++ {
		// Map: result[n-1-i] = tokens[storeLen-1-i]
		// i=0: result[n-1] = tokens[storeLen-1] (most recent → last position)
		// i=1: result[n-2] = tokens[storeLen-2] (second most recent → second-to-last)
		result[n-1-i] = s.tokens[storeLen-1-i].tokType
	}
	return result
}

// last returns the most recent token type, or TOK_UNUSED if empty.
func (s *tokenStore) last() int {
	if len(s.tokens) == 0 {
		return TOK_UNUSED
	}
	return s.tokens[len(s.tokens)-1].tokType
}

// len returns the number of tokens in the store.
func (s *tokenStore) len() int {
	return len(s.tokens)
}

// computeHash returns the SHA-256 hash of the binary token array as a hex string.
func (s *tokenStore) computeHash() string {
	hash := sha256.Sum256(s.tokenArray)
	return hex.EncodeToString(hash[:])
}

// buildText converts the token stack to a normalized SQL string.
// Uses MySQL's delayed-space approach: space is added before a token
// only if the previous token requested it.
func (s *tokenStore) buildText(maxLen int) string {
	var b strings.Builder
	addSpace := false

	for _, tok := range s.tokens {
		text := tokenToText(tok)
		if text == "" {
			continue
		}
		if addSpace {
			b.WriteByte(' ')
		}
		b.WriteString(text)
		addSpace = TokenAppendSpace(tok.tokType)
	}

	result := b.String()
	if maxLen > 0 && len(result) > maxLen {
		result = result[:maxLen] + "..."
	}
	return result
}

// removeTrailingSemicolon removes the trailing semicolon token if present.
func (s *tokenStore) removeTrailingSemicolon() {
	if len(s.tokens) > 0 && s.tokens[len(s.tokens)-1].tokType == ';' {
		s.pop(1)
	}
}

// tokenToText converts a stored token to its output string representation.
func tokenToText(tok storedToken) string {
	if tok.tokType == TOK_IDENT {
		return "`" + escapeBackticks(tok.text) + "`"
	}
	text := TokenString(tok.tokType)
	if text == "(unknown)" {
		return ""
	}
	return text
}
