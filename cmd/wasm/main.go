//go:build js && wasm

package main

import (
	"encoding/json"
	"math"
	"syscall/js"

	digest "github.com/rashiq/mysql-digest"
)

func computeDigest(this js.Value, args []js.Value) any {
	if len(args) != 2 || args[0].Type() != js.TypeString || args[1].Type() != js.TypeNumber {
		return `{"error":"expected SQL text and a numeric version"}`
	}

	sql := args[0].String()
	version := args[1].Float()
	if version != math.Trunc(version) || version < 0 || version > float64(digest.MySQL57) {
		return `{"error":"unsupported MySQL version"}`
	}

	result, err := digest.Compute(sql, digest.Options{
		Version: digest.MySQLVersion(version),
	})

	if err != nil {
		b, _ := json.Marshal(map[string]string{"error": err.Error()})
		return string(b)
	}

	b, _ := json.Marshal(map[string]string{
		"text": result.Text,
		"hash": result.Hash,
	})
	return string(b)
}

func main() {
	js.Global().Set("computeDigest", js.FuncOf(computeDigest))
	select {}
}
