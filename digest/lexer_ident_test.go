package digest

import (
	"testing"
)

// ============================================================================
// MY_LEX_IDENT_OR_NCHAR tests
// ============================================================================

// TestLexer_IDENT_OR_NCHAR_NcharString tests N'string' returns NCHAR_STRING
func TestLexer_IDENT_OR_NCHAR_NcharString(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"uppercase N", "N'hello'"},
		{"lowercase n", "n'hello'"},
		{"empty string", "N''"},
		{"with escape", "N'it''s'"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(tc.input)
			tok := l.Lex()

			if tok.Type != NCHAR_STRING {
				t.Errorf("input %q: expected NCHAR_STRING (%d), got %d",
					tc.input, NCHAR_STRING, tok.Type)
			}
		})
	}
}

// TestLexer_IDENT_OR_NCHAR_NotNchar tests N not followed by ' is IDENT
func TestLexer_IDENT_OR_NCHAR_NotNchar(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		wantType int
	}{
		{"N alone", "N", IDENT},
		{"Nident", "Nident", IDENT},
		{"NULL keyword", "NULL", NULL_SYM},
		{"NOT keyword", "NOT", NOT_SYM},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(tc.input)
			tok := l.Lex()

			if tok.Type != tc.wantType {
				t.Errorf("input %q: expected %d, got %d",
					tc.input, tc.wantType, tok.Type)
			}
		})
	}
}

// ============================================================================
// MY_LEX_IDENT_OR_HEX tests
// ============================================================================

// TestLexer_IDENT_OR_HEX_HexString tests X'hex' transitions to HEX_NUMBER
func TestLexer_IDENT_OR_HEX_HexString(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"uppercase X", "X'1A2B'"},
		{"lowercase x", "x'1a2b'"},
		{"empty hex", "X''"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(tc.input)
			tok := l.Lex()

			if tok.Type != HEX_NUM {
				t.Errorf("input %q: expected HEX_NUM (%d), got %d",
					tc.input, HEX_NUM, tok.Type)
			}
		})
	}
}

// TestLexer_IDENT_OR_HEX_NotHex tests X not followed by ' is IDENT
func TestLexer_IDENT_OR_HEX_NotHex(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"X alone", "X"},
		{"Xident", "Xident"},
		{"XOR keyword", "XOR"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(tc.input)
			tok := l.Lex()

			// Should be IDENT or a keyword, not HEX_NUM
			if tok.Type == HEX_NUM {
				t.Errorf("input %q: should not be HEX_NUM", tc.input)
			}
		})
	}
}

// ============================================================================
// MY_LEX_IDENT_OR_BIN tests
// ============================================================================

// TestLexer_IDENT_OR_BIN_BinString tests B'bin' transitions to BIN_NUMBER
func TestLexer_IDENT_OR_BIN_BinString(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"uppercase B", "B'1010'"},
		{"lowercase b", "b'1010'"},
		{"empty bin", "B''"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(tc.input)
			tok := l.Lex()

			if tok.Type != BIN_NUM {
				t.Errorf("input %q: expected BIN_NUM (%d), got %d",
					tc.input, BIN_NUM, tok.Type)
			}
		})
	}
}

// TestLexer_IDENT_OR_BIN_NotBin tests B not followed by ' is IDENT
func TestLexer_IDENT_OR_BIN_NotBin(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"B alone", "B"},
		{"Bident", "Bident"},
		{"BY keyword", "BY"},
		{"BETWEEN keyword", "BETWEEN"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(tc.input)
			tok := l.Lex()

			// Should be IDENT or a keyword, not BIN_NUM
			if tok.Type == BIN_NUM {
				t.Errorf("input %q: should not be BIN_NUM", tc.input)
			}
		})
	}
}

// ============================================================================
// MY_LEX_IDENT tests
// ============================================================================

// TestLexer_IDENT_Simple tests simple identifiers
func TestLexer_IDENT_Simple(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"lowercase", "foo"},
		{"uppercase", "FOO"},
		{"mixed case", "FooBar"},
		{"with digits", "foo123"},
		{"with underscore", "foo_bar"},
		{"underscore prefix", "_foo"},
		{"all underscore", "___"},
		{"single char", "a"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(tc.input)
			tok := l.Lex()

			if tok.Type != IDENT && tok.Type != IDENT_QUOTED {
				t.Errorf("input %q: expected IDENT (%d) or IDENT_QUOTED (%d), got %d",
					tc.input, IDENT, IDENT_QUOTED, tok.Type)
			}
			text := l.TokenText(tok)
			if text != tc.input {
				t.Errorf("input %q: expected text %q, got %q",
					tc.input, tc.input, text)
			}
		})
	}
}

