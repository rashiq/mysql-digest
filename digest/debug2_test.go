package digest

import (
"crypto/md5"
"encoding/hex"
"fmt"
"testing"
)

func TestDebugActualHashes(t *testing.T) {
// Let's verify the hash computation is correct
sql := "SELECT * FROM users WHERE id = 42"

d57 := Normalize(sql, Options{Version: MySQL57})

fmt.Println("=== Verification ===")
fmt.Printf("SQL: %s\n", sql)
fmt.Printf("Text: %s\n", d57.Text)
fmt.Printf("Hash: %s\n", d57.Hash)
fmt.Printf("Expected: 731d9efe96031900ba2a36667f4718d0\n")

// Let's also print what the token array looks like
lexer := NewLexer(sql)
store := newTokenStore(MySQL57)
reducer := newReducer(store)
handler := newTokenHandler(lexer, store, reducer)
handler.processAll()

fmt.Printf("\nToken array hex: %s\n", hex.EncodeToString(store.tokenArray))
fmt.Printf("Token array len: %d bytes\n", len(store.tokenArray))

// Verify MD5
hash := md5.Sum(store.tokenArray)
fmt.Printf("MD5 of token array: %s\n", hex.EncodeToString(hash[:]))

// Print each 2-byte token
fmt.Printf("\nTokens (little-endian 16-bit):\n")
i := 0
for i < len(store.tokenArray) {
if i+1 < len(store.tokenArray) {
tokVal := int(store.tokenArray[i]) | (int(store.tokenArray[i+1]) << 8)
fmt.Printf("  [%d-%d]: 0x%04x (%d)\n", i, i+1, tokVal, tokVal)

// Check if this is an identifier (look for TOK_IDENT = 939)
if tokVal == 939 && i+3 < len(store.tokenArray) {
// Next 2 bytes are length
strLen := int(store.tokenArray[i+2]) | (int(store.tokenArray[i+3]) << 8)
if i+4+strLen <= len(store.tokenArray) {
str := string(store.tokenArray[i+4 : i+4+strLen])
fmt.Printf("    -> Identifier length=%d: %q\n", strLen, str)
i += 4 + strLen
continue
}
}
}
i += 2
}
}
