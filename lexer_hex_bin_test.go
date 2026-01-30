package digest

import "testing"

func TestLexer_HEX_NUMBER_Valid(t *testing.T) {
	tests := []struct {
		name  string
		input string
		text  string
	}{
		{"simple hex", "x'1A2B'", "x'1A2B'"},
		{"lowercase x", "x'abcd'", "x'abcd'"},
		{"uppercase X", "X'ABCD'", "X'ABCD'"},
		{"mixed case hex", "x'AbCd'", "x'AbCd'"},
		{"single byte", "x'FF'", "x'FF'"},
		{"two bytes", "x'1234'", "x'1234'"},
		{"empty hex", "x''", "x''"},
		{"all zeros", "x'0000'", "x'0000'"},
		{"all f's", "x'FFFF'", "x'FFFF'"},
		{"long hex", "x'0123456789ABCDEF'", "x'0123456789ABCDEF'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.Lex()
			if tok.Type != HEX_NUM {
				t.Errorf("expected HEX_NUM (%d), got %d", HEX_NUM, tok.Type)
			}
			if got := l.MustTokenText(tok); got != tt.text {
				t.Errorf("expected text %q, got %q", tt.text, got)
			}
		})
	}
}

func TestLexer_HEX_NUMBER_OddDigits(t *testing.T) {
	// MySQL requires even number of hex digits (each byte is 2 hex chars)
	tests := []struct {
		name  string
		input string
	}{
		{"single digit", "x'A'"},
		{"three digits", "x'ABC'"},
		{"five digits", "x'12345'"},
		{"seven digits", "x'1234567'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.Lex()
			if tok.Type != ABORT_SYM {
				t.Errorf("expected ABORT_SYM (%d) for odd hex digits, got %d", ABORT_SYM, tok.Type)
			}
		})
	}
}

func TestLexer_HEX_NUMBER_InvalidChars(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"invalid G", "x'GG'"},
		{"invalid Z", "x'1Z'"},
		{"space in hex", "x'12 34'"},
		{"special char", "x'12#34'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.Lex()
			if tok.Type != ABORT_SYM {
				t.Errorf("expected ABORT_SYM (%d) for invalid hex chars, got %d", ABORT_SYM, tok.Type)
			}
		})
	}
}

func TestLexer_HEX_NUMBER_Unclosed(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"no closing quote", "x'1A2B"},
		{"no closing quote empty", "x'"},
		{"EOF after digits", "x'ABCD"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.Lex()
			if tok.Type != ABORT_SYM {
				t.Errorf("expected ABORT_SYM (%d) for unclosed hex, got %d", ABORT_SYM, tok.Type)
			}
		})
	}
}

func TestLexer_HEX_NUMBER_InContext(t *testing.T) {
	// Test hex literals in typical SQL contexts
	l := NewLexer("SELECT x'48454C4C4F'")

	tok := l.Lex()
	if tok.Type != SELECT_SYM {
		t.Errorf("expected SELECT_SYM, got %d", tok.Type)
	}

	tok = l.Lex()
	if tok.Type != HEX_NUM {
		t.Errorf("expected HEX_NUM, got %d", tok.Type)
	}
	if got := l.MustTokenText(tok); got != "x'48454C4C4F'" {
		t.Errorf("expected x'48454C4C4F', got %q", got)
	}
}

func TestLexer_BIN_NUMBER_Valid(t *testing.T) {
	tests := []struct {
		name  string
		input string
		text  string
	}{
		{"simple binary", "b'1010'", "b'1010'"},
		{"lowercase b", "b'1111'", "b'1111'"},
		{"uppercase B", "B'0000'", "B'0000'"},
		{"single bit", "b'1'", "b'1'"},
		{"empty binary", "b''", "b''"},
		{"all zeros", "b'00000000'", "b'00000000'"},
		{"all ones", "b'11111111'", "b'11111111'"},
		{"long binary", "b'10101010101010101010'", "b'10101010101010101010'"},
		{"mixed", "b'10110100'", "b'10110100'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.Lex()
			if tok.Type != BIN_NUM {
				t.Errorf("expected BIN_NUM (%d), got %d", BIN_NUM, tok.Type)
			}
			if got := l.MustTokenText(tok); got != tt.text {
				t.Errorf("expected text %q, got %q", tt.text, got)
			}
		})
	}
}

func TestLexer_BIN_NUMBER_InvalidChars(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"digit 2", "b'102'"},
		{"digit 9", "b'1091'"},
		{"letter a", "b'10a1'"},
		{"space", "b'10 10'"},
		{"hex digit", "b'10F0'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.Lex()
			if tok.Type != ABORT_SYM {
				t.Errorf("expected ABORT_SYM (%d) for invalid binary chars, got %d", ABORT_SYM, tok.Type)
			}
		})
	}
}

func TestLexer_BIN_NUMBER_Unclosed(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"no closing quote", "b'1010"},
		{"no closing quote empty", "b'"},
		{"EOF after digits", "b'11110000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.Lex()
			if tok.Type != ABORT_SYM {
				t.Errorf("expected ABORT_SYM (%d) for unclosed binary, got %d", ABORT_SYM, tok.Type)
			}
		})
	}
}

func TestLexer_BIN_NUMBER_InContext(t *testing.T) {
	// Test binary literals in typical SQL contexts
	l := NewLexer("SELECT b'10101010'")

	tok := l.Lex()
	if tok.Type != SELECT_SYM {
		t.Errorf("expected SELECT_SYM, got %d", tok.Type)
	}

	tok = l.Lex()
	if tok.Type != BIN_NUM {
		t.Errorf("expected BIN_NUM, got %d", tok.Type)
	}
	if got := l.MustTokenText(tok); got != "b'10101010'" {
		t.Errorf("expected b'10101010', got %q", got)
	}
}

func TestLexer_HEX_BIN_NotLiteral(t *testing.T) {
	// When X or B is not followed by quote, it should be an identifier
	tests := []struct {
		name      string
		input     string
		firstType int
		firstText string
	}{
		{"X alone", "X", IDENT, "X"},
		{"x alone", "x", IDENT, "x"},
		{"B alone", "B", IDENT, "B"},
		{"b alone", "b", IDENT, "b"},
		{"X followed by space", "X 123", IDENT, "X"},
		{"B followed by space", "B 123", IDENT, "B"},
		{"Xidentifier", "Xfoo", IDENT, "Xfoo"},
		{"Bidentifier", "Bbar", IDENT, "Bbar"},
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

func TestLexer_HEX_BIN_Sequence(t *testing.T) {
	// Test a sequence with hex and binary literals
	l := NewLexer("x'FF', b'1010', X'00'")

	expected := []struct {
		typ  int
		text string
	}{
		{HEX_NUM, "x'FF'"},
		{int(','), ","},
		{BIN_NUM, "b'1010'"},
		{int(','), ","},
		{HEX_NUM, "X'00'"},
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
