package digest

// Digest represents a computed SQL digest.
type Digest struct {
	// Hash is the SHA-256 hash of the normalized token array (hex-encoded).
	// This matches MySQL's STATEMENT_DIGEST() function output.
	Hash string

	// Text is the human-readable normalized SQL with literals replaced by placeholders.
	// This matches MySQL's STATEMENT_DIGEST_TEXT() function output.
	Text string
}

// Options controls digest computation behavior.
type Options struct {
	// SQLMode affects lexer behavior (ANSI_QUOTES, NO_BACKSLASH_ESCAPES).
	SQLMode SQLMode

	// MaxLength limits the digest text length (0 = unlimited).
	// If exceeded, the text is truncated with "..." appended.
	MaxLength int
}

// Normalize computes the digest of a SQL statement.
//
// It normalizes the SQL by:
//   - Replacing literal values (strings, numbers) with placeholders (?)
//   - Collapsing multiple values in IN(...) to a single placeholder
//   - Collapsing multiple rows in VALUES(...) to a single row with comment
//   - Preserving keywords and identifiers
//   - Normalizing whitespace
//
// Options can be provided to customize behavior. If no options are provided,
// default settings are used.
func Normalize(sql string, opts ...Options) Digest {
	opt := Options{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	lexer := NewLexer(sql)
	lexer.SetSQLMode(opt.SQLMode)

	store := newTokenStore()
	reducer := newReducer(store)
	handler := newTokenHandler(lexer, store, reducer)

	handler.processAll()

	return Digest{
		Hash: store.computeHash(),
		Text: store.buildText(opt.MaxLength),
	}
}

// Compute calculates the digest of a SQL statement with default options.
// This is a convenience wrapper around Normalize.
func Compute(sql string) Digest {
	return Normalize(sql)
}

// ComputeWithOptions calculates digest with custom options.
// Deprecated: Use Normalize(sql, opts) instead.
func ComputeWithOptions(sql string, opts Options) Digest {
	return Normalize(sql, opts)
}
