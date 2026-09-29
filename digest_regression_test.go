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
