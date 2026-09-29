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
