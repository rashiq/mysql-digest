package digest

import (
	"testing"
)

func TestLexer_START_EmptyInput(t *testing.T) {
	l := NewLexer("")
	tok := l.Lex()

	if tok.Type != END_OF_INPUT {
		t.Errorf("expected END_OF_INPUT (%d), got %d", END_OF_INPUT, tok.Type)
	}
	if tok.Start != 0 || tok.End != 0 {
		t.Errorf("expected Start=0, End=0, got Start=%d, End=%d", tok.Start, tok.End)
	}
}

func TestLexer_START_WhitespaceOnly(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"single space", " "},
		{"multiple spaces", "    "},
		{"tab", "\t"},
		{"newline", "\n"},
		{"carriage return", "\r"},
		{"mixed whitespace", "  \t\n\r  "},
		{"vertical tab", "\v"},
		{"form feed", "\f"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(tc.input)
			tok := l.Lex()

			if tok.Type != END_OF_INPUT {
				t.Errorf("input %q: expected END_OF_INPUT (%d), got %d",
					tc.input, END_OF_INPUT, tok.Type)
			}
		})
	}
}

func TestLexer_START_LeadingWhitespaceStripped(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expectedStart int
	}{
		{"space before char", " +", 1},
		{"tab before char", "\t+", 1},
		{"newline before char", "\n+", 1},
		{"multiple spaces before char", "   +", 3},
		{"mixed whitespace before char", " \t\n+", 3},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(tc.input)
			tok := l.Lex()

			if tok.Start != tc.expectedStart {
				t.Errorf("input %q: expected Start=%d, got Start=%d",
					tc.input, tc.expectedStart, tok.Start)
			}
			// The token should be '+' (ASCII 43)
			if tok.Type != int('+') {
				t.Errorf("input %q: expected type '+' (%d), got %d",
					tc.input, int('+'), tok.Type)
			}
		})
	}
}

func TestLexer_START_StateDispatch(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expectedState LexState
	}{
		// Letters -> MY_LEX_IDENT (or variants)
		{"lowercase letter", "a", MY_LEX_IDENT},
		{"uppercase letter", "A", MY_LEX_IDENT},
		{"underscore", "_", MY_LEX_IDENT},

		// Special identifier starters
		{"x for hex", "x", MY_LEX_IDENT_OR_HEX},
		{"X for hex", "X", MY_LEX_IDENT_OR_HEX},
		{"b for bin", "b", MY_LEX_IDENT_OR_BIN},
		{"B for bin", "B", MY_LEX_IDENT_OR_BIN},
		{"n for nchar", "n", MY_LEX_IDENT_OR_NCHAR},
		{"N for nchar", "N", MY_LEX_IDENT_OR_NCHAR},
		{"dollar", "$", MY_LEX_IDENT_OR_DOLLAR_QUOTED_TEXT},

		// Digits -> MY_LEX_NUMBER_IDENT
		{"digit 0", "0", MY_LEX_NUMBER_IDENT},
		{"digit 5", "5", MY_LEX_NUMBER_IDENT},
		{"digit 9", "9", MY_LEX_NUMBER_IDENT},

		// String delimiters
		{"single quote", "'", MY_LEX_STRING},
		{"double quote", "\"", MY_LEX_STRING_OR_DELIMITER},
		{"backtick", "`", MY_LEX_USER_VARIABLE_DELIMITER},

		// Operators
		{"greater than", ">", MY_LEX_CMP_OP},
		{"equals", "=", MY_LEX_CMP_OP},
		{"exclamation", "!", MY_LEX_CMP_OP},
		{"less than", "<", MY_LEX_LONG_CMP_OP},
		{"ampersand", "&", MY_LEX_BOOL},
		{"pipe", "|", MY_LEX_BOOL},

		// Comments
		{"hash", "#", MY_LEX_COMMENT},
		{"slash", "/", MY_LEX_LONG_COMMENT},
		{"asterisk", "*", MY_LEX_END_LONG_COMMENT},

		// Special
		{"semicolon", ";", MY_LEX_SEMICOLON},
		{"colon", ":", MY_LEX_SET_VAR},
		{"at sign", "@", MY_LEX_USER_END},
		{"dot", ".", MY_LEX_REAL_OR_POINT},

		// Regular characters -> MY_LEX_CHAR
		{"plus", "+", MY_LEX_CHAR},
		{"minus", "-", MY_LEX_CHAR},
		{"comma", ",", MY_LEX_CHAR},
		{"open paren", "(", MY_LEX_CHAR},
		{"close paren", ")", MY_LEX_CHAR},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualState := getStateMap(tc.input[0])
			if actualState != tc.expectedState {
				t.Errorf("getStateMap(%q): expected %d, got %d",
					tc.input, tc.expectedState, actualState)
			}
		})
	}
}

