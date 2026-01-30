package internal

import "testing"

func TestLexer_Hint_Basic(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		tokens []int
	}{
		{
			name:   "select_with_max_execution_time",
			input:  "SELECT /*+ MAX_EXECUTION_TIME(1000) */ 1",
			tokens: []int{SELECT_SYM, TOK_HINT_COMMENT_OPEN, MAX_EXECUTION_TIME_HINT, '(', NUM, ')', TOK_HINT_COMMENT_CLOSE, NUM},
		},
		{
			name:   "select_with_no_index",
			input:  "SELECT /*+ NO_INDEX(t) */ * FROM t",
			tokens: []int{SELECT_SYM, TOK_HINT_COMMENT_OPEN, NO_INDEX_HINT, '(', IDENT, ')', TOK_HINT_COMMENT_CLOSE, '*', FROM, IDENT},
		},
		{
			name:   "update_with_hint",
			input:  "UPDATE /*+ NO_MERGE(t1) */ t1 SET x=1",
			tokens: []int{UPDATE_SYM, TOK_HINT_COMMENT_OPEN, NO_DERIVED_MERGE_HINT, '(', IDENT, ')', TOK_HINT_COMMENT_CLOSE, IDENT, SET_SYM, IDENT, EQ, NUM},
		},
		{
			name:   "delete_with_hint",
			input:  "DELETE /*+ BKA(t1) */ FROM t1",
			tokens: []int{DELETE_SYM, TOK_HINT_COMMENT_OPEN, BKA_HINT, '(', IDENT, ')', TOK_HINT_COMMENT_CLOSE, FROM, IDENT},
		},
		{
			name:   "insert_with_hint",
			input:  "INSERT /*+ SET_VAR(sort_buffer_size=1) */ INTO t VALUES(1)",
			tokens: []int{INSERT_SYM, TOK_HINT_COMMENT_OPEN, SET_VAR_HINT, '(', IDENT, '=', NUM, ')', TOK_HINT_COMMENT_CLOSE, INTO, IDENT, VALUES, '(', NUM, ')'},
		},
		{
			name:   "replace_with_hint",
			input:  "REPLACE /*+ QB_NAME(qb1) */ INTO t VALUES(1)",
			tokens: []int{REPLACE_SYM, TOK_HINT_COMMENT_OPEN, QB_NAME_HINT, '(', IDENT, ')', TOK_HINT_COMMENT_CLOSE, INTO, IDENT, VALUES, '(', NUM, ')'},
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

func TestLexer_Hint_NotAfterHintable(t *testing.T) {
	// /*+ ... */ not after hintable keyword should be skipped as regular comment
	tests := []struct {
		name   string
		input  string
		tokens []int
	}{
		{
			name:   "hint_at_start",
			input:  "/*+ MAX_EXECUTION_TIME(1000) */ SELECT 1",
			tokens: []int{SELECT_SYM, NUM},
		},
		{
			name:   "hint_after_from",
			input:  "SELECT * FROM /*+ NO_INDEX(t) */ t",
			tokens: []int{SELECT_SYM, '*', FROM, IDENT},
		},
		{
			name:   "hint_after_where",
			input:  "SELECT * FROM t WHERE /*+ INDEX(t) */ x=1",
			tokens: []int{SELECT_SYM, '*', FROM, IDENT, WHERE, IDENT, EQ, NUM},
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

func TestLexer_Hint_MultipleHints(t *testing.T) {
	// Multiple hints in a single comment
	input := "SELECT /*+ MAX_EXECUTION_TIME(100) NO_INDEX(t1) */ * FROM t1"
	l := NewLexer(input)

	expected := []int{
		SELECT_SYM,
		TOK_HINT_COMMENT_OPEN,
		MAX_EXECUTION_TIME_HINT,
		'(', NUM, ')',
		NO_INDEX_HINT,
		'(', IDENT, ')',
		TOK_HINT_COMMENT_CLOSE,
		'*', FROM, IDENT,
	}

	for i, want := range expected {
		tok := l.Lex()
		if tok.Type != want {
			t.Errorf("token %d: got %d (%s), want %d (%s)",
				i, tok.Type, TokenString(tok.Type), want, TokenString(want))
		}
	}
}

func TestLexer_Hint_EmptyHint(t *testing.T) {
	// Empty hint comment
	input := "SELECT /*+ */ 1"
	l := NewLexer(input)

	expected := []int{
		SELECT_SYM,
		TOK_HINT_COMMENT_OPEN,
		TOK_HINT_COMMENT_CLOSE,
		NUM,
	}

	for i, want := range expected {
		tok := l.Lex()
		if tok.Type != want {
			t.Errorf("token %d: got %d (%s), want %d (%s)",
				i, tok.Type, TokenString(tok.Type), want, TokenString(want))
		}
	}
}

func TestLexer_Hint_WithStrings(t *testing.T) {
	// Hint with string values (SET_VAR with sql_mode)
	input := "SELECT /*+ SET_VAR(sql_mode='ANSI') */ 1"
	l := NewLexer(input)

	expected := []int{
		SELECT_SYM,
		TOK_HINT_COMMENT_OPEN,
		SET_VAR_HINT,
		'(',
		IDENT,       // sql_mode
		'=',         // =
		TEXT_STRING, // 'ANSI'
		')',
		TOK_HINT_COMMENT_CLOSE,
		NUM, // 1
	}

	for i, want := range expected {
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
}

func TestLexer_Hint_WithBacktickIdent(t *testing.T) {
	// Hint with backtick-quoted identifier
	input := "SELECT /*+ SET_VAR(`sql_mode`='STRICT') */ 1"
	l := NewLexer(input)

	expected := []int{
		SELECT_SYM,
		TOK_HINT_COMMENT_OPEN,
		SET_VAR_HINT,
		'(',
		IDENT, // `sql_mode` - backtick-quoted
		'=',
		TEXT_STRING, // 'STRICT'
		')',
		TOK_HINT_COMMENT_CLOSE,
		NUM, // 1
	}

	for i, want := range expected {
		tok := l.Lex()
		if tok.Type != want {
			t.Errorf("token %d: got %d (%s), want %d (%s)",
				i, tok.Type, TokenString(tok.Type), want, TokenString(want))
		}
	}
}

func TestLexer_Hint_EscapedQuotes(t *testing.T) {
	// Test escaped quotes inside strings
	tests := []struct {
		name   string
		input  string
		tokens []int
	}{
		{
			name:   "escaped_single_quote",
			input:  "SELECT /*+ SET_VAR(x='a''b') */ 1",
			tokens: []int{SELECT_SYM, TOK_HINT_COMMENT_OPEN, SET_VAR_HINT, '(', IDENT, '=', TEXT_STRING, ')', TOK_HINT_COMMENT_CLOSE, NUM},
		},
		{
			name:   "escaped_backtick",
			input:  "SELECT /*+ SET_VAR(`col``name`=1) */ 1",
			tokens: []int{SELECT_SYM, TOK_HINT_COMMENT_OPEN, SET_VAR_HINT, '(', IDENT, '=', NUM, ')', TOK_HINT_COMMENT_CLOSE, NUM},
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
		})
	}
}
