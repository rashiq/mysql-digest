package digest

import (
	"testing"
)

func TestDigest_LiteralReplacement(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		wantText string
	}{
		{
			name:     "integer literal",
			sql:      "SELECT 1",
			wantText: "SELECT ?",
		},
		{
			name:     "integer in WHERE",
			sql:      "SELECT * FROM users WHERE id = 123",
			wantText: "SELECT * FROM `users` WHERE `id` = ?",
		},
		{
			name:     "multiple integers",
			sql:      "SELECT * FROM users WHERE id = 1 AND age = 25",
			wantText: "SELECT * FROM `users` WHERE `id` = ? AND `age` = ?",
		},
		{
			name:     "negative integer",
			sql:      "SELECT * FROM t WHERE x = -5",
			wantText: "SELECT * FROM `t` WHERE `x` = ?",
		},
		{
			name:     "float literal",
			sql:      "SELECT 3.14",
			wantText: "SELECT ?",
		},
		{
			name:     "scientific notation",
			sql:      "SELECT 1.5e10",
			wantText: "SELECT ?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Compute(tt.sql)
			if d.Text != tt.wantText {
				t.Errorf("Compute(%q).Text = %q, want %q", tt.sql, d.Text, tt.wantText)
			}
		})
	}
}

func TestDigest_StringReplacement(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		wantText string
	}{
		{
			name:     "single quoted string",
			sql:      "SELECT 'hello'",
			wantText: "SELECT ?",
		},
		{
			name:     "double quoted string",
			sql:      `SELECT "world"`,
			wantText: "SELECT ?",
		},
		{
			name:     "string in WHERE",
			sql:      "SELECT * FROM users WHERE name = 'john'",
			wantText: "SELECT * FROM `users` WHERE NAME = ?",
		},
		{
			name:     "escaped quotes",
			sql:      "SELECT 'it''s a test'",
			wantText: "SELECT ?",
		},
		{
			name:     "backslash escape",
			sql:      `SELECT 'hello\nworld'`,
			wantText: "SELECT ?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Compute(tt.sql)
			if d.Text != tt.wantText {
				t.Errorf("Compute(%q).Text = %q, want %q", tt.sql, d.Text, tt.wantText)
			}
		})
	}
}

func TestDigest_INClauseCollapsing(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		wantText string
	}{
		{
			name:     "IN with integers",
			sql:      "SELECT * FROM t WHERE x IN (1, 2, 3)",
			wantText: "SELECT * FROM `t` WHERE `x` IN (...)",
		},
		{
			name:     "IN with strings",
			sql:      "SELECT * FROM t WHERE x IN ('a', 'b', 'c')",
			wantText: "SELECT * FROM `t` WHERE `x` IN (...)",
		},
		{
			name:     "IN with single value",
			sql:      "SELECT * FROM t WHERE x IN (1)",
			wantText: "SELECT * FROM `t` WHERE `x` IN (...)",
		},
		{
			name:     "IN with many values",
			sql:      "SELECT * FROM t WHERE x IN (1, 2, 3, 4, 5, 6, 7, 8, 9, 10)",
			wantText: "SELECT * FROM `t` WHERE `x` IN (...)",
		},
		{
			name:     "NOT IN",
			sql:      "SELECT * FROM t WHERE x NOT IN (1, 2, 3)",
			wantText: "SELECT * FROM `t` WHERE `x` NOT IN (...)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Compute(tt.sql)
			if d.Text != tt.wantText {
				t.Errorf("Compute(%q).Text = %q, want %q", tt.sql, d.Text, tt.wantText)
			}
		})
	}
}

func TestDigest_IdentifierPreservation(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		wantText string
	}{
		{
			name:     "simple identifier",
			sql:      "SELECT foo FROM bar",
			wantText: "SELECT `foo` FROM `bar`",
		},
		{
			name:     "quoted identifier",
			sql:      "SELECT `name` FROM `table`",
			wantText: "SELECT `name` FROM `table`",
		},
		{
			name:     "mixed case identifier",
			sql:      "SELECT MyColumn FROM MyTable",
			wantText: "SELECT `MyColumn` FROM `MyTable`",
		},
		{
			name:     "table.column",
			sql:      "SELECT t.id FROM users t",
			wantText: "SELECT `t`.`id` FROM `users` `t`",
		},
		{
			name:     "database.table.column",
			sql:      "SELECT db.t.id FROM db.users",
			wantText: "SELECT `db`.`t`.`id` FROM `db`.`users`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Compute(tt.sql)
			if d.Text != tt.wantText {
				t.Errorf("Compute(%q).Text = %q, want %q", tt.sql, d.Text, tt.wantText)
			}
		})
	}
}

