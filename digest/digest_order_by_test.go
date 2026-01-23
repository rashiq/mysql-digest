package digest

import (
	"testing"
)

func TestOrderByNumeric(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		expected string
	}{
		{
			name:     "Order by numeric column",
			sql:      "SELECT * FROM t ORDER BY 1",
			expected: "SELECT * FROM `t` ORDER BY 1",
		},
		{
			name:     "Group by numeric column",
			sql:      "SELECT * FROM t GROUP BY 1",
			expected: "SELECT * FROM `t` GROUP BY 1",
		},
		{
			name:     "Order by multiple numeric columns",
			sql:      "SELECT * FROM t ORDER BY 1, 2",
			expected: "SELECT * FROM `t` ORDER BY 1, 2",
		},
		{
			name:     "Order by mixed columns",
			sql:      "SELECT * FROM t ORDER BY a, 2, b",
			expected: "SELECT * FROM `t` ORDER BY `a`, 2, `b`",
		},
		{
			name:     "Order by numeric and limit",
			sql:      "SELECT * FROM t ORDER BY 1 LIMIT 10",
			expected: "SELECT * FROM `t` ORDER BY 1 LIMIT ?",
		},
		{
			name:     "Order by expression (first part looks like col)",
			sql:      "SELECT * FROM t ORDER BY 1 + 1",
			expected: "SELECT * FROM `t` ORDER BY 1 + ?",
		},
		{
			name:     "Order by expression (start with non-col)",
			sql:      "SELECT * FROM t ORDER BY a + 1",
			expected: "SELECT * FROM `t` ORDER BY `a` + ?",
		},
		{
			name:     "Partition by key",
			sql:      "SELECT * FROM t PARTITION BY KEY(id) PARTITIONS 4",
			expected: "SELECT * FROM `t` PARTITION BY KEY (`id`) PARTITIONS ?",
		},
		{
			// Note: PARTITION BY usually takes expressions or columns, not positional numbers like ORDER BY.
			// But our logic is broad: "in ORDER/GROUP/PARTITION clause".
			// MySQL syntax: PARTITION BY RANGE (store_id) ...
			// PARTITION BY HASH(id) partitions 4.
			// The `4` in `PARTITIONS 4` is NOT a column position.
			// My logic: `PARTITIONS` is not `BY`, `,`, or `(`.
			// So `4` will be `?`. Correct.
			name:     "Partition by partitions count",
			sql:      "CREATE TABLE t (id int) PARTITION BY HASH(id) PARTITIONS 4",
			expected: "CREATE TABLE `t` (`id` INTEGER) PARTITION BY HASH (`id`) PARTITIONS ?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Compute(tt.sql)
			if d.Text != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, d.Text)
			}
		})
	}
}
