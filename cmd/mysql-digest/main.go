package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	digest "github.com/rashiq/mysql-digest"
	"github.com/spf13/cobra"
)

var (
	sqlInput   string
	fileInput  string
	jsonOutput bool
	textOnly   bool
	hashOnly   bool
)

func main() {
	cmd := &cobra.Command{
		Use:   "mysql-digest [sql]",
		Short: "Compute MySQL query digest",
		Long:  "Compute MySQL query digest from SQL statements, matching MySQL's Performance Schema.",
		Example: `  mysql-digest "SELECT * FROM users WHERE id = 123"
  mysql-digest --sql "SELECT * FROM users WHERE id = 123"
  mysql-digest --file query.sql
  echo "SELECT 1" | mysql-digest
  mysql-digest "SELECT 1" --json`,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
		RunE:         run,
	}

	cmd.Flags().StringVar(&sqlInput, "sql", "", "SQL statement to compute digest for")
	cmd.Flags().StringVarP(&fileInput, "file", "f", "", "file containing SQL statement")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "output in JSON format")
	cmd.Flags().BoolVar(&textOnly, "text-only", false, "output only the normalized text")
	cmd.Flags().BoolVar(&hashOnly, "hash-only", false, "output only the digest hash")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	sql, err := getSQL(args)
	if err != nil {
		return err
	}

	result := digest.Compute(sql)

	return output(result)
}

func getSQL(args []string) (string, error) {
	var sql string

	switch {
	case sqlInput != "":
		sql = sqlInput
	case fileInput != "":
		data, err := os.ReadFile(fileInput)
		if err != nil {
			return "", fmt.Errorf("reading file: %w", err)
		}
		sql = string(data)
	case len(args) > 0:
		sql = args[0]
	default:
		if isPipe() {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return "", fmt.Errorf("reading stdin: %w", err)
			}
			sql = string(data)
		}
	}

	sql = strings.TrimSpace(sql)
	if sql == "" {
		return "", fmt.Errorf("no SQL input provided")
	}

	return sql, nil
}

func isPipe() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func output(result digest.Digest) error {
	switch {
	case jsonOutput:
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(map[string]string{
			"digest":      result.Hash,
			"digest_text": result.Text,
		})
	case textOnly:
		fmt.Println(result.Text)
	case hashOnly:
		fmt.Println(result.Hash)
	default:
		fmt.Printf("DIGEST: %s\n", result.Hash)
		fmt.Printf("DIGEST_TEXT: %s\n", result.Text)
	}
	return nil
}
