package digest

import (
	"testing"
)

// TestLexer_OP_Comparison tests comparison operators
func TestLexer_OP_Comparison(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType int
		wantText string
	}{
		// Greater than and greater-or-equal
		{"greater than", ">", GT_SYM, ">"},
		{"greater or equal", ">=", GE, ">="},

		// Less than and variations
		{"less than", "<", LT, "<"},
		{"less or equal", "<=", LE, "<="},
		{"not equal <>", "<>", NE, "<>"},
		{"null-safe equal", "<=>", EQUAL_SYM, "<=>"},

		// Equal
		{"equal", "=", EQ, "="},

		// Not equal (with !)
		{"not equal !=", "!=", NE, "!="},

		// Shift operators
		{"shift left", "<<", SHIFT_LEFT, "<<"},
		{"shift right", ">>", SHIFT_RIGHT, ">>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.Lex()

			if tok.Type != tt.wantType {
				t.Errorf("got type %d (%s), want %d (%s)",
					tok.Type, TokenString(tok.Type),
					tt.wantType, TokenString(tt.wantType))
			}

			gotText := tt.input[tok.Start:tok.End]
			if gotText != tt.wantText {
				t.Errorf("got text %q, want %q", gotText, tt.wantText)
			}
		})
	}
}

// TestLexer_OP_Boolean tests boolean operators && and ||
func TestLexer_OP_Boolean(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType int
		wantText string
	}{
		{"and and", "&&", AND_AND_SYM, "&&"},
		{"or or", "||", OR_OR_SYM, "||"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.Lex()

			if tok.Type != tt.wantType {
				t.Errorf("got type %d (%s), want %d (%s)",
					tok.Type, TokenString(tok.Type),
					tt.wantType, TokenString(tt.wantType))
			}

			gotText := tt.input[tok.Start:tok.End]
			if gotText != tt.wantText {
				t.Errorf("got text %q, want %q", gotText, tt.wantText)
			}
		})
	}
}

// TestLexer_OP_SingleChar tests single character operators
func TestLexer_OP_SingleChar(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType int
	}{
		{"single ampersand", "&", int('&')},
		{"single pipe", "|", int('|')},
		{"colon", ":", int(':')},
		{"semicolon", ";", int(';')},
		{"exclamation alone", "! ", int('!')},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.Lex()

			if tok.Type != tt.wantType {
				t.Errorf("got type %d (%s), want %d",
					tok.Type, TokenString(tok.Type), tt.wantType)
			}
		})
	}
}

// TestLexer_OP_SetVar tests the := operator
func TestLexer_OP_SetVar(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType int
		wantText string
	}{
		{"set var", ":=", SET_VAR, ":="},
		{"colon alone", ": ", int(':'), ":"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.Lex()

			if tok.Type != tt.wantType {
				t.Errorf("got type %d (%s), want %d (%s)",
					tok.Type, TokenString(tok.Type),
					tt.wantType, TokenString(tt.wantType))
			}

			gotText := tt.input[tok.Start:tok.End]
			if gotText != tt.wantText {
				t.Errorf("got text %q, want %q", gotText, tt.wantText)
			}
		})
	}
}

// TestLexer_OP_InContext tests operators in context with other tokens
func TestLexer_OP_InContext(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantTypes []int
		wantTexts []string
	}{
		{
			"a >= b",
			"a >= b",
			[]int{IDENT, GE, IDENT, END_OF_INPUT},
			[]string{"a", ">=", "b", ""},
		},
		{
			"x<>y",
			"x<>y",
			[]int{IDENT, NE, IDENT, END_OF_INPUT},
			[]string{"x", "<>", "y", ""},
		},
		{
			"a<=>b",
			"a<=>b",
			[]int{IDENT, EQUAL_SYM, IDENT, END_OF_INPUT},
			[]string{"a", "<=>", "b", ""},
		},
		{
			"x && y || z",
			"x && y || z",
			[]int{IDENT, AND_AND_SYM, IDENT, OR_OR_SYM, IDENT, END_OF_INPUT},
			[]string{"x", "&&", "y", "||", "z", ""},
		},
		{
			"a<<1",
			"a<<1",
			[]int{IDENT, SHIFT_LEFT, NUM, END_OF_INPUT},
			[]string{"a", "<<", "1", ""},
		},
		{
			"b>>2",
			"b>>2",
			[]int{IDENT, SHIFT_RIGHT, NUM, END_OF_INPUT},
			[]string{"b", ">>", "2", ""},
		},
		// Note: @x tokenization requires Phase 11 (USER_END state) - skipped for now
		// {
		// 	"SET @x := 5",
		// 	"SET @x := 5",
		// 	[]int{SET_SYM, IDENT_QUOTED, SET_VAR, NUM, END_OF_INPUT},
		// 	[]string{"SET", "@x", ":=", "5", ""},
		// },
		{
			"a & b | c",
			"a & b | c",
			[]int{IDENT, int('&'), IDENT, int('|'), IDENT, END_OF_INPUT},
			[]string{"a", "&", "b", "|", "c", ""},
		},
		{
			"statement;",
			"SELECT 1;",
			[]int{SELECT_SYM, NUM, int(';'), END_OF_INPUT},
			[]string{"SELECT", "1", ";", ""},
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

// TestLexer_OP_EdgeCases tests edge cases for operators
func TestLexer_OP_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantTypes []int
	}{
		// Multiple operators in a row
		{">>>=", ">>>=", []int{SHIFT_RIGHT, GE, END_OF_INPUT}},
		{"<<<<", "<<<<", []int{SHIFT_LEFT, SHIFT_LEFT, END_OF_INPUT}},

		// Operators at end of input
		{">EOF", ">", []int{GT_SYM, END_OF_INPUT}},
		{"<EOF", "<", []int{LT, END_OF_INPUT}},
		{"!EOF", "!", []int{int('!'), END_OF_INPUT}},

		// Mixed with identifiers
		{"a>b", "a>b", []int{IDENT, GT_SYM, IDENT, END_OF_INPUT}},
		{"a<b", "a<b", []int{IDENT, LT, IDENT, END_OF_INPUT}},
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
			}
		})
	}
}
