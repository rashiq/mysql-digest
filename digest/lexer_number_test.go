package digest

import (
	"testing"
)

// ============================================================================
// MY_LEX_NUMBER_IDENT tests - Numbers that start with a digit
// ============================================================================

// TestLexer_NUMBER_HexLiteral tests 0x hex numbers
func TestLexer_NUMBER_HexLiteral(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"simple hex", "0x1AF"},
		{"lowercase", "0xabc"},
		{"uppercase", "0xABC"},
		{"mixed case", "0xAbCdEf"},
		{"single digit", "0x1"},
		{"zero", "0x0"},
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

// TestLexer_NUMBER_HexInvalid tests invalid hex that becomes identifier
func TestLexer_NUMBER_HexInvalid(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"0x alone", "0x"},
		{"0x followed by non-hex", "0xGG"},
		{"0x followed by ident", "0xident"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(tc.input)
			tok := l.Lex()

			// Should become an identifier, not HEX_NUM
			if tok.Type == HEX_NUM {
				t.Errorf("input %q: should not be HEX_NUM", tc.input)
			}
		})
	}
}

// TestLexer_NUMBER_BinLiteral tests 0b binary numbers
func TestLexer_NUMBER_BinLiteral(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"simple binary", "0b1010"},
		{"all zeros", "0b0000"},
		{"all ones", "0b1111"},
		{"single bit", "0b1"},
		{"zero", "0b0"},
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

// TestLexer_NUMBER_BinInvalid tests invalid binary that becomes identifier
func TestLexer_NUMBER_BinInvalid(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"0b alone", "0b"},
		{"0b followed by non-binary", "0b222"},
		{"0b followed by ident", "0bident"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(tc.input)
			tok := l.Lex()

			// Should become an identifier, not BIN_NUM
			if tok.Type == BIN_NUM {
				t.Errorf("input %q: should not be BIN_NUM", tc.input)
			}
		})
	}
}

// TestLexer_NUMBER_Integer tests simple integer numbers
func TestLexer_NUMBER_Integer(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		wantType int
	}{
		{"single digit", "5", NUM},
		{"multiple digits", "123", NUM},
		{"zero", "0", NUM},
		{"large number", "123456789", NUM},
		// LONG_NUM is for numbers that don't fit in regular int
		{"very large", "12345678901234567890", LONG_NUM},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(tc.input)
			tok := l.Lex()

			// Accept NUM, LONG_NUM, or ULONGLONG_NUM
			validTypes := []int{NUM, LONG_NUM, ULONGLONG_NUM}
			valid := false
			for _, vt := range validTypes {
				if tok.Type == vt {
					valid = true
					break
				}
			}
			if !valid {
				t.Errorf("input %q: expected integer type, got %d", tc.input, tok.Type)
			}
		})
	}
}

// TestLexer_NUMBER_IntegerIdent tests numbers followed by letters (becomes ident)
func TestLexer_NUMBER_IntegerIdent(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"digit then letter", "123abc"},
		{"digit then underscore", "123_abc"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(tc.input)
			tok := l.Lex()

			// Should become IDENT, not a number
			if tok.Type != IDENT && tok.Type != IDENT_QUOTED {
				t.Errorf("input %q: expected IDENT, got %d", tc.input, tok.Type)
			}
			text := l.MustTokenText(tok)
			if text != tc.input {
				t.Errorf("input %q: expected text %q, got %q", tc.input, tc.input, text)
			}
		})
	}
}

// TestLexer_NUMBER_Float tests floating point numbers with exponent
func TestLexer_NUMBER_Float(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"simple exponent", "1e10"},
		{"uppercase E", "1E10"},
		{"with plus", "1e+10"},
		{"with minus", "1e-10"},
		{"uppercase with plus", "1E+10"},
		{"decimal with exponent", "1.5e10"},
		{"multi digit base", "123e4"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(tc.input)
			tok := l.Lex()

			if tok.Type != FLOAT_NUM {
				t.Errorf("input %q: expected FLOAT_NUM (%d), got %d",
					tc.input, FLOAT_NUM, tok.Type)
			}
		})
	}
}

// TestLexer_NUMBER_FloatInvalid tests invalid float notation becomes ident
func TestLexer_NUMBER_FloatInvalid(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"1e alone", "1e"},
		{"1e followed by letter", "1ex"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(tc.input)
			tok := l.Lex()

			// Should become IDENT
			if tok.Type == FLOAT_NUM {
				t.Errorf("input %q: should not be FLOAT_NUM", tc.input)
			}
		})
	}
}

