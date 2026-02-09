//go:build js && wasm

package main

import (
	"encoding/json"
	"syscall/js"

	digest "github.com/rashiq/mysql-digest"
)

func computeDigest(this js.Value, args []js.Value) any {
	if len(args) < 2 {
		return `{"error":"missing arguments"}`
	}

	sql := args[0].String()
	version := args[1].Int()

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
