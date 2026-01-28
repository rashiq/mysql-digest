package digest

import (
"encoding/hex"
"fmt"
"testing"
)

func TestDebugMySQL57Token(t *testing.T) {
sql := "SELECT * FROM users WHERE id = 42"

// MySQL 8.0 version
opt80 := Options{Version: MySQL80}
lexer80 := NewLexer(sql)
store80 := newTokenStore(MySQL80)
reducer80 := newReducer(store80)
handler80 := newTokenHandler(lexer80, store80, reducer80)
handler80.processAll()

fmt.Println("=== MySQL 8.0 ===")
fmt.Printf("Token array (%d bytes): %s\n", len(store80.tokenArray), hex.EncodeToString(store80.tokenArray))
fmt.Printf("Tokens: ")
for i, tok := range store80.tokens {
if i > 0 {
fmt.Print(", ")
}
if tok.text != "" {
fmt.Printf("%d(%s)", tok.tokType, tok.text)
} else {
fmt.Printf("%d", tok.tokType)
}
}
fmt.Println()
d80 := Normalize(sql, opt80)
fmt.Printf("Hash: %s\n", d80.Hash)
fmt.Printf("Text: %s\n", d80.Text)

// MySQL 5.7 version
opt57 := Options{Version: MySQL57}
lexer57 := NewLexer(sql)
store57 := newTokenStore(MySQL57)
reducer57 := newReducer(store57)
handler57 := newTokenHandler(lexer57, store57, reducer57)
handler57.processAll()

fmt.Println("\n=== MySQL 5.7 ===")
fmt.Printf("Token array (%d bytes): %s\n", len(store57.tokenArray), hex.EncodeToString(store57.tokenArray))
fmt.Printf("Tokens: ")
for i, tok := range store57.tokens {
if i > 0 {
fmt.Print(", ")
}
if tok.text != "" {
fmt.Printf("%d(%s)->%d", tok.tokType, tok.text, mysql57TokenID[tok.tokType])
} else {
mapped := tok.tokType
if m, ok := mysql57TokenID[tok.tokType]; ok {
mapped = m
}
fmt.Printf("%d->%d", tok.tokType, mapped)
}
}
fmt.Println()
d57 := Normalize(sql, opt57)
fmt.Printf("Hash: %s\n", d57.Hash)
fmt.Printf("Text: %s\n", d57.Text)
fmt.Printf("Expected: 731d9efe96031900ba2a36667f4718d0\n")
}