// TestLexer_NUMBER_ExponentWithSignButNoDigits tests the edge case where
// exponent has sign but no digits (e.g., "1e+x"). This was a bug where
// position wasn't properly restored.
func TestLexer_NUMBER_ExponentWithSignButNoDigits(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		wantType  int
		wantText  string
		wantNext  int // Expected next token type
		wantNext2 string
	}{
		{
			name:      "1e+ followed by ident",
			input:     "1e+x",
			wantType:  IDENT,
			wantText:  "1e",
			wantNext:  int('+'),
			wantNext2: "+",
		},
		{
			name:      "1e- followed by ident",
			input:     "1e-x",
			wantType:  IDENT,
			wantText:  "1e",
			wantNext:  int('-'),
			wantNext2: "-",
		},
		{
			name:      "decimal with e+ but no digits",
			input:     "1.5e+x",
			wantType:  DECIMAL_NUM,
			wantText:  "1.5",
			wantNext:  IDENT,
			wantNext2: "e",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(tc.input)
			tok := l.Lex()

			if tok.Type != tc.wantType {
				t.Errorf("first token type: got %d (%s), want %d (%s)",
					tok.Type, TokenString(tok.Type), tc.wantType, TokenString(tc.wantType))
			}

			text := l.MustTokenText(tok)
			if text != tc.wantText {
				t.Errorf("first token text: got %q, want %q", text, tc.wantText)
			}

			// Check next token to verify position was properly restored
			tok2 := l.Lex()
			if tok2.Type != tc.wantNext {
				t.Errorf("second token type: got %d (%s), want %d (%s)",
					tok2.Type, TokenString(tok2.Type), tc.wantNext, TokenString(tc.wantNext))
			}
		})
	}
}

// TestLexer_NUMBER_Decimal tests decimal numbers (with decimal point)
func TestLexer_NUMBER_Decimal(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"simple decimal", "123.456"},
		{"leading zero", "0.123"},
		{"trailing zeros", "1.000"},
		{"single decimal place", "1.5"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(tc.input)
			tok := l.Lex()

			if tok.Type != DECIMAL_NUM {
				t.Errorf("input %q: expected DECIMAL_NUM (%d), got %d",
					tc.input, DECIMAL_NUM, tok.Type)
			}
		})
	}
}

// TestLexer_NUMBER_DecimalWithExponent tests decimal with exponent = float
func TestLexer_NUMBER_DecimalWithExponent(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"decimal with e", "123.456e7"},
		{"decimal with E+", "1.5E+10"},
		{"decimal with e-", "0.5e-3"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(tc.input)
			tok := l.Lex()

			if tok.Type != FLOAT_NUM {
				t.Errorf("input %q: expected FLOAT_NUM (%d), got %d",
					tc.input, FLOAT_NUM, tok.Type)
			}
		})
	}
}

// TestLexer_NUMBER_LeadingDot tests numbers starting with decimal point
func TestLexer_NUMBER_LeadingDot(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		wantType int
	}{
		{"dot then digits", ".456", DECIMAL_NUM},
		{"dot then digits then e", ".5e10", FLOAT_NUM},
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

// TestLexer_NUMBER_DotNotNumber tests dot followed by non-digit
func TestLexer_NUMBER_DotNotNumber(t *testing.T) {
	l := NewLexer(".abc")

	// First token should be '.'
	tok1 := l.Lex()
	if tok1.Type != int('.') {
		t.Errorf("expected '.' (%d), got %d", int('.'), tok1.Type)
	}

	// Second should be IDENT
	tok2 := l.Lex()
	if tok2.Type != IDENT {
		t.Errorf("expected IDENT, got %d", tok2.Type)
	}
}

// TestLexer_NUMBER_Sequence tests number in context
func TestLexer_NUMBER_Sequence(t *testing.T) {
	l := NewLexer("SELECT 123, 45.67")

	// SELECT
	tok := l.Lex()
	if tok.Type != SELECT_SYM {
		t.Errorf("expected SELECT_SYM, got %d", tok.Type)
	}

	// 123 (integer)
	tok = l.Lex()
	if tok.Type != NUM && tok.Type != LONG_NUM {
		t.Errorf("expected NUM, got %d", tok.Type)
	}

	// comma
	tok = l.Lex()
	if tok.Type != int(',') {
		t.Errorf("expected ',', got %d", tok.Type)
	}

	// 45.67 (decimal)
	tok = l.Lex()
	if tok.Type != DECIMAL_NUM {
		t.Errorf("expected DECIMAL_NUM, got %d", tok.Type)
	}
}
