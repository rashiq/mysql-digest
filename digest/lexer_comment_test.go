package digest

import "testing"

// Phase 8: MY_LEX_COMMENT, MY_LEX_LONG_COMMENT, MY_LEX_END_LONG_COMMENT tests
// Tests for single-line and multi-line comments

// ============================================================================
// MY_LEX_COMMENT: Single-line comments (-- and #)
// ============================================================================

func TestLexer_COMMENT_DoubleDash(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		firstType int
		firstText string
	}{
		{"comment then SELECT", "-- comment\nSELECT", SELECT_SYM, "SELECT"},
		{"comment with space", "--  spaces  \nSELECT", SELECT_SYM, "SELECT"},
		{"comment at EOF", "-- comment", END_OF_INPUT, ""},
		{"empty comment", "--\nSELECT", SELECT_SYM, "SELECT"},
		{"comment with tab", "-- \tcomment\nSELECT", SELECT_SYM, "SELECT"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.Lex()
			if tok.Type != tt.firstType {
				t.Errorf("expected type %d, got %d", tt.firstType, tok.Type)
			}
			if tt.firstText != "" {
				if got := l.MustTokenText(tok); got != tt.firstText {
					t.Errorf("expected text %q, got %q", tt.firstText, got)
				}
			}
		})
	}
}

func TestLexer_COMMENT_DoubleDashNoSpace(t *testing.T) {
	// -- without space after is NOT a comment, it's two minus signs
	l := NewLexer("--x")

	tok := l.Lex()
	if tok.Type != int('-') {
		t.Errorf("expected '-' (%d), got %d", int('-'), tok.Type)
	}

	tok = l.Lex()
	if tok.Type != int('-') {
		t.Errorf("expected '-' (%d), got %d", int('-'), tok.Type)
	}

	tok = l.Lex()
	if tok.Type != IDENT {
		t.Errorf("expected IDENT, got %d", tok.Type)
	}
	if got := l.MustTokenText(tok); got != "x" {
		t.Errorf("expected 'x', got %q", got)
	}
}

func TestLexer_COMMENT_Hash(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		firstType int
		firstText string
	}{
		{"hash comment then SELECT", "# comment\nSELECT", SELECT_SYM, "SELECT"},
		{"hash with content", "#this is a comment\nSELECT", SELECT_SYM, "SELECT"},
		{"hash at EOF", "# comment", END_OF_INPUT, ""},
		{"empty hash comment", "#\nSELECT", SELECT_SYM, "SELECT"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.Lex()
			if tok.Type != tt.firstType {
				t.Errorf("expected type %d, got %d", tt.firstType, tok.Type)
			}
			if tt.firstText != "" {
				if got := l.MustTokenText(tok); got != tt.firstText {
					t.Errorf("expected text %q, got %q", tt.firstText, got)
				}
			}
		})
	}
}

// ============================================================================
// MY_LEX_LONG_COMMENT: Multi-line C-style comments /* */
// ============================================================================

func TestLexer_COMMENT_CStyle(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		firstType int
		firstText string
	}{
		{"simple comment", "/* comment */ SELECT", SELECT_SYM, "SELECT"},
		{"empty comment", "/**/ SELECT", SELECT_SYM, "SELECT"},
		{"multiline comment", "/* line1\nline2 */ SELECT", SELECT_SYM, "SELECT"},
		{"comment with stars", "/* * * * */ SELECT", SELECT_SYM, "SELECT"},
		{"comment at EOF", "/* comment */", END_OF_INPUT, ""},
		{"multiple comments", "/* c1 */ /* c2 */ SELECT", SELECT_SYM, "SELECT"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.Lex()
			if tok.Type != tt.firstType {
				t.Errorf("expected type %d, got %d", tt.firstType, tok.Type)
			}
			if tt.firstText != "" {
				if got := l.MustTokenText(tok); got != tt.firstText {
					t.Errorf("expected text %q, got %q", tt.firstText, got)
				}
			}
		})
	}
}

func TestLexer_COMMENT_CStyle_Unclosed(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"no closing", "/* unclosed"},
		{"partial close", "/* comment *"},
		{"empty unclosed", "/*"},
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

func TestLexer_COMMENT_SlashNotComment(t *testing.T) {
	// Single '/' is division, not start of comment
	l := NewLexer("a / b")

	tok := l.Lex()
	if tok.Type != IDENT {
		t.Errorf("expected IDENT, got %d", tok.Type)
	}

	tok = l.Lex()
	if tok.Type != int('/') {
		t.Errorf("expected '/' (%d), got %d", int('/'), tok.Type)
	}

	tok = l.Lex()
	if tok.Type != IDENT {
		t.Errorf("expected IDENT, got %d", tok.Type)
	}
}

// ============================================================================
// Version comments /*!50000 ... */
// ============================================================================

func TestLexer_COMMENT_Version_Execute(t *testing.T) {
	// Version comments with version <= current should execute content
	tests := []struct {
		name      string
		input     string
		firstType int
		firstText string
	}{
		{"version 50000", "/*!50000 SELECT */ 1", SELECT_SYM, "SELECT"},
		{"version 80000", "/*!80000 SELECT */ 1", SELECT_SYM, "SELECT"},
		{"version 80400", "/*!80400 SELECT */ 1", SELECT_SYM, "SELECT"},
		{"version 32302", "/*!32302 SELECT */ 1", SELECT_SYM, "SELECT"},
		{"6-digit version", "/*!080000 SELECT */ 1", SELECT_SYM, "SELECT"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.Lex()
			if tok.Type != tt.firstType {
				t.Errorf("expected type %d, got %d", tt.firstType, tok.Type)
			}
			if got := l.MustTokenText(tok); got != tt.firstText {
				t.Errorf("expected text %q, got %q", tt.firstText, got)
			}
		})
	}
}