func TestLexer_START_MultipleTokens(t *testing.T) {
	l := NewLexer("+ - *")

	// First token: '+'
	tok1 := l.Lex()
	if tok1.Type != int('+') {
		t.Errorf("first token: expected '+' (%d), got %d", int('+'), tok1.Type)
	}

	// Second token: '-'
	tok2 := l.Lex()
	if tok2.Type != int('-') {
		t.Errorf("second token: expected '-' (%d), got %d", int('-'), tok2.Type)
	}

	// Third token: '*'
	// Note: '*' maps to MY_LEX_END_LONG_COMMENT state in the state map,
	// but as a standalone token it should still return the character value
	tok3 := l.Lex()
	text := l.MustTokenText(tok3)
	if text != "*" {
		t.Errorf("third token: expected '*', got %q", text)
	}

	// Fourth token: END_OF_INPUT
	tok4 := l.Lex()
	if tok4.Type != END_OF_INPUT {
		t.Errorf("fourth token: expected END_OF_INPUT (%d), got %d", END_OF_INPUT, tok4.Type)
	}
}

func TestLexer_TokenText(t *testing.T) {
	l := NewLexer("  abc")
	tok := l.Lex()

	text := l.MustTokenText(tok)
	// Note: at this stage, we only have basic CHAR handling,
	// so 'a' will dispatch to IDENT state which falls through to default
	// Just verify TokenText works with positions
	if tok.Start < 0 || tok.End > len("  abc") {
		t.Errorf("invalid token positions: Start=%d, End=%d", tok.Start, tok.End)
	}
	if len(text) != tok.End-tok.Start {
		t.Errorf("TokenText length mismatch: got %d, expected %d", len(text), tok.End-tok.Start)
	}
}

func TestLexer_TokenText_Error(t *testing.T) {
	l := NewLexer("abc")

	tests := []struct {
		name  string
		token Token
	}{
		{
			name:  "negative start",
			token: Token{Type: IDENT, Start: -1, End: 2},
		},
		{
			name:  "end beyond input",
			token: Token{Type: IDENT, Start: 0, End: 100},
		},
		{
			name:  "start after end",
			token: Token{Type: IDENT, Start: 2, End: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text, err := l.TokenText(tt.token)
			if err == nil {
				t.Errorf("expected error for invalid token bounds, got text %q", text)
			}
			if text != "" {
				t.Errorf("expected empty text on error, got %q", text)
			}
		})
	}
}

func TestLexer_SKIP_WhitespaceCharacters(t *testing.T) {
	whitespaceChars := []byte{' ', '\t', '\n', '\r', '\v', '\f'}

	for _, c := range whitespaceChars {
		state := getStateMap(c)
		if state != MY_LEX_SKIP {
			t.Errorf("character %q (0x%02x) should map to MY_LEX_SKIP, got %d",
				c, c, state)
		}
	}
}

func TestLexer_EOL_NullByte(t *testing.T) {
	state := getStateMap(0)
	if state != MY_LEX_EOL {
		t.Errorf("null byte should map to MY_LEX_EOL, got %d", state)
	}
}
