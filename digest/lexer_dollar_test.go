package digest

import "testing"

// Phase 12: MY_LEX_IDENT_OR_DOLLAR_QUOTED_TEXT tests
// MySQL Reference: sql/sql_lex.cc:2138-2172
//
// Dollar-quoted strings are used in MySQL 8.0 for stored routine bodies.
// Format: $$text$$ or $tag$text$tag$
// - $$ starts an anonymous dollar-quoted string
// - $tag$ starts a tagged dollar-quoted string (tag must match at end)
// - $ followed by identifier chars (not $$) is just an identifier

func TestLexer_Dollar_SimpleString(t *testing.T) {
	// $$text$$ -> DOLLAR_QUOTED_STRING_SYM
	tests := []struct {
		name     string
		input    string
		wantType int
		wantText string
	}{
		{
			name:     "simple_dollar_string",
			input:    "$$hello$$",
			wantType: DOLLAR_QUOTED_STRING_SYM,
			wantText: "$$hello$$",
		},
		{
			name:     "empty_dollar_string",
			input:    "$$$$",
			wantType: DOLLAR_QUOTED_STRING_SYM,
			wantText: "$$$$",
		},
		{
			name:     "dollar_string_with_spaces",
			input:    "$$hello world$$",
			wantType: DOLLAR_QUOTED_STRING_SYM,
			wantText: "$$hello world$$",
		},
		{
			name:     "dollar_string_with_newline",
			input:    "$$line1\nline2$$",
			wantType: DOLLAR_QUOTED_STRING_SYM,
			wantText: "$$line1\nline2$$",
		},
		{
			name:     "dollar_string_with_single_dollar",
			input:    "$$has $ inside$$",
			wantType: DOLLAR_QUOTED_STRING_SYM,
			wantText: "$$has $ inside$$",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.Lex()
			if tok.Type != tt.wantType {
				t.Errorf("got token type %d (%s), want %d (%s)",
					tok.Type, TokenString(tok.Type), tt.wantType, TokenString(tt.wantType))
			}
			gotText := l.MustTokenText(tok)
			if gotText != tt.wantText {
				t.Errorf("got text %q, want %q", gotText, tt.wantText)
			}
		})
	}
}

func TestLexer_Dollar_TaggedString(t *testing.T) {
	// $tag$text$tag$ -> DOLLAR_QUOTED_STRING_SYM
	tests := []struct {
		name     string
		input    string
		wantType int
		wantText string
	}{
		{
			name:     "simple_tag",
			input:    "$x$hello$x$",
			wantType: DOLLAR_QUOTED_STRING_SYM,
			wantText: "$x$hello$x$",
		},
		{
			name:     "longer_tag",
			input:    "$body$CREATE PROCEDURE...$body$",
			wantType: DOLLAR_QUOTED_STRING_SYM,
			wantText: "$body$CREATE PROCEDURE...$body$",
		},
		{
			name:     "tag_with_nested_dollars",
			input:    "$sql$SELECT $$value$$ FROM t$sql$",
			wantType: DOLLAR_QUOTED_STRING_SYM,
			wantText: "$sql$SELECT $$value$$ FROM t$sql$",
		},
		{
			name:     "empty_tagged",
			input:    "$tag$$tag$",
			wantType: DOLLAR_QUOTED_STRING_SYM,
			wantText: "$tag$$tag$",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.Lex()
			if tok.Type != tt.wantType {
				t.Errorf("got token type %d (%s), want %d (%s)",
					tok.Type, TokenString(tok.Type), tt.wantType, TokenString(tt.wantType))
			}
			gotText := l.MustTokenText(tok)
			if gotText != tt.wantText {
				t.Errorf("got text %q, want %q", gotText, tt.wantText)
			}
		})
	}
}

func TestLexer_Dollar_Identifier(t *testing.T) {
	// $ident or ident$ident -> IDENT (not a dollar-quoted string)
	tests := []struct {
		name     string
		input    string
		wantType int
		wantText string
	}{
		{
			name:     "dollar_alone",
			input:    "$",
			wantType: IDENT,
			wantText: "$",
		},
		{
			name:     "dollar_ident",
			input:    "$var",
			wantType: IDENT,
			wantText: "$var",
		},
		{
			name:     "dollar_in_middle",
			input:    "var$name",
			wantType: IDENT,
			wantText: "var$name",
		},
		{
			name:     "dollar_followed_by_non_ident_non_dollar",
			input:    "$ ",
			wantType: IDENT,
			wantText: "$",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.Lex()
			if tok.Type != tt.wantType {
				t.Errorf("got token type %d (%s), want %d (%s)",
					tok.Type, TokenString(tok.Type), tt.wantType, TokenString(tt.wantType))
			}
			gotText := l.MustTokenText(tok)
			if gotText != tt.wantText {
				t.Errorf("got text %q, want %q", gotText, tt.wantText)
			}
		})
	}
}

func TestLexer_Dollar_Unclosed(t *testing.T) {
	// Unclosed dollar-quoted strings -> ABORT_SYM
	tests := []struct {
		name     string
		input    string
		wantType int
	}{
		{
			name:     "unclosed_double_dollar",
			input:    "$$hello",
			wantType: ABORT_SYM,
		},
		{
			name:     "unclosed_with_single_dollar",
			input:    "$$hello$",
			wantType: ABORT_SYM,
		},
		{
			name:     "unclosed_tagged",
			input:    "$tag$hello",
			wantType: ABORT_SYM,
		},
		{
			name:     "unclosed_tagged_wrong_end",
			input:    "$tag$hello$other$",
			wantType: ABORT_SYM,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tok := l.Lex()
			if tok.Type != tt.wantType {
				t.Errorf("got token type %d (%s), want %d (%s)",
					tok.Type, TokenString(tok.Type), tt.wantType, TokenString(tt.wantType))
			}
		})
	}
}

func TestLexer_Dollar_InContext(t *testing.T) {
	// Dollar-quoted strings in SQL context
	tests := []struct {
		name   string
		input  string
		tokens []int
	}{
		{
			name:   "create_procedure",
			input:  "CREATE PROCEDURE p() $$SELECT 1$$",
			tokens: []int{CREATE, PROCEDURE_SYM, IDENT, '(', ')', DOLLAR_QUOTED_STRING_SYM},
		},
		{
			name:   "dollar_string_followed_by_semicolon",
			input:  "$$body$$;",
			tokens: []int{DOLLAR_QUOTED_STRING_SYM, ';'},
		},
		{
			name:   "multiple_dollar_strings",
			input:  "$$first$$ $$second$$",
			tokens: []int{DOLLAR_QUOTED_STRING_SYM, DOLLAR_QUOTED_STRING_SYM},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			for i, want := range tt.tokens {
				tok := l.Lex()
				if tok.Type != want {
					t.Errorf("token %d: got %d (%s), want %d (%s)",
						i, tok.Type, TokenString(tok.Type), want, TokenString(want))
				}
			}
			// Verify EOF
			tok := l.Lex()
			if tok.Type != END_OF_INPUT {
				t.Errorf("expected END_OF_INPUT, got %d (%s)", tok.Type, TokenString(tok.Type))
			}
		})
	}
}