func TestLexer_COMMENT_Version_Skip(t *testing.T) {
	// Version comments with version > current should be skipped
	tests := []struct {
		name      string
		input     string
		firstType int
		firstText string
	}{
		{"future version", "/*!99999 SELECT */ 1", NUM, "1"},
		{"very future", "/*!999999 SELECT */ 1", NUM, "1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.Lex()
			if tok.Type != tt.firstType {
				t.Errorf("expected type %d, got %d", tt.firstType, tok.Type)
			}
			if got := l.MustTokenText(tok); got != tt.firstText {
				t.Errorf("expected text %q, got %q", tt.firstText, got)
			}
		})
	}
}

func TestLexer_COMMENT_Version_NoVersion(t *testing.T) {
	// /*! without version number - always execute
	tests := []struct {
		name      string
		input     string
		firstType int
		firstText string
	}{
		{"no version", "/*! SELECT */ 1", SELECT_SYM, "SELECT"},
		{"space after bang", "/*! SELECT */ 1", SELECT_SYM, "SELECT"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.Lex()
			if tok.Type != tt.firstType {
				t.Errorf("expected type %d, got %d", tt.firstType, tok.Type)
			}
			if got := l.MustTokenText(tok); got != tt.firstText {
				t.Errorf("expected text %q, got %q", tt.firstText, got)
			}
		})
	}
}

// ============================================================================
// MY_LEX_END_LONG_COMMENT: Asterisk handling
// ============================================================================

func TestLexer_COMMENT_Asterisk(t *testing.T) {
	// Outside of comments, '*' is just multiplication
	l := NewLexer("a * b")

	tok := l.Lex()
	if tok.Type != IDENT {
		t.Errorf("expected IDENT, got %d", tok.Type)
	}

	tok = l.Lex()
	if tok.Type != int('*') {
		t.Errorf("expected '*' (%d), got %d", int('*'), tok.Type)
	}

	tok = l.Lex()
	if tok.Type != IDENT {
		t.Errorf("expected IDENT, got %d", tok.Type)
	}
}

func TestLexer_COMMENT_SelectStar(t *testing.T) {
	// SELECT * FROM - common pattern
	l := NewLexer("SELECT * FROM t")

	expected := []struct {
		typ  int
		text string
	}{
		{SELECT_SYM, "SELECT"},
		{int('*'), "*"},
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

// ============================================================================
// Context tests - comments in typical SQL
// ============================================================================

func TestLexer_COMMENT_InContext(t *testing.T) {
	l := NewLexer("SELECT /* comment */ a FROM /* another */ t")

	expected := []struct {
		typ  int
		text string
	}{
		{SELECT_SYM, "SELECT"},
		{IDENT, "a"},
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

func TestLexer_COMMENT_MixedTypes(t *testing.T) {
	// Mix of comment types
	l := NewLexer("-- line comment\nSELECT /* block */ 1 # end")

	tok := l.Lex()
	if tok.Type != SELECT_SYM {
		t.Errorf("expected SELECT_SYM, got %d", tok.Type)
	}

	tok = l.Lex()
	if tok.Type != NUM {
		t.Errorf("expected NUM, got %d", tok.Type)
	}

	tok = l.Lex()
	if tok.Type != END_OF_INPUT {
		t.Errorf("expected END_OF_INPUT, got %d", tok.Type)
	}
}

// ============================================================================
// Version comment token sequence tests
// ============================================================================

func TestLexer_COMMENT_Version_FullTokenSequence(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []struct {
			typ  int
			text string
		}
	}{
		{
			name:  "versioned comment consumes closing",
			input: "/*!50000 SELECT */ 1",
			expected: []struct {
				typ  int
				text string
			}{
				{SELECT_SYM, "SELECT"},
				{NUM, "1"},
				{END_OF_INPUT, ""},
			},
		},
		{
			name:  "no-version comment consumes closing",
			input: "/*! SELECT */ 1",
			expected: []struct {
				typ  int
				text string
			}{
				{SELECT_SYM, "SELECT"},
				{NUM, "1"},
				{END_OF_INPUT, ""},
			},
		},
		{
			name:  "multiplication inside version comment",
			input: "/*!50000 SELECT 2 * 3 */ AS result",
			expected: []struct {
				typ  int
				text string
			}{
				{SELECT_SYM, "SELECT"},
				{NUM, "2"},
				{int('*'), "*"},
				{NUM, "3"},
				{AS, "AS"},
				{IDENT, "result"},
				{END_OF_INPUT, ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			for i, exp := range tt.expected {
				tok := l.Lex()
				if tok.Type != exp.typ {
					t.Errorf("token %d: expected type %d, got %d", i, exp.typ, tok.Type)
				}
				got := ""
				if tok.Type != END_OF_INPUT {
					got = l.MustTokenText(tok)
				}
				if got != exp.text {
					t.Errorf("token %d: expected text %q, got %q", i, exp.text, got)
				}
			}
		})
	}
}