func TestDigest_HashConsistency(t *testing.T) {
	// Same normalized SQL should produce same hash
	sqls := []string{
		"SELECT * FROM users WHERE id = 1",
		"SELECT * FROM users WHERE id = 2",
		"SELECT * FROM users WHERE id = 999",
	}

	var firstHash string
	for i, sql := range sqls {
		d := Compute(sql)
		if i == 0 {
			firstHash = d.Hash
		} else {
			if d.Hash != firstHash {
				t.Errorf("Hash mismatch: %q produced %q, expected %q", sql, d.Hash, firstHash)
			}
		}
	}

	// Different normalized SQL should produce different hash
	d1 := Compute("SELECT * FROM users WHERE id = 1")
	d2 := Compute("SELECT * FROM orders WHERE id = 1")
	if d1.Hash == d2.Hash {
		t.Errorf("Different queries should have different hashes")
	}
}

func TestDigest_HashFormat(t *testing.T) {
	d := Compute("SELECT 1")
	// SHA-256 hash should be 64 hex characters
	if len(d.Hash) != 64 {
		t.Errorf("Hash length = %d, want 64", len(d.Hash))
	}
	// Should be valid hex
	for _, c := range d.Hash {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			t.Errorf("Hash contains invalid character: %c", c)
		}
	}
}

func TestDigest_Comments(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		wantText string
	}{
		{
			name:     "line comment at end",
			sql:      "SELECT 1 -- comment",
			wantText: "SELECT ?",
		},
		{
			name:     "block comment",
			sql:      "SELECT /* comment */ 1",
			wantText: "SELECT ?",
		},
		{
			name:     "hash comment",
			sql:      "SELECT 1 # comment",
			wantText: "SELECT ?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Compute(tt.sql)
			if d.Text != tt.wantText {
				t.Errorf("Compute(%q).Text = %q, want %q", tt.sql, d.Text, tt.wantText)
			}
		})
	}
}

func TestDigest_HexAndBinary(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		wantText string
	}{
		{
			name:     "hex literal",
			sql:      "SELECT 0xDEADBEEF",
			wantText: "SELECT ?",
		},
		{
			name:     "hex string",
			sql:      "SELECT x'DEADBEEF'",
			wantText: "SELECT ?",
		},
		{
			name:     "binary literal",
			sql:      "SELECT 0b1010",
			wantText: "SELECT ?",
		},
		{
			name:     "binary string",
			sql:      "SELECT b'1010'",
			wantText: "SELECT ?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Compute(tt.sql)
			if d.Text != tt.wantText {
				t.Errorf("Compute(%q).Text = %q, want %q", tt.sql, d.Text, tt.wantText)
			}
		})
	}
}

func TestDigest_VALUESCollapsing(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		wantText string
	}{
		{
			name:     "single row",
			sql:      "INSERT INTO t VALUES (1)",
			wantText: "INSERT INTO `t` VALUES (?)",
		},
		{
			name:     "single row multiple columns",
			sql:      "INSERT INTO t VALUES (1, 'a', 3)",
			wantText: "INSERT INTO `t` VALUES (...)",
		},
		{
			name:     "multiple rows",
			sql:      "INSERT INTO t VALUES (1), (2), (3)",
			wantText: "INSERT INTO `t` VALUES (?) /* , ... */",
		},
		{
			name:     "multiple rows with columns",
			sql:      "INSERT INTO t VALUES (1, 'a'), (2, 'b'), (3, 'c')",
			wantText: "INSERT INTO `t` VALUES (...) /* , ... */",
		},
		{
			name:     "with column list",
			sql:      "INSERT INTO t (col1, col2) VALUES (1, 'a'), (2, 'b')",
			wantText: "INSERT INTO `t` (`col1`, `col2`) VALUES (...) /* , ... */",
		},
		{
			name:     "with ON DUPLICATE KEY",
			sql:      "INSERT INTO t VALUES (1, 'a'), (2, 'b') ON DUPLICATE KEY UPDATE x = 1",
			wantText: "INSERT INTO `t` VALUES (...) /* , ... */ ON DUPLICATE KEY UPDATE `x` = ?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Compute(tt.sql)
			if d.Text != tt.wantText {
				t.Errorf("Compute(%q).Text = %q, want %q", tt.sql, d.Text, tt.wantText)
			}
		})
	}
}

func TestDigest_NullHandling(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		wantText string
	}{
		{
			name:     "NULL in WHERE",
			sql:      "SELECT * FROM t WHERE x = NULL",
			wantText: "SELECT * FROM `t` WHERE `x` = ?",
		},
		{
			name:     "IS NULL",
			sql:      "SELECT * FROM t WHERE x IS NULL",
			wantText: "SELECT * FROM `t` WHERE `x` IS ?",
		},
		{
			name:     "NULL in IN clause",
			sql:      "SELECT * FROM t WHERE x IN (1, NULL, 3)",
			wantText: "SELECT * FROM `t` WHERE `x` IN (...)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Compute(tt.sql)
			if d.Text != tt.wantText {
				t.Errorf("Compute(%q).Text = %q, want %q", tt.sql, d.Text, tt.wantText)
			}
		})
	}
}

