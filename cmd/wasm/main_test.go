//go:build js && wasm

package main

import (
	"encoding/json"
	"math"
	"syscall/js"
	"testing"
)

func TestComputeDigest(t *testing.T) {
	for _, tc := range []struct {
		args []any
		fail bool
	}{
		{args: []any{"SELECT 1", 0}},
		{args: []any{"SELECT 1", 3}},
		{args: []any{"SELECT 'unfinished", 0}, fail: true},
		{fail: true},
		{args: []any{nil, 0}, fail: true},
		{args: []any{"SELECT 1", "0"}, fail: true},
		{args: []any{"SELECT 1", 0.5}, fail: true},
		{args: []any{"SELECT 1", math.NaN()}, fail: true},
		{args: []any{"SELECT 1", math.Inf(1)}, fail: true},
		{args: []any{"SELECT 1", -1}, fail: true},
		{args: []any{"SELECT 1", 4}, fail: true},
	} {
		args := make([]js.Value, len(tc.args))
		for i, arg := range tc.args {
			args[i] = js.ValueOf(arg)
		}
		var result map[string]string
		response := computeDigest(js.Undefined(), args).(string)
		if err := json.Unmarshal([]byte(response), &result); err != nil {
			t.Fatal(err)
		}
		if (result["error"] != "") != tc.fail || (!tc.fail && result["text"] != "SELECT ?") {
			t.Errorf("%v: %s", tc.args, response)
		}
	}
}
