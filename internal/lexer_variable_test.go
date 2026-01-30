package internal

import (
	"testing"
)

func TestLexer_UserVar_At(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantTypes []int
		wantTexts []string
	}{
		{
			"at alone",
			"@",
			[]int{int('@'), END_OF_INPUT},
			[]string{"@", ""},
		},
		{
			// User variable with unquoted identifier returns LEX_HOSTNAME
			// (matches MySQL sql_lexer.cc MY_LEX_HOSTNAME state)
			"at with identifier",
			"@var",
			[]int{int('@'), LEX_HOSTNAME, END_OF_INPUT},
			[]string{"@", "var", ""},
		},
		{
			// Single-quoted user variable returns TEXT_STRING
			"at with quoted string single",
			"@'var'",
			[]int{int('@'), TEXT_STRING, END_OF_INPUT},
			[]string{"@", "'var'", ""},
		},
		{
			// Double-quoted user variable returns TEXT_STRING
			"at with quoted string double",
			"@\"var\"",
			[]int{int('@'), TEXT_STRING, END_OF_INPUT},
			[]string{"@", "\"var\"", ""},
		},
		{
			// Backtick-quoted user variable returns IDENT_QUOTED
			// (MySQL preserves these as identifiers, not normalized)
			"at with backtick",
			"@`var`",
			[]int{int('@'), IDENT_QUOTED, END_OF_INPUT},
			[]string{"@", "`var`", ""},
		},
		{
			// User variable with underscore prefix returns LEX_HOSTNAME
			"at with underscore var",
			"@_myvar",
			[]int{int('@'), LEX_HOSTNAME, END_OF_INPUT},
			[]string{"@", "_myvar", ""},
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

func TestLexer_UserVar_SystemVar(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantTypes []int
		wantTexts []string
	}{
		{
			"double at",
			"@@",
			[]int{int('@'), int('@'), END_OF_INPUT},
			[]string{"@", "@", ""},
		},
		{
			"double at with var",
			"@@var",
			[]int{int('@'), int('@'), IDENT, END_OF_INPUT},
			[]string{"@", "@", "var", ""},
		},
		{
			"double at global",
			"@@global",
			[]int{int('@'), int('@'), GLOBAL_SYM, END_OF_INPUT},
			[]string{"@", "@", "global", ""},
		},
		{
			"double at session",
			"@@session",
			[]int{int('@'), int('@'), SESSION_SYM, END_OF_INPUT},
			[]string{"@", "@", "session", ""},
		},
		{
			"double at global.var",
			"@@global.var",
			[]int{int('@'), int('@'), GLOBAL_SYM, int('.'), IDENT, END_OF_INPUT},
			[]string{"@", "@", "global", ".", "var", ""},
		},
		{
			"double at session.var",
			"@@session.var",
			[]int{int('@'), int('@'), SESSION_SYM, int('.'), IDENT, END_OF_INPUT},
			[]string{"@", "@", "session", ".", "var", ""},
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

func TestLexer_UserVar_InContext(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantTypes []int
	}{
		{
			"SELECT @var",
			"SELECT @var",
			[]int{SELECT_SYM, int('@'), LEX_HOSTNAME, END_OF_INPUT},
		},
		{
			"SET @var = 5",
			"SET @var = 5",
			[]int{SET_SYM, int('@'), LEX_HOSTNAME, EQ, NUM, END_OF_INPUT},
		},
		{
			"SET @var := 5",
			"SET @var := 5",
			[]int{SET_SYM, int('@'), LEX_HOSTNAME, SET_VAR, NUM, END_OF_INPUT},
		},
		{
			"SELECT @@version",
			"SELECT @@version",
			[]int{SELECT_SYM, int('@'), int('@'), IDENT, END_OF_INPUT},
		},
		{
			"SET @@session.sql_mode = ''",
			"SET @@session.sql_mode = ''",
			[]int{SET_SYM, int('@'), int('@'), SESSION_SYM, int('.'), IDENT, EQ, TEXT_STRING, END_OF_INPUT},
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
			}
		})
	}
}

func TestLexer_UserVar_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantTypes []int
	}{
		{
			"triple at",
			"@@@",
			[]int{int('@'), int('@'), int('@'), END_OF_INPUT},
		},
		{
			"at in expression",
			"@a + @b",
			[]int{int('@'), LEX_HOSTNAME, int('+'), int('@'), LEX_HOSTNAME, END_OF_INPUT},
		},
		{
			"at with number",
			"@1var",
			[]int{int('@'), LEX_HOSTNAME, END_OF_INPUT}, // 1var is a hostname starting with digit
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
			}
		})
	}
}
