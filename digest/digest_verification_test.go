package digest

import (
	"testing"
)

func TestVerificationSteps(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		expected string
	}{
		{
			name:     "Order by hex literal",
			sql:      "SELECT * FROM t ORDER BY 0x1",
			expected: "SELECT * FROM `t` ORDER BY 0x1",
		},
		{
			name:     "Order by binary literal",
			sql:      "SELECT * FROM t ORDER BY 0b1",
			expected: "SELECT * FROM `t` ORDER BY 0b1",
		},
		{
			name:     "Group by numeric with rollup",
			sql:      "SELECT * FROM t GROUP BY 1 WITH ROLLUP",
			expected: "SELECT * FROM `t` GROUP BY 1 WITH ROLLUP",
		},
		{
			name:     "Group by numeric and limit with rollup",
			sql:      "SELECT * FROM t GROUP BY 1 WITH ROLLUP LIMIT 1",
			expected: "SELECT * FROM `t` GROUP BY 1 WITH ROLLUP LIMIT ?",
		},
		{
			name: "Order by non-numeric hex",
			// If it's a hex string not acting as a column reference, does it matter?
			// In MySQL, ORDER BY 0x41 (ASCII 'A') sorts by the string 'A' if treated as string,
			// or the number 65 if treated as number.
			// Usually valid as a constant.
			sql:      "SELECT * FROM t ORDER BY 0x123",
			expected: "SELECT * FROM `t` ORDER BY 0x123",
		},
		{
			name:     "Partition by numeric literal",
			sql:      "CREATE TABLE t (id int) PARTITION BY 1",
			expected: "CREATE TABLE `t` ( `id` INTEGER ) PARTITION BY 1",
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
