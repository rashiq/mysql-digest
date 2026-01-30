package digest

type Digest struct {
	Hash string
	Text string
}

type MySQLVersion int

const (
	MySQL80 MySQLVersion = iota
	MySQL84
	MySQL57
)

type Options struct {
	SQLMode   SQLMode      // ANSI_QUOTES, NO_BACKSLASH_ESCAPES
	MaxLength int          // 0 = unlimited, otherwise truncates with "..."
	Version   MySQLVersion // MySQL 5.7 uses MD5, MySQL 8+ use SHA-256
}

func Normalize(sql string, opts ...Options) Digest {
	opt := Options{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	lexer := NewLexer(sql)
	lexer.SetSQLMode(opt.SQLMode)
	lexer.SetDigestVersion(opt.Version)

	store := newTokenStore(opt.Version)
	reducer := newReducer(store)
	handler := newTokenHandler(lexer, store, reducer)

	handler.processAll()

	return Digest{
		Hash: store.computeHash(),
		Text: store.buildText(opt.MaxLength),
	}
}

func Compute(sql string) Digest {
	return Normalize(sql)
}
