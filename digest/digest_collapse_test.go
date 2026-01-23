package digest

import (
	"testing"
)

func TestDigest_Collapsing(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		expected string
	}{
		{
			name:     "IN with simple values",
			sql:      "SELECT * FROM t WHERE a IN (1, 2, 3)",
			expected: "SELECT * FROM `t` WHERE `a` IN (?)",
		},
		{
			name:     "IN with row constructors",
			sql:      "SELECT * FROM t WHERE (a, b) IN ((1, 2), (3, 4))",
			expected: "SELECT * FROM `t` WHERE (`a`, `b`) IN (?)",
			// MySQL 8.0 digest might output IN (?) or IN ((?)) or something else.
			// Based on TOK_IN_GENERIC_VALUE_EXPRESSION, it likely collapses the whole thing.
		},
		{
			name:     "VALUES with multiple rows",
			sql:      "INSERT INTO t VALUES (1, 2), (3, 4), (5, 6)",
			expected: "INSERT INTO `t` VALUES (?)",
		},
		{
			name:     "VALUES with multiple rows and columns",
			sql:      "INSERT INTO t (a, b) VALUES (1, 2), (3, 4)",
			expected: "INSERT INTO `t` (`a`, `b`) VALUES (?)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Compute(tt.sql)
			if d.Text != tt.expected {
				t.Errorf("\nInput:    %s\nExpected: %s\nGot:      %s", tt.sql, tt.expected, d.Text)
			}
		})
	}
}
