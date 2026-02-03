// Package digest lets you build a query digest just like MySQL does.
//
// It normalizes SQL queries by replacing literal values with placeholders
// and computes a hash digest compatible with MySQL's STATEMENT_DIGEST()
// and STATEMENT_DIGEST_TEXT() functions.
package digest

import (
	"github.com/rashiq/mysql-digest/internal"
)

// Digest represents a normalized SQL query and its hash.
type Digest struct {
	Hash string // SHA-256 hash (MD5 for MySQL 5.7)
	Text string // Normalized query text
}

// MySQLVersion represents the target MySQL version for digest computation.
type MySQLVersion = internal.MySQLVersion

const (
	MySQL80 = internal.MySQL80
	MySQL84 = internal.MySQL84
	MySQL57 = internal.MySQL57
)

// SQLMode represents MySQL SQL mode flags that affect parsing.
type SQLMode = internal.SQLMode

const (
	// MODE_NO_BACKSLASH_ESCAPES disables backslash as escape character in strings.
	MODE_NO_BACKSLASH_ESCAPES = internal.MODE_NO_BACKSLASH_ESCAPES
	// MODE_ANSI_QUOTES treats " as identifier delimiter instead of string delimiter.
	MODE_ANSI_QUOTES = internal.MODE_ANSI_QUOTES
)

// Options configures digest computation.
type Options struct {
	SQLMode   SQLMode      // SQL mode flags
	MaxLength int          // 0 = unlimited, otherwise truncates with "..."
	Version   MySQLVersion // MySQL version (affects hash algorithm and token handling)
}

// Compute computes a digest with default options, if no options specified.
func Compute(sql string, opts ...Options) Digest {
	opt := Options{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	lexer := internal.NewLexer(sql)
	lexer.SetSQLMode(opt.SQLMode)
	lexer.SetDigestVersion(opt.Version)

	store := internal.NewTokenStore(opt.Version)
	reducer := internal.NewReducer(store)
	handler := internal.NewTokenHandler(lexer, store, reducer)

	handler.ProcessAll()

	return Digest{
		Hash: store.ComputeHash(),
		Text: store.BuildText(opt.MaxLength),
	}
}
