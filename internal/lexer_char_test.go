package internal

import (
	"testing"
)

func TestLexer_CHAR_SingleCharTokens(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		char  byte
	}{
		{"plus", "+", '+'},
		{"comma", ",", ','},
		{"open paren", "(", '('},
		{"close paren", ")", ')'},
		{"open bracket", "[", '['},
		{"close bracket", "]", ']'},
		{"open brace", "{", '{'},
		{"close brace", "}", '}'},
		{"tilde", "~", '~'},
		{"percent", "%", '%'},
		{"caret", "^", '^'},
		{"question mark", "?", '?'},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(tc.input)
			tok := l.Lex()

			if tok.Type != int(tc.char) {
				t.Errorf("input %q: expected type %d (%c), got %d",
					tc.input, int(tc.char), tc.char, tok.Type)
			}
			text := l.MustTokenText(tok)
			if text != string(tc.char) {
				t.Errorf("input %q: expected text %q, got %q",
					tc.input, string(tc.char), text)
			}
		})
	}
}

func TestLexer_CHAR_MinusAlone(t *testing.T) {
	l := NewLexer("-")
	tok := l.Lex()

	if tok.Type != int('-') {
		t.Errorf("expected '-' (%d), got %d", int('-'), tok.Type)
	}
}

func TestLexer_CHAR_MinusNotComment(t *testing.T) {
	// "--x" should produce "-", "-", then identifier "x"
	l := NewLexer("--x")

	// First token: '-'
	tok1 := l.Lex()
	if tok1.Type != int('-') {
		t.Errorf("first token: expected '-' (%d), got %d", int('-'), tok1.Type)
	}

	// Second token: '-'
	tok2 := l.Lex()
	if tok2.Type != int('-') {
		t.Errorf("second token: expected '-' (%d), got %d", int('-'), tok2.Type)
	}
}

func TestLexer_CHAR_DoubleDashComment(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"space after", "-- comment\nSELECT"},
		{"tab after", "--\tcomment\nSELECT"},
		{"newline immediate", "--\nSELECT"},
		{"ctrl char after", "--\x01comment\nSELECT"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLexer(tc.input)
			tok := l.Lex()

			// The comment should be skipped, and we should get the next token
			// For now, since we haven't implemented IDENT yet, we just check
			// it's not returning '-'
			if tok.Type == int('-') {
				t.Errorf("input %q: expected comment to be skipped, got '-'", tc.input)
			}
		})
	}
}

func TestLexer_CHAR_JSONSeparator(t *testing.T) {
	l := NewLexer("->")
	tok := l.Lex()

	if tok.Type != JSON_SEPARATOR_SYM {
		t.Errorf("expected JSON_SEPARATOR_SYM (%d), got %d", JSON_SEPARATOR_SYM, tok.Type)
	}
	text := l.MustTokenText(tok)
	if text != "->" {
		t.Errorf("expected text %q, got %q", "->", text)
	}
}

func TestLexer_CHAR_JSONUnquotedSeparator(t *testing.T) {
	l := NewLexer("->>")
	tok := l.Lex()

	if tok.Type != JSON_UNQUOTED_SEPARATOR_SYM {
		t.Errorf("expected JSON_UNQUOTED_SEPARATOR_SYM (%d), got %d", JSON_UNQUOTED_SEPARATOR_SYM, tok.Type)
	}
	text := l.MustTokenText(tok)
	if text != "->>" {
		t.Errorf("expected text %q, got %q", "->>", text)
	}
}

func TestLexer_CHAR_JSONArrowFollowedByOther(t *testing.T) {
	l := NewLexer("->x")

	// First token: ->
	tok1 := l.Lex()
	if tok1.Type != JSON_SEPARATOR_SYM {
		t.Errorf("first token: expected JSON_SEPARATOR_SYM (%d), got %d", JSON_SEPARATOR_SYM, tok1.Type)
	}

	// Second token should start 'x' (will be IDENT in later phases)
	tok2 := l.Lex()
	text := l.MustTokenText(tok2)
	if text == "" {
		t.Errorf("second token: expected identifier 'x', got empty")
	}
}

func TestLexer_CHAR_CloseParenNoSignedNumbers(t *testing.T) {
	l := NewLexer(")-1")

	// First token: ')'
	tok1 := l.Lex()
	if tok1.Type != int(')') {
		t.Errorf("first token: expected ')' (%d), got %d", int(')'), tok1.Type)
	}

	// Second token: '-' (not a negative number, just minus)
	tok2 := l.Lex()
	if tok2.Type != int('-') {
		t.Errorf("second token: expected '-' (%d), got %d", int('-'), tok2.Type)
	}
}

func TestLexer_CHAR_OtherCharsAllowSignedNumbers(t *testing.T) {
	// For now, we just verify the tokens are parsed correctly
	// The "allow signed numbers" behavior will be tested more fully when number parsing is implemented
	l := NewLexer("+-1")

	tok1 := l.Lex()
	if tok1.Type != int('+') {
		t.Errorf("first token: expected '+' (%d), got %d", int('+'), tok1.Type)
	}

	tok2 := l.Lex()
	if tok2.Type != int('-') {
		t.Errorf("second token: expected '-' (%d), got %d", int('-'), tok2.Type)
	}
}

func TestLexer_CHAR_ParamMarkerInPrepareMode(t *testing.T) {
	l := NewLexer("?")
	l.SetPrepareMode(true)
	tok := l.Lex()

	if tok.Type != PARAM_MARKER {
		t.Errorf("expected PARAM_MARKER (%d), got %d", PARAM_MARKER, tok.Type)
	}
}

func TestLexer_CHAR_ParamMarkerWithWhitespace(t *testing.T) {
	l := NewLexer("? ")
	l.SetPrepareMode(true)
	tok := l.Lex()

	if tok.Type != PARAM_MARKER {
		t.Errorf("expected PARAM_MARKER (%d), got %d", PARAM_MARKER, tok.Type)
	}
}

func TestLexer_CHAR_ParamMarkerBeforeIdent(t *testing.T) {
	l := NewLexer("?x")
	l.SetPrepareMode(true)
	tok := l.Lex()

	// Should return '?' as character, not PARAM_MARKER
	if tok.Type != int('?') {
		t.Errorf("expected '?' (%d) before ident, got %d", int('?'), tok.Type)
	}
}

func TestLexer_CHAR_ParamMarkerNotInPrepareMode(t *testing.T) {
	l := NewLexer("?")
	// Default mode, prepare mode is off
	tok := l.Lex()

	// Should return '?' as character, not PARAM_MARKER
	if tok.Type != int('?') {
		t.Errorf("expected '?' (%d) when not in prepare mode, got %d", int('?'), tok.Type)
	}
}

func TestLexer_CHAR_MultipleOperators(t *testing.T) {
	l := NewLexer("(+,-)")

	expected := []byte{'(', '+', ',', '-', ')'}
	for i, exp := range expected {
		tok := l.Lex()
		if tok.Type != int(exp) {
			t.Errorf("token %d: expected %q (%d), got %d", i, exp, int(exp), tok.Type)
		}
	}

	// Should end with END_OF_INPUT
	tok := l.Lex()
	if tok.Type != END_OF_INPUT {
		t.Errorf("expected END_OF_INPUT, got %d", tok.Type)
	}
}
