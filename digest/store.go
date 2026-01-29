package digest

import (
	"crypto/md5"
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
	version    MySQLVersion
}

// newTokenStore creates a new token store with pre-allocated capacity.
func newTokenStore(version MySQLVersion) *tokenStore {
	return &tokenStore{
		tokens:     make([]storedToken, 0, 256),
		tokenArray: make([]byte, 0, 1024),
		version:    version,
	}
}

// push adds a simple token to both stores.
// Binary format: 2 bytes (little-endian token type).
func (s *tokenStore) push(tokType int) {
	s.tokens = append(s.tokens, storedToken{tokType: tokType})
	binTok := s.translateToken(tokType)
	s.tokenArray = append(s.tokenArray,
		byte(binTok&0xff),
		byte((binTok>>8)&0xff))
}

// pushIdent adds an identifier token with its text to both stores.
// Binary format: 2 bytes (token) + 2 bytes (length) + N bytes (text).
func (s *tokenStore) pushIdent(text string) {
	s.tokens = append(s.tokens, storedToken{tokType: TOK_IDENT, text: text})

	binTok := s.translateToken(TOK_IDENT)
	s.tokenArray = append(s.tokenArray,
		byte(binTok&0xff),
		byte((binTok>>8)&0xff),
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

// translateToken converts a MySQL 8.0 token ID to the appropriate version.
// For MySQL 5.7, it uses the mysql80To57TokenMap mapping.
// For MySQL 8.0, it returns the token unchanged.
// ASCII tokens (0-255) are version-independent and pass through unchanged.
func (s *tokenStore) translateToken(tokType int) int {
	if s.version == MySQL57 {
		// ASCII characters are identical across versions
		if tokType < 256 {
			return tokType
		}
		if mapped, ok := mysql80To57TokenMap[tokType]; ok {
			return mapped
		}
		// Token not found in mapping, use TOK_UNUSED equivalent
		return m57TOK_UNUSED
	}
	return tokType
}

// computeHash returns the hash of the binary token array as a hex string.
// Uses SHA-256 for MySQL 8.0 (default) or MD5 for MySQL 5.7.
func (s *tokenStore) computeHash() string {
	if s.version == MySQL57 {
		hash := md5.Sum(s.tokenArray)
		return hex.EncodeToString(hash[:])
	}
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
