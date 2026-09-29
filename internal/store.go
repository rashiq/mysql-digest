package internal

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"strings"
	"unicode/utf8"
)

type storedToken struct {
	tokType int
	text    string
}

type tokenStore struct {
	tokens      []storedToken
	tokenConfig *TokenConfig
}

func NewTokenStore(config *TokenConfig) *tokenStore {
	return &tokenStore{
		tokens:      make([]storedToken, 0, 32),
		tokenConfig: config,
	}
}

func (s *tokenStore) push(tokType int) {
	s.tokens = append(s.tokens, storedToken{tokType: tokType})
}

func (s *tokenStore) pushIdent(text string) {
	s.tokens = append(s.tokens, storedToken{tokType: TOK_IDENT, text: text})
}

func (s *tokenStore) pop(n int) {
	if n <= 0 || n > len(s.tokens) {
		return
	}
	s.tokens = s.tokens[:len(s.tokens)-n]
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

// ComputeHash returns the digest hash.
func (s *tokenStore) ComputeHash() string {
	size := 2 * len(s.tokens)
	for _, tok := range s.tokens {
		if tok.tokType == TOK_IDENT {
			size += 2 + len(tok.text)
		}
	}
	data := make([]byte, 0, size)
	for _, tok := range s.tokens {
		data = binary.LittleEndian.AppendUint16(data, uint16(s.tokenConfig.TranslateForHash(tok.tokType)))
		if tok.tokType == TOK_IDENT {
			data = binary.LittleEndian.AppendUint16(data, uint16(len(tok.text)))
			data = append(data, tok.text...)
		}
	}
	if s.tokenConfig.Version == MySQL57 {
		hash := md5.Sum(data)
		return hex.EncodeToString(hash[:])
	}
	hash := sha256.Sum256(data)
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
		if maxLen <= 3 {
			return strings.Repeat(".", maxLen)
		}
		end := maxLen - 3
		for end > 0 && !utf8.RuneStart(result[end]) {
			end--
		}
		result = result[:end] + "..."
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
