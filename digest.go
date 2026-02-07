// Package digest lets you build a query digest just like MySQL does.
//
// It normalizes SQL queries by replacing literal values with placeholders
// and computes a hash digest compatible with MySQL's STATEMENT_DIGEST()
// and STATEMENT_DIGEST_TEXT() functions.
package digest

import (
	"github.com/rashiq/mysql-digest/internal"
)

type Digest struct {
	Hash string
	Text string
}

type MySQLVersion = internal.MySQLVersion

const (
	MySQL80 = internal.MySQL80
	MySQL84 = internal.MySQL84
	MySQL57 = internal.MySQL57
)

type SQLMode = internal.SQLMode

const (
	MODE_NO_BACKSLASH_ESCAPES = internal.MODE_NO_BACKSLASH_ESCAPES
	MODE_ANSI_QUOTES          = internal.MODE_ANSI_QUOTES
)

type Options struct {
	SQLMode   SQLMode
	MaxLength int
	Version   MySQLVersion
}

type Digester struct {
	opts Options
}

func NewDigester(opts ...Options) *Digester {
	var o Options
	if len(opts) > 0 {
		o = opts[0]
	}
	return &Digester{opts: o}
}

func (d *Digester) Digest(sql string) (Digest, error) {
	return compute(sql, d.opts)
}

func Compute(sql string, opts ...Options) (Digest, error) {
	var opt Options
	if len(opts) > 0 {
		opt = opts[0]
	}
	return compute(sql, opt)
}

func compute(sql string, opt Options) (Digest, error) {
	lexer := internal.NewLexer(sql)
	lexer.SetSQLMode(opt.SQLMode)
	lexer.SetDigestVersion(opt.Version)

	store := internal.NewTokenStore(opt.Version)
	reducer := internal.NewReducer(store)
	handler := internal.NewTokenHandler(lexer, store, reducer)

	err := handler.ProcessAll()

	return Digest{
		Hash: store.ComputeHash(),
		Text: store.BuildText(opt.MaxLength),
	}, err
}
