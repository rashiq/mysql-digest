package digest

import (
"fmt"
"testing"
)

func TestDebugHint(t *testing.T) {
sql := "SELECT /*+ MAX_EXECUTION_TIME(5000) */ id, name FROM large_table WHERE status = 1 LIMIT 1000"

lexer := NewLexer(sql)
store := newTokenStore(MySQL57)
reducer := newReducer(store)
handler := newTokenHandler(lexer, store, reducer)
handler.processAll()

fmt.Printf("SQL: %s\n", sql)

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

d := Normalize(sql, Options{Version: MySQL57})
fmt.Printf("\nOur hash: %s\n", d.Hash)
fmt.Printf("Expected: 0673f7e6d08e8d422618b2ee0e6700dd\n")
fmt.Printf("Our text: %s\n", d.Text)
}