// TestLexer_IDENT_Keywords tests that keywords are recognized
func TestLexer_IDENT_Keywords(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		wantType int
	}{
		{"SELECT uppercase", "SELECT", SELECT_SYM},
		{"select lowercase", "select", SELECT_SYM},
		{"Select mixed", "Select", SELECT_SYM},
		{"FROM", "FROM", FROM},
		{"WHERE", "WHERE", WHERE},
		{"AND", "AND", AND_SYM},
		{"OR", "OR", OR_SYM},
		{"INSERT", "INSERT", INSERT_SYM},
		{"UPDATE", "UPDATE", UPDATE_SYM},
		{"DELETE", "DELETE", DELETE_SYM},
		{"CREATE", "CREATE", CREATE},
		{"TABLE", "TABLE", TABLE_SYM},
		{"NULL", "NULL", NULL_SYM},
		{"TRUE", "TRUE", TRUE_SYM},
		{"FALSE", "FALSE", FALSE_SYM},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(tc.input)
			tok := l.Lex()

			if tok.Type != tc.wantType {
				t.Errorf("input %q: expected %d, got %d",
					tc.input, tc.wantType, tok.Type)
			}
		})
	}
}

// TestLexer_IDENT_FollowedByDot tests identifier followed by dot sets IDENT_SEP
func TestLexer_IDENT_FollowedByDot(t *testing.T) {
	l := NewLexer("a.b")

	// First token: 'a' as IDENT
	tok1 := l.Lex()
	if tok1.Type != IDENT && tok1.Type != IDENT_QUOTED {
		t.Errorf("first token: expected IDENT, got %d", tok1.Type)
	}
	text1 := l.TokenText(tok1)
	if text1 != "a" {
		t.Errorf("first token: expected text 'a', got %q", text1)
	}

	// Second token: '.'
	tok2 := l.Lex()
	if tok2.Type != int('.') {
		t.Errorf("second token: expected '.' (%d), got %d", int('.'), tok2.Type)
	}

	// Third token: 'b' as IDENT
	tok3 := l.Lex()
	if tok3.Type != IDENT && tok3.Type != IDENT_QUOTED {
		t.Errorf("third token: expected IDENT, got %d", tok3.Type)
	}
}

// TestLexer_IDENT_WithDollar tests identifiers containing $
func TestLexer_IDENT_WithDollar(t *testing.T) {
	// Note: $ as first char goes to IDENT_OR_DOLLAR_QUOTED_TEXT
	// But $ in middle of ident should work
	l := NewLexer("foo$bar")
	tok := l.Lex()

	if tok.Type != IDENT && tok.Type != IDENT_QUOTED {
		t.Errorf("expected IDENT, got %d", tok.Type)
	}
	text := l.TokenText(tok)
	if text != "foo$bar" {
		t.Errorf("expected 'foo$bar', got %q", text)
	}
}

// TestLexer_IDENT_MultipleIdents tests sequence of identifiers
func TestLexer_IDENT_MultipleIdents(t *testing.T) {
	l := NewLexer("foo bar baz")

	expected := []string{"foo", "bar", "baz"}
	for i, exp := range expected {
		tok := l.Lex()
		if tok.Type != IDENT && tok.Type != IDENT_QUOTED {
			t.Errorf("token %d: expected IDENT, got %d", i, tok.Type)
		}
		text := l.TokenText(tok)
		if text != exp {
			t.Errorf("token %d: expected %q, got %q", i, exp, text)
		}
	}
}

// TestLexer_IDENT_IdentThenOperator tests identifier followed by operator
func TestLexer_IDENT_IdentThenOperator(t *testing.T) {
	l := NewLexer("foo+bar")

	// First: foo
	tok1 := l.Lex()
	if tok1.Type != IDENT && tok1.Type != IDENT_QUOTED {
		t.Errorf("first token: expected IDENT, got %d", tok1.Type)
	}

	// Second: +
	tok2 := l.Lex()
	if tok2.Type != int('+') {
		t.Errorf("second token: expected '+', got %d", tok2.Type)
	}

	// Third: bar
	tok3 := l.Lex()
	if tok3.Type != IDENT && tok3.Type != IDENT_QUOTED {
		t.Errorf("third token: expected IDENT, got %d", tok3.Type)
	}
}

// TestLexer_IDENT_FunctionCall tests identifier followed by parenthesis
func TestLexer_IDENT_FunctionCall(t *testing.T) {
	l := NewLexer("COUNT(")

	tok := l.Lex()
	// COUNT followed by ( should be recognized as COUNT_SYM (function keyword)
	if tok.Type != COUNT_SYM {
		t.Errorf("expected COUNT_SYM (%d), got %d", COUNT_SYM, tok.Type)
	}
}

// TestLexer_IDENT_NotFunctionCall tests keyword not followed by paren
func TestLexer_IDENT_NotFunctionCall(t *testing.T) {
	// Some keywords are only keywords when followed by (
	// For now, test that regular keywords work without (
	l := NewLexer("SELECT foo")

	tok := l.Lex()
	if tok.Type != SELECT_SYM {
		t.Errorf("expected SELECT_SYM (%d), got %d", SELECT_SYM, tok.Type)
	}
}
