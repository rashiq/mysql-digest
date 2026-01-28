package digest

import (
"encoding/hex"
"fmt"
"testing"
)

func TestDebugJoin(t *testing.T) {
sql := "SELECT a.id, b.name FROM users a JOIN orders b ON a.id = b.user_id WHERE b.total > 100 AND a.status = 'active'"

lexer := NewLexer(sql)
store := newTokenStore(MySQL57)
reducer := newReducer(store)
handler := newTokenHandler(lexer, store, reducer)
handler.processAll()

fmt.Printf("SQL: %s\n", sql)
fmt.Printf("Token array len: %d bytes\n", len(store.tokenArray))
fmt.Printf("Token array hex: %s\n", hex.EncodeToString(store.tokenArray))

// List the tokens
fmt.Printf("\nTokens:\n")
for i, tok := range store.tokens {
m57, ok := mysql57TokenID[tok.tokType]
if !ok {
m57 = tok.tokType
}
if tok.text != "" {
fmt.Printf("  %d: %d -> %d (%s)\n", i, tok.tokType, m57, tok.text)
} else {
name := TokenString(tok.tokType)
fmt.Printf("  %d: %d -> %d (%s)\n", i, tok.tokType, m57, name)
}
}

fmt.Printf("\nOur hash: %s\n", store.computeHash())
fmt.Printf("Expected: 5b83e9f6a57a0c372448944f7ffaafb3\n")
}
