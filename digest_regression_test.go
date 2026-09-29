package digest

import (
	"fmt"
	"testing"
	"unicode/utf8"
)

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

func TestVersionConfiguration(t *testing.T) {
	for _, tc := range []struct {
		version MySQLVersion
		server  int
		keyword string
	}{
		{MySQL57, 50700, "SQL_CACHE"},
		{MySQL80, 80000, "MASTER_HOST"},
		{MySQL84, 80400, "QUALIFY"},
		{MySQL90, 90000, "VECTOR"},
	} {
		for _, offset := range []int{0, 1} {
			sql := fmt.Sprintf("SELECT /*!%d %s */ 1", tc.server+offset, tc.keyword)
			want := "SELECT ?"
			if offset == 0 {
				want = "SELECT " + tc.keyword + " ?"
			}
			got, err := Compute(sql, Options{Version: tc.version})
			if err != nil || got.Text != want {
				t.Errorf("version %d, %q: got %q, %v; want %q", tc.version, sql, got.Text, err, want)
			}
		}
	}
	for _, opts := range []Options{{Version: -1}, {Version: 4}, {SQLMode: 4}} {
		if _, err := Compute("SELECT 1", opts); err == nil {
			t.Errorf("%+v: expected an error", opts)
		}
	}
}

func TestTextLimit(t *testing.T) {
	const sql = "SELECT `名`"
	want, err := Compute(sql)
	if err != nil {
		t.Fatal(err)
	}
	for limit := -1; limit <= len(want.Text)+1; limit++ {
		got, err := Compute(sql, Options{MaxLength: limit})
		if err != nil {
			t.Fatal(err)
		}
		if got.Hash != want.Hash || !utf8.ValidString(got.Text) || (limit > 0 && len(got.Text) > limit) {
			t.Errorf("limit %d: invalid result %+v", limit, got)
		}
		if (limit <= 0 || limit >= len(want.Text)) && got.Text != want.Text {
			t.Errorf("limit %d: unexpected truncation %q", limit, got.Text)
		}
	}
}

func TestParameterDigest(t *testing.T) {
	for _, version := range []MySQLVersion{MySQL57, MySQL80, MySQL84, MySQL90} {
		d := NewDigester(Options{Version: version})
		want, err := d.Digest("SELECT * FROM t WHERE id = 1")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := d.Digest("SELECT 'unfinished"); err == nil {
			t.Fatal("expected an error")
		}
		for _, sql := range []string{want.Text, "SELECT * FROM t WHERE id = 2"} {
			got, err := d.Digest(sql)
			if err != nil || got != want {
				t.Errorf("version %d, %q: got %+v, %v; want %+v", version, sql, got, err, want)
			}
		}
	}
}
