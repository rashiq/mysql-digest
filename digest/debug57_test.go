package digest

import (
"encoding/hex"
"fmt"
"testing"
)

func TestDebug57Simple(t *testing.T) {
sql := "SELECT 1"

lexer := NewLexer(sql)
store := newTokenStore(MySQL57)
reducer := newReducer(store)
handler := newTokenHandler(lexer, store, reducer)
handler.processAll()

fmt.Printf("SQL: %s\n", sql)
fmt.Printf("Token array hex: %s\n", hex.EncodeToString(store.tokenArray))
fmt.Printf("Token array len: %d bytes\n", len(store.tokenArray))

// Print each 2-byte token
fmt.Printf("\nTokens:\n")
for i := 0; i+1 < len(store.tokenArray); i += 2 {
tokVal := int(store.tokenArray[i]) | (int(store.tokenArray[i+1]) << 8)
fmt.Printf("  [%d-%d]: 0x%04x (%d)\n", i, i+1, tokVal, tokVal)
}

fmt.Printf("\nOur hash: %s\n", store.computeHash())
fmt.Printf("Expected: 3d4fc22e33e10d7235eced3c75a84c2c\n")
}
