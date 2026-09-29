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

type options struct {
	sqlInput   string
	fileInput  string
	jsonOutput bool
	textOnly   bool
	hashOnly   bool
}

func main() {
	if err := newCommand().Execute(); err != nil {
		os.Exit(1)
	}
}

func newCommand() *cobra.Command {
	var opts options
	cmd := &cobra.Command{
		Use:          "mysql-digest [sql]",
		Short:        "Compute a MySQL query digest",
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			sql, err := opts.getSQL(cmd, args)
			if err != nil {
				return err
			}
			result, err := digest.Compute(sql)
			if err != nil {
				return fmt.Errorf("compute digest: %w", err)
			}
			return opts.output(cmd.OutOrStdout(), result)
		},
	}
	cmd.Flags().StringVar(&opts.sqlInput, "sql", "", "SQL statement")
	cmd.Flags().StringVarP(&opts.fileInput, "file", "f", "", "SQL file")
	cmd.Flags().BoolVar(&opts.jsonOutput, "json", false, "output JSON")
	cmd.Flags().BoolVar(&opts.textOnly, "text-only", false, "output only the normalized text")
	cmd.Flags().BoolVar(&opts.hashOnly, "hash-only", false, "output only the digest hash")
	cmd.MarkFlagsMutuallyExclusive("sql", "file")
	cmd.MarkFlagsMutuallyExclusive("json", "text-only", "hash-only")
	return cmd
}

func (o options) getSQL(cmd *cobra.Command, args []string) (string, error) {
	if len(args) > 0 && (cmd.Flags().Changed("sql") || cmd.Flags().Changed("file")) {
		return "", fmt.Errorf("use one SQL input source")
	}
	var sql string
	switch {
	case cmd.Flags().Changed("sql"):
		sql = o.sqlInput
	case cmd.Flags().Changed("file"):
		data, err := os.ReadFile(o.fileInput)
		if err != nil {
			return "", fmt.Errorf("read file: %w", err)
		}
		sql = string(data)
	case len(args) > 0:
		sql = args[0]
	default:
		input := cmd.InOrStdin()
		if file, ok := input.(*os.File); ok {
			stat, err := file.Stat()
			if err != nil {
				return "", fmt.Errorf("read stdin: %w", err)
			}
			if stat.Mode()&os.ModeCharDevice != 0 {
				return "", fmt.Errorf("no SQL input provided")
			}
		}
		data, err := io.ReadAll(input)
		if err != nil {
			return "", fmt.Errorf("read stdin: %w", err)
		}
		sql = string(data)
	}
	if strings.TrimSpace(sql) == "" {
		return "", fmt.Errorf("no SQL input provided")
	}
	return sql, nil
}

func (o options) output(w io.Writer, result digest.Digest) error {
	var err error
	switch {
	case o.jsonOutput:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(map[string]string{
			"digest": result.Hash, "digest_text": result.Text,
		})
	case o.textOnly:
		_, err = fmt.Fprintln(w, result.Text)
	case o.hashOnly:
		_, err = fmt.Fprintln(w, result.Hash)
	default:
		_, err = fmt.Fprintf(w, "DIGEST: %s\nDIGEST_TEXT: %s\n", result.Hash, result.Text)
	}
	return err
}
