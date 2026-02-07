package internal

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

type storedToken struct {
	tokType int
	text    string
}

type tokenStore struct {
	tokens      []storedToken
	tokenArray  []byte
	version     MySQLVersion
	tokenConfig *TokenConfig
}

// TokenStore holds the normalized tokens for digest computation.
type TokenStore = tokenStore

// NewTokenStore creates a new token store for the given MySQL version.
func NewTokenStore(version MySQLVersion) *tokenStore {
	return &tokenStore{
		tokens:      make([]storedToken, 0, 256),
		tokenArray:  make([]byte, 0, 1024),
		version:     version,
		tokenConfig: GetTokenConfig(version),
	}
}

func (s *tokenStore) push(tokType int) {
	s.tokens = append(s.tokens, storedToken{tokType: tokType})
	binTok := s.translateToken(tokType)
	s.tokenArray = append(s.tokenArray,
		byte(binTok&0xff),
		byte((binTok>>8)&0xff))
}

// Binary format for identifiers: 2 bytes (token) + 2 bytes (length) + N bytes (text).
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

// peek2 returns the last two token types (second-to-last, last).
// Returns TOK_UNUSED for missing positions.
func (s *tokenStore) peek2() (int, int) {
	n := len(s.tokens)
	t1, t0 := TOK_UNUSED, TOK_UNUSED
	if n >= 1 {
		t0 = s.tokens[n-1].tokType
	}
	if n >= 2 {
		t1 = s.tokens[n-2].tokType
	}
	return t1, t0
}

// peek3 returns the last three token types (third-to-last, second-to-last, last).
// Returns TOK_UNUSED for missing positions.
func (s *tokenStore) peek3() (int, int, int) {
	n := len(s.tokens)
	t2, t1, t0 := TOK_UNUSED, TOK_UNUSED, TOK_UNUSED
	if n >= 1 {
		t0 = s.tokens[n-1].tokType
	}
	if n >= 2 {
		t1 = s.tokens[n-2].tokType
	}
	if n >= 3 {
		t2 = s.tokens[n-3].tokType
	}
	return t2, t1, t0
}

func (s *tokenStore) last() int {
	if len(s.tokens) == 0 {
		return TOK_UNUSED
	}
	return s.tokens[len(s.tokens)-1].tokType
}

func (s *tokenStore) len() int {
	return len(s.tokens)
}

func (s *tokenStore) translateToken(tokType int) int {
	return s.tokenConfig.TranslateForHash(tokType)
}

// ComputeHash returns the digest hash.
func (s *tokenStore) ComputeHash() string {
	if s.version == MySQL57 {
		hash := md5.Sum(s.tokenArray)
		return hex.EncodeToString(hash[:])
	}
	hash := sha256.Sum256(s.tokenArray)
	return hex.EncodeToString(hash[:])
}

// BuildText returns the normalized query text.
func (s *tokenStore) BuildText(maxLen int) string {
	var b strings.Builder
	addSpace := false

	for _, tok := range s.tokens {
		text := s.tokenToText(tok)
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

func (s *tokenStore) removeTrailingSemicolon() {
	if len(s.tokens) > 0 && s.tokens[len(s.tokens)-1].tokType == ';' {
		s.pop(1)
	}
}

func (s *tokenStore) tokenToText(tok storedToken) string {
	if tok.tokType == TOK_IDENT {
		return "`" + escapeBackticks(tok.text) + "`"
	}
	text := s.tokenConfig.GetString(tok.tokType)
	if text == "(unknown)" {
		return ""
	}
	return text
}
