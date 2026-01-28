package digest

import (
"encoding/hex"
"fmt"
"testing"
)

func TestDebugDelete(t *testing.T) {
sql := "DELETE FROM sessions WHERE expires_at < NOW() AND user_id IN (1, 2, 3, 4, 5)"

lexer := NewLexer(sql)
store := newTokenStore(MySQL57)
reducer := newReducer(store)
handler := newTokenHandler(lexer, store, reducer)
handler.processAll()

fmt.Printf("SQL: %s\n", sql)
fmt.Printf("Token array hex: %s\n", hex.EncodeToString(store.tokenArray))

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
fmt.Printf("Expected: b5f970dd3ff955057de7f2d8b6e7e2aa\n")
}
