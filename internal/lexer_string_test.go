package internal

import "testing"

func TestLexer_STRING_SingleQuote_Valid(t *testing.T) {
	tests := []struct {
		name  string
		input string
		text  string
	}{
		{"simple string", "'hello'", "'hello'"},
		{"empty string", "''", "''"},
		{"with spaces", "'hello world'", "'hello world'"},
		{"with digits", "'abc123'", "'abc123'"},
		{"single char", "'x'", "'x'"},
		{"unicode", "'日本語'", "'日本語'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.Lex()
			if tok.Type != TEXT_STRING {
				t.Errorf("expected TEXT_STRING (%d), got %d", TEXT_STRING, tok.Type)
			}
			if got := l.MustTokenText(tok); got != tt.text {
				t.Errorf("expected text %q, got %q", tt.text, got)
			}
		})
	}
}

func TestLexer_STRING_EscapedQuote(t *testing.T) {
	tests := []struct {
		name  string
		input string
		text  string
	}{
		{"doubled quote", "'it''s'", "'it''s'"},
		{"multiple doubled", "'a''b''c'", "'a''b''c'"},
		{"doubled at start", "'''hello'", "'''hello'"},
		{"doubled at end", "'hello'''", "'hello'''"},
		{"just doubled quotes", "''''", "''''"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.Lex()
			if tok.Type != TEXT_STRING {
				t.Errorf("expected TEXT_STRING (%d), got %d", TEXT_STRING, tok.Type)
			}
			if got := l.MustTokenText(tok); got != tt.text {
				t.Errorf("expected text %q, got %q", tt.text, got)
			}
		})
	}
}

func TestLexer_STRING_BackslashEscape(t *testing.T) {
	// Default mode: backslash is escape character
	tests := []struct {
		name  string
		input string
		text  string
	}{
		{"escaped n", "'hello\\nworld'", "'hello\\nworld'"},
		{"escaped t", "'hello\\tworld'", "'hello\\tworld'"},
		{"escaped quote", "'don\\'t'", "'don\\'t'"},
		{"escaped backslash", "'path\\\\file'", "'path\\\\file'"},
		{"escaped at end", "'test\\\\'", "'test\\\\'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.Lex()
			if tok.Type != TEXT_STRING {
				t.Errorf("expected TEXT_STRING (%d), got %d", TEXT_STRING, tok.Type)
			}
			if got := l.MustTokenText(tok); got != tt.text {
				t.Errorf("expected text %q, got %q", tt.text, got)
			}
		})
	}
}

func TestLexer_STRING_NoBackslashEscapes(t *testing.T) {
	// NO_BACKSLASH_ESCAPES mode: backslash is literal
	tests := []struct {
		name  string
		input string
		text  string
	}{
		{"backslash n literal", "'hello\\nworld'", "'hello\\nworld'"},
		{"backslash quote", "'don\\'t'", "'don\\'"}, // Stops at the quote!
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			l.SetSQLMode(MODE_NO_BACKSLASH_ESCAPES)
			tok := l.Lex()
			if tok.Type != TEXT_STRING {
				t.Errorf("expected TEXT_STRING (%d), got %d", TEXT_STRING, tok.Type)
			}
			if got := l.MustTokenText(tok); got != tt.text {
				t.Errorf("expected text %q, got %q", tt.text, got)
			}
		})
	}
}

func TestLexer_STRING_Unclosed(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"no closing quote", "'hello"},
		{"empty unclosed", "'"},
		{"with content unclosed", "'some text here"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.Lex()
			if tok.Type != ABORT_SYM {
				t.Errorf("expected ABORT_SYM (%d), got %d", ABORT_SYM, tok.Type)
			}
		})
	}
}

func TestLexer_STRING_DoubleQuote_Default(t *testing.T) {
	// Default mode: " is a string delimiter
	tests := []struct {
		name  string
		input string
		text  string
	}{
		{"simple string", "\"hello\"", "\"hello\""},
		{"empty string", "\"\"", "\"\""},
		{"with spaces", "\"hello world\"", "\"hello world\""},
		{"escaped double quote", "\"say \"\"hello\"\"\"", "\"say \"\"hello\"\"\""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.Lex()
			if tok.Type != TEXT_STRING {
				t.Errorf("expected TEXT_STRING (%d), got %d", TEXT_STRING, tok.Type)
			}
			if got := l.MustTokenText(tok); got != tt.text {
				t.Errorf("expected text %q, got %q", tt.text, got)
			}
		})
	}
}

func TestLexer_STRING_DoubleQuote_AnsiQuotes(t *testing.T) {
	// ANSI_QUOTES mode: " is an identifier delimiter
	tests := []struct {
		name  string
		input string
		text  string
	}{
		{"simple ident", "\"column\"", "\"column\""},
		{"with underscore", "\"my_column\"", "\"my_column\""},
		{"reserved word", "\"select\"", "\"select\""},
		{"escaped double quote", "\"col\"\"name\"", "\"col\"\"name\""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			l.SetSQLMode(MODE_ANSI_QUOTES)
			tok := l.Lex()
			if tok.Type != IDENT_QUOTED {
				t.Errorf("expected IDENT_QUOTED (%d), got %d", IDENT_QUOTED, tok.Type)
			}
			if got := l.MustTokenText(tok); got != tt.text {
				t.Errorf("expected text %q, got %q", tt.text, got)
			}
		})
	}
}