func TestDigest_ComplexQueries(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		wantText string
	}{
		{
			name:     "JOIN query",
			sql:      "SELECT u.id, o.total FROM users u JOIN orders o ON u.id = o.user_id WHERE o.total > 100",
			wantText: "SELECT `u`.`id`, `o`.`total` FROM `users` `u` JOIN `orders` `o` ON `u`.`id` = `o`.`user_id` WHERE `o`.`total` > ?",
		},
		{
			name:     "GROUP BY with HAVING",
			sql:      "SELECT status, COUNT(*) FROM orders GROUP BY status HAVING COUNT(*) > 5",
			wantText: "SELECT STATUS, COUNT (*) FROM `orders` GROUP BY STATUS HAVING COUNT (*) > ?",
		},
		{
			name:     "LIMIT and OFFSET",
			sql:      "SELECT * FROM users LIMIT 10 OFFSET 20",
			wantText: "SELECT * FROM `users` LIMIT ? OFFSET ?",
		},
		{
			name:     "INSERT",
			sql:      "INSERT INTO users (name, age) VALUES ('John', 25)",
			wantText: "INSERT INTO `users` (NAME, `age`) VALUES (...)",
		},
		{
			name:     "UPDATE",
			sql:      "UPDATE users SET name = 'Jane', age = 30 WHERE id = 1",
			wantText: "UPDATE `users` SET NAME = ?, `age` = ? WHERE `id` = ?",
		},
		{
			name:     "DELETE",
			sql:      "DELETE FROM users WHERE id = 1",
			wantText: "DELETE FROM `users` WHERE `id` = ?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Compute(tt.sql)
			if d.Text != tt.wantText {
				t.Errorf("Compute(%q).Text = %q, want %q", tt.sql, d.Text, tt.wantText)
			}
		})
	}
}

func TestDigest_OptimizerHints(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		wantText string
	}{
		{
			name:     "MAX_EXECUTION_TIME hint",
			sql:      "SELECT /*+ MAX_EXECUTION_TIME(1000) */ * FROM t",
			wantText: "SELECT /*+ MAX_EXECUTION_TIME (?) */ * FROM `t`",
		},
		{
			name:     "SET_VAR with string value",
			sql:      "SELECT /*+ SET_VAR(sql_mode = 'STRICT') */ 1",
			wantText: "SELECT /*+ SET_VAR (`sql_mode` = ?) */ ?",
		},
		{
			name:     "SET_VAR with numeric value",
			sql:      "SELECT /*+ SET_VAR(sort_buffer_size=1000000) */ * FROM t",
			wantText: "SELECT /*+ SET_VAR (`sort_buffer_size` = ?) */ * FROM `t`",
		},
		{
			name:     "multiple hints",
			sql:      "SELECT /*+ MAX_EXECUTION_TIME(8000) SET_VAR(sql_mode = 'STRICT_ALL_TABLES') */ * FROM t WHERE x IN (1, 2, 6)",
			wantText: "SELECT /*+ MAX_EXECUTION_TIME (?) SET_VAR (`sql_mode` = ?) */ * FROM `t` WHERE `x` IN (...)",
		},
		{
			name:     "NO_INDEX hint",
			sql:      "SELECT /*+ NO_INDEX(t1 idx1) */ * FROM t1",
			wantText: "SELECT /*+ NO_INDEX (`t1` `idx1`) */ * FROM `t1`",
		},
		{
			name:     "empty hint",
			sql:      "SELECT /*+ */ 1",
			wantText: "SELECT /*+ */ ?",
		},
		{
			name:     "hint in UPDATE",
			sql:      "UPDATE /*+ NO_MERGE(t1) */ t1 SET x = 1",
			wantText: "UPDATE /*+ NO_MERGE (`t1`) */ `t1` SET `x` = ?",
		},
		{
			name:     "hint in DELETE",
			sql:      "DELETE /*+ BKA(t1) */ FROM t1 WHERE id = 5",
			wantText: "DELETE /*+ BKA (`t1`) */ FROM `t1` WHERE `id` = ?",
		},
		{
			name:     "hint in INSERT",
			sql:      "INSERT /*+ SET_VAR(foreign_key_checks=0) */ INTO t VALUES (1, 2)",
			wantText: "INSERT /*+ SET_VAR (`foreign_key_checks` = ?) */ INTO `t` VALUES (...)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Compute(tt.sql)
			if d.Text != tt.wantText {
				t.Errorf("Compute(%q).Text = %q, want %q", tt.sql, d.Text, tt.wantText)
			}
		})
	}
}
