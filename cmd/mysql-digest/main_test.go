package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCommand(t *testing.T) {
	path := filepath.Join(t.TempDir(), "query.sql")
	if err := os.WriteFile(path, []byte("SELECT 1"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		args  []string
		input string
		want  string
		fail  bool
	}{
		{args: []string{"SELECT 1", "--text-only"}, want: "SELECT ?\n"},
		{args: []string{"--sql", "SELECT 1", "--text-only"}, want: "SELECT ?\n"},
		{args: []string{"--file", path, "--text-only"}, want: "SELECT ?\n"},
		{args: []string{"--text-only"}, input: "SELECT 1", want: "SELECT ?\n"},
		{args: []string{"SELECT 1", "--json"}, want: `"digest_text": "SELECT ?"`},
		{args: []string{"SELECT 1"}, want: "DIGEST_TEXT: SELECT ?\n"},
		{args: []string{"SELECT 1", "--json", "--hash-only"}, fail: true},
		{args: []string{"--sql", "SELECT 1", "--file", path}, fail: true},
		{args: []string{"SELECT 1", "--sql", "SELECT 2"}, fail: true},
		{args: []string{"--sql", ""}, input: "SELECT 1", fail: true},
		{args: []string{"SELECT 'unfinished"}, fail: true},
		{fail: true},
	} {
		cmd := newCommand()
		cmd.SetArgs(tc.args)
		cmd.SetIn(strings.NewReader(tc.input))
		var out bytes.Buffer
		cmd.SetOut(&out)
		cmd.SetErr(io.Discard)
		err := cmd.Execute()
		if (err != nil) != tc.fail || (!tc.fail && !strings.Contains(out.String(), tc.want)) {
			t.Errorf("%v: output %q, error %v", tc.args, out.String(), err)
		}
	}
}

type failedIO struct{ err error }

func (f failedIO) Read([]byte) (int, error)  { return 0, f.err }
func (f failedIO) Write([]byte) (int, error) { return 0, f.err }

func TestCommandIOErrors(t *testing.T) {
	failure := errors.New("I/O failure")
	for _, args := range [][]string{
		{}, {"SELECT 1"}, {"SELECT 1", "--json"},
		{"SELECT 1", "--text-only"}, {"SELECT 1", "--hash-only"},
	} {
		cmd := newCommand()
		cmd.SetArgs(args)
		cmd.SetIn(failedIO{failure})
		cmd.SetOut(failedIO{failure})
		cmd.SetErr(io.Discard)
		if err := cmd.Execute(); !errors.Is(err, failure) {
			t.Errorf("%v: got %v, want %v", args, err, failure)
		}
	}
	file, err := os.CreateTemp(t.TempDir(), "closed")
	if err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	cmd := newCommand()
	cmd.SetArgs([]string{})
	cmd.SetIn(file)
	cmd.SetErr(io.Discard)
	if err := cmd.Execute(); !errors.Is(err, os.ErrClosed) {
		t.Fatalf("got %v, want a closed file error", err)
	}
}
