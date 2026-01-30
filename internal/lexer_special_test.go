package internal

import (
	"testing"
)

func TestLexer_Special_Semicolon(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType int
	}{
		{"semicolon alone", ";", int(';')},
		{"semicolon after statement", "SELECT 1;", int(';')}, // third token
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)

			// For "SELECT 1;", we need to skip to the semicolon
			if tt.input == "SELECT 1;" {
				l.Lex() // SELECT
				l.Lex() // 1
			}

			tok := l.Lex()
			if tok.Type != tt.wantType {
				t.Errorf("got type %d (%s), want %d",
					tok.Type, TokenString(tok.Type), tt.wantType)
			}
		})
	}
}

func TestLexer_Special_EOF(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty input", ""},
		{"whitespace only", "   "},
		{"after statement", "SELECT"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)

			// Skip tokens until we get END_OF_INPUT
			var tok Token
			for i := 0; i < 10; i++ {
				tok = l.Lex()
				if tok.Type == END_OF_INPUT {
					break
				}
			}

			if tok.Type != END_OF_INPUT {
				t.Errorf("expected END_OF_INPUT, got %d (%s)",
					tok.Type, TokenString(tok.Type))
			}
		})
	}
}

func TestLexer_Special_RepeatedEOF(t *testing.T) {
	l := NewLexer("SELECT")

	// First call: SELECT
	tok1 := l.Lex()
	if tok1.Type != SELECT_SYM {
		t.Fatalf("expected SELECT_SYM, got %d", tok1.Type)
	}

	// Second call: END_OF_INPUT
	tok2 := l.Lex()
	if tok2.Type != END_OF_INPUT {
		t.Fatalf("expected END_OF_INPUT, got %d", tok2.Type)
	}

	// Third call: should still return END_OF_INPUT
	tok3 := l.Lex()
	if tok3.Type != END_OF_INPUT {
		t.Errorf("repeated EOF: expected END_OF_INPUT, got %d", tok3.Type)
	}

	// Fourth call: should still return END_OF_INPUT
	tok4 := l.Lex()
	if tok4.Type != END_OF_INPUT {
		t.Errorf("repeated EOF: expected END_OF_INPUT, got %d", tok4.Type)
	}
}

func TestLexer_Special_RealOrPoint(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantTypes []int
		wantTexts []string
	}{
		{
			"dot five",
			".5",
			[]int{DECIMAL_NUM, END_OF_INPUT},
			[]string{".5", ""},
		},
		{
			"dot five with more digits",
			".123456",
			[]int{DECIMAL_NUM, END_OF_INPUT},
			[]string{".123456", ""},
		},
		{
			"dot with exponent",
			".5e10",
			[]int{FLOAT_NUM, END_OF_INPUT},
			[]string{".5e10", ""},
		},
		{
			"dot ident",
			".foo",
			[]int{int('.'), IDENT, END_OF_INPUT},
			[]string{".", "foo", ""},
		},
		{
			"dot alone",
			".",
			[]int{int('.'), END_OF_INPUT},
			[]string{".", ""},
		},
		{
			"qualified identifier",
			"a.b",
			[]int{IDENT, int('.'), IDENT, END_OF_INPUT},
			[]string{"a", ".", "b", ""},
		},
		{
			"three part identifier",
			"a.b.c",
			[]int{IDENT, int('.'), IDENT, int('.'), IDENT, END_OF_INPUT},
			[]string{"a", ".", "b", ".", "c", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)

			for i, wantType := range tt.wantTypes {
				tok := l.Lex()

				if tok.Type != wantType {
					t.Errorf("token %d: got type %d (%s), want %d (%s)",
						i, tok.Type, TokenString(tok.Type),
						wantType, TokenString(wantType))
				}

				if i < len(tt.wantTexts) && wantType != END_OF_INPUT {
					gotText := tt.input[tok.Start:tok.End]
					if gotText != tt.wantTexts[i] {
						t.Errorf("token %d: got text %q, want %q",
							i, gotText, tt.wantTexts[i])
					}
				}
			}
		})
	}
}

func TestLexer_Special_MultipleStatements(t *testing.T) {
	input := "SELECT 1; SELECT 2"
	wantTypes := []int{SELECT_SYM, NUM, int(';'), SELECT_SYM, NUM, END_OF_INPUT}
	wantTexts := []string{"SELECT", "1", ";", "SELECT", "2", ""}

	l := NewLexer(input)

	for i, wantType := range wantTypes {
		tok := l.Lex()

		if tok.Type != wantType {
			t.Errorf("token %d: got type %d (%s), want %d (%s)",
				i, tok.Type, TokenString(tok.Type),
				wantType, TokenString(wantType))
		}

		if wantType != END_OF_INPUT {
			gotText := input[tok.Start:tok.End]
			if gotText != wantTexts[i] {
				t.Errorf("token %d: got text %q, want %q",
					i, gotText, wantTexts[i])
			}
		}
	}
}
