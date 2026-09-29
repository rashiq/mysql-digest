package digest

import "testing"

func TestVersionCommentBoundary(t *testing.T) {
	for _, pair := range [][2]string{
		{"SELECT /*!50000 1 */+2", "SELECT 1+2"},
		{"SELECT /*! 1 */AS x", "SELECT 1 AS x"},
		{"SELECT /*!50000 1 */", "SELECT 1"},
	} {
		got, err := Compute(pair[0])
		if err != nil {
			t.Fatal(err)
		}
		want, err := Compute(pair[1])
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Errorf("%q: got %+v, want %+v", pair[0], got, want)
		}
	}
}

func TestUnterminatedInput(t *testing.T) {
	for _, sql := range []string{
		"SELECT $$unterminated", "SELECT $tag$unterminated",
		"SELECT /*!50000 1", "SELECT /*! 1", "SELECT /* comment",
		"SELECT /*+ INDEX(t", "SELECT 'text", "SELECT `name",
	} {
		if _, err := Compute(sql); err == nil {
			t.Errorf("%q: expected an error", sql)
		}
	}
}

func TestQuotedHintIdentifiers(t *testing.T) {
	got, err := Compute("SELECT /*+ INDEX(`t` `idx`) */ * FROM t")
	if err != nil {
		t.Fatal(err)
	}
	want, err := Compute("SELECT /*+ INDEX(t idx) */ * FROM t")
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestHintAfterValue(t *testing.T) {
	for _, value := range []string{"1", "'text'", "*", "(1)"} {
		sql := "SELECT " + value
		want, err := Compute(sql)
		if err != nil {
			t.Fatal(err)
		}
		got, err := Compute(sql + " /*+ MAX_EXECUTION_TIME(1000) */")
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Errorf("%q: got %+v, want %+v", sql, got, want)
		}
	}
}

func TestNationalStringEscapes(t *testing.T) {
	for _, mode := range []SQLMode{0, MODE_NO_BACKSLASH_ESCAPES} {
		for _, literal := range []string{`'it\'s'`, `'it''s'`, `'text\'`} {
			want, wantErr := Compute("SELECT "+literal, Options{SQLMode: mode})
			got, err := Compute("SELECT N"+literal, Options{SQLMode: mode})
			if (err == nil) != (wantErr == nil) || (err == nil && got != want) {
				t.Errorf("mode %d, %s: got %+v, %v; want %+v, %v", mode, literal, got, err, want, wantErr)
			}
		}
	}
}