func TestLexer_STRING_DoubleQuote_Unclosed(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"no closing quote", "\"hello"},
		{"empty unclosed", "\""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.Lex()
			if tok.Type != ABORT_SYM {
				t.Errorf("expected ABORT_SYM (%d), got %d", ABORT_SYM, tok.Type)
			}
		})
	}
}

func TestLexer_STRING_Backtick_Valid(t *testing.T) {
	tests := []struct {
		name  string
		input string
		text  string
	}{
		{"simple ident", "`column`", "`column`"},
		{"with spaces", "`my column`", "`my column`"},
		{"with special chars", "`col-name`", "`col-name`"},
		{"reserved word", "`select`", "`select`"},
		{"with digits", "`col123`", "`col123`"},
		{"empty", "``", "``"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.Lex()
			if tok.Type != IDENT_QUOTED {
				t.Errorf("expected IDENT_QUOTED (%d), got %d", IDENT_QUOTED, tok.Type)
			}
			if got := l.MustTokenText(tok); got != tt.text {
				t.Errorf("expected text %q, got %q", tt.text, got)
			}
		})
	}
}

func TestLexer_STRING_Backtick_EscapedBacktick(t *testing.T) {
	tests := []struct {
		name  string
		input string
		text  string
	}{
		{"doubled backtick", "`col``name`", "`col``name`"},
		{"multiple doubled", "`a``b``c`", "`a``b``c`"},
		{"at start", "```col`", "```col`"},
		{"at end", "`col```", "`col```"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.Lex()
			if tok.Type != IDENT_QUOTED {
				t.Errorf("expected IDENT_QUOTED (%d), got %d", IDENT_QUOTED, tok.Type)
			}
			if got := l.MustTokenText(tok); got != tt.text {
				t.Errorf("expected text %q, got %q", tt.text, got)
			}
		})
	}
}

func TestLexer_STRING_Backtick_Unclosed(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"no closing backtick", "`column"},
		{"empty unclosed", "`"},
		{"with content unclosed", "`my_column"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.Lex()
			if tok.Type != ABORT_SYM {
				t.Errorf("expected ABORT_SYM (%d), got %d", ABORT_SYM, tok.Type)
			}
		})
	}
}

func TestLexer_STRING_InContext(t *testing.T) {
	l := NewLexer("SELECT 'hello', \"world\", `column` FROM t")

	expected := []struct {
		typ  int
		text string
	}{
		{SELECT_SYM, "SELECT"},
		{TEXT_STRING, "'hello'"},
		{int(','), ","},
		{TEXT_STRING, "\"world\""},
		{int(','), ","},
		{IDENT_QUOTED, "`column`"},
		{FROM, "FROM"},
		{IDENT, "t"},
	}

	for i, exp := range expected {
		tok := l.Lex()
		if tok.Type != exp.typ {
			t.Errorf("token %d: expected type %d, got %d", i, exp.typ, tok.Type)
		}
		if got := l.MustTokenText(tok); got != exp.text {
			t.Errorf("token %d: expected text %q, got %q", i, exp.text, got)
		}
	}
}

func TestLexer_STRING_AnsiQuotesContext(t *testing.T) {
	l := NewLexer("SELECT \"column\" FROM t")
	l.SetSQLMode(MODE_ANSI_QUOTES)

	expected := []struct {
		typ  int
		text string
	}{
		{SELECT_SYM, "SELECT"},
		{IDENT_QUOTED, "\"column\""},
		{FROM, "FROM"},
		{IDENT, "t"},
	}

	for i, exp := range expected {
		tok := l.Lex()
		if tok.Type != exp.typ {
			t.Errorf("token %d: expected type %d, got %d", i, exp.typ, tok.Type)
		}
		if got := l.MustTokenText(tok); got != exp.text {
			t.Errorf("token %d: expected text %q, got %q", i, exp.text, got)
		}
	}
}

func TestLexer_STRING_Consecutive(t *testing.T) {
	// MySQL allows consecutive string literals (they get concatenated by parser)
	l := NewLexer("'hello' 'world'")

	tok := l.Lex()
	if tok.Type != TEXT_STRING {
		t.Errorf("expected TEXT_STRING, got %d", tok.Type)
	}
	if got := l.MustTokenText(tok); got != "'hello'" {
		t.Errorf("expected 'hello', got %q", got)
	}

	tok = l.Lex()
	if tok.Type != TEXT_STRING {
		t.Errorf("expected TEXT_STRING, got %d", tok.Type)
	}
	if got := l.MustTokenText(tok); got != "'world'" {
		t.Errorf("expected 'world', got %q", got)
	}
}
