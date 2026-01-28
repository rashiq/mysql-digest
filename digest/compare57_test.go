package digest

import (
"fmt"
"testing"
)

func TestCompare57Hashes(t *testing.T) {
tests := []struct {
sql      string
expected string
expectedText string
}{
{"SELECT 1", "3d4fc22e33e10d7235eced3c75a84c2c", "SELECT ?"},
{"SELECT 1 FROM DUAL", "6147128d66c29360df7f61f7b29450e0", "SELECT ? FROM DUAL"},
{"SELECT * FROM users", "a70c32b06c430b191593fe98cf4069a1", "SELECT * FROM `users`"},
{"SELECT * FROM users WHERE id = 42", "731d9efe96031900ba2a36667f4718d0", "SELECT * FROM `users` WHERE `id` = ?"},
}

for _, tt := range tests {
d := Normalize(tt.sql, Options{Version: MySQL57})
status := "PASS"
if d.Hash != tt.expected {
status = "FAIL"
}
fmt.Printf("%s: %s\n", status, tt.sql)
fmt.Printf("  Expected: %s\n", tt.expected)
fmt.Printf("  Got:      %s\n", d.Hash)
fmt.Printf("  Text:     %s\n", d.Text)
fmt.Printf("  Expected: %s\n", tt.expectedText)
fmt.Println()
}
}
