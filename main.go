package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"mysql-digest/digest"
)

func main() {
	// Define flags
	sqlFlag := flag.String("sql", "", "SQL statement to compute digest for")
	fileFlag := flag.String("file", "", "File containing SQL statement(s)")
	stdinFlag := flag.Bool("stdin", false, "Read SQL from stdin")
	jsonFlag := flag.Bool("json", false, "Output in JSON format")
	textOnlyFlag := flag.Bool("text-only", false, "Output only the normalized text")
	hashOnlyFlag := flag.Bool("hash-only", false, "Output only the digest hash")
	debugFlag := flag.Bool("debug", false, "Print lexer tokens for debugging")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: mysql-digest [options]\n\n")
		fmt.Fprintf(os.Stderr, "Compute MySQL query digest from SQL statements.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  mysql-digest -sql \"SELECT * FROM users WHERE id = 123\"\n")
		fmt.Fprintf(os.Stderr, "  mysql-digest -file query.sql\n")
		fmt.Fprintf(os.Stderr, "  echo \"SELECT 1\" | mysql-digest -stdin\n")
		fmt.Fprintf(os.Stderr, "  mysql-digest -sql \"SELECT 1\" -json\n")
	}

	flag.Parse()

	// Determine input source
	var sql string
	var err error

	switch {
	case *sqlFlag != "":
		sql = *sqlFlag
	case *fileFlag != "":
		sql, err = readFile(*fileFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
	case *stdinFlag:
		sql, err = readStdin()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}
	default:
		// Check if there are positional arguments
		if flag.NArg() > 0 {
			sql = strings.Join(flag.Args(), " ")
		} else {
			// No input provided, try reading from stdin if it's a pipe
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				sql, err = readStdin()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
					os.Exit(1)
				}
			} else {
				flag.Usage()
				os.Exit(1)
			}
		}
	}

	sql = strings.TrimSpace(sql)
	if sql == "" {
		fmt.Fprintf(os.Stderr, "Error: empty SQL input\n")
		os.Exit(1)
	}

	// Compute digest
	result := digest.Compute(sql)

	// Debug output - print tokens
	if *debugFlag {
		fmt.Println("Lexer Tokens (from scanner):")
		l := digest.NewLexer(sql)
		for {
			tok := l.Lex()
			if tok.Type == digest.END_OF_INPUT {
				break
			}
			rawText := sql[tok.Start:tok.End]
			normalized := digest.TokenString(tok.Type)
			if normalized != rawText && normalized != "" && normalized != "(unknown)" {
				fmt.Printf("  {id=%d norm=%q raw=%q}\n", tok.Type, normalized, rawText)
			} else {
				fmt.Printf("  {id=%d raw=%q}\n", tok.Type, rawText)
			}
		}
		fmt.Println()
	}

	// Output result
	switch {
	case *jsonFlag:
		output := struct {
			Hash string `json:"digest"`
			Text string `json:"digest_text"`
		}{
			Hash: result.Hash,
			Text: result.Text,
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(output)
	case *textOnlyFlag:
		fmt.Println(result.Text)
	case *hashOnlyFlag:
		fmt.Println(result.Hash)
	default:
		fmt.Printf("DIGEST: %s\n", result.Hash)
		fmt.Printf("DIGEST_TEXT: %s\n", result.Text)
	}
}

func readFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func readStdin() (string, error) {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
