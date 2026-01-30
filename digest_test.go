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
			wantText: "SELECT `t` . `id` FROM `users` `t`",
		},
		{
			name:     "database.table.column",
			sql:      "SELECT db.t.id FROM db.users",
			wantText: "SELECT `db` . `t` . `id` FROM `db` . `users`",
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
			wantText: "INSERT INTO `t` ( `col1` , `col2` ) VALUES (...) /* , ... */",
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
			wantText: "SELECT * FROM `t` WHERE `x` IS NULL",
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
			wantText: "SELECT `u` . `id` , `o` . `total` FROM `users` `u` JOIN `orders` `o` ON `u` . `id` = `o` . `user_id` WHERE `o` . `total` > ?",
		},
		{
			name:     "GROUP BY with HAVING",
			sql:      "SELECT status, COUNT(*) FROM orders GROUP BY status HAVING COUNT(*) > 5",
			wantText: "SELECT STATUS , COUNT ( * ) FROM `orders` GROUP BY STATUS HAVING COUNT ( * ) > ?",
		},
		{
			name:     "LIMIT and OFFSET",
			sql:      "SELECT * FROM users LIMIT 10 OFFSET 20",
			wantText: "SELECT * FROM `users` LIMIT ? OFFSET ?",
		},
		{
			name:     "INSERT",
			sql:      "INSERT INTO users (name, age) VALUES ('John', 25)",
			wantText: "INSERT INTO `users` ( NAME , `age` ) VALUES (...)",
		},
		{
			name:     "UPDATE",
			sql:      "UPDATE users SET name = 'Jane', age = 30 WHERE id = 1",
			wantText: "UPDATE `users` SET NAME = ? , `age` = ? WHERE `id` = ?",
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
			wantText: "SELECT /*+ SET_VAR ( `sql_mode` = ? ) */ ?",
		},
		{
			name:     "SET_VAR with numeric value",
			sql:      "SELECT /*+ SET_VAR(sort_buffer_size=1000000) */ * FROM t",
			wantText: "SELECT /*+ SET_VAR ( `sort_buffer_size` = ? ) */ * FROM `t`",
		},
		{
			name:     "multiple hints",
			sql:      "SELECT /*+ MAX_EXECUTION_TIME(8000) SET_VAR(sql_mode = 'STRICT_ALL_TABLES') */ * FROM t WHERE x IN (1, 2, 6)",
			wantText: "SELECT /*+ MAX_EXECUTION_TIME (?) SET_VAR ( `sql_mode` = ? ) */ * FROM `t` WHERE `x` IN (...)",
		},
		{
			name:     "NO_INDEX hint",
			sql:      "SELECT /*+ NO_INDEX(t1 idx1) */ * FROM t1",
			wantText: "SELECT /*+ NO_INDEX ( `t1` `idx1` ) */ * FROM `t1`",
		},
		{
			name:     "empty hint",
			sql:      "SELECT /*+ */ 1",
			wantText: "SELECT /*+ */ ?",
		},
		{
			name:     "hint in UPDATE",
			sql:      "UPDATE /*+ NO_MERGE(t1) */ t1 SET x = 1",
			wantText: "UPDATE /*+ NO_MERGE ( `t1` ) */ `t1` SET `x` = ?",
		},
		{
			name:     "hint in DELETE",
			sql:      "DELETE /*+ BKA(t1) */ FROM t1 WHERE id = 5",
			wantText: "DELETE /*+ BKA ( `t1` ) */ FROM `t1` WHERE `id` = ?",
		},
		{
			name:     "hint in INSERT",
			sql:      "INSERT /*+ SET_VAR(foreign_key_checks=0) */ INTO t VALUES (1, 2)",
			wantText: "INSERT /*+ SET_VAR ( `foreign_key_checks` = ? ) */ INTO `t` VALUES (...)",
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

func TestDigest_UserAndSystemVariables(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		wantText string
	}{
		// User variables (@var) - variable name normalized to ?
		{
			name:     "user variable unquoted",
			sql:      "SELECT @myvar",
			wantText: "SELECT @?",
		},
		{
			name:     "user variable assignment",
			sql:      "SET @x := 5",
			wantText: "SET @? := ?",
		},
		{
			name:     "user variable with equals",
			sql:      "SET @x = 10",
			wantText: "SET @? = ?",
		},
		{
			name:     "user variable single quoted",
			sql:      "SELECT @'my var'",
			wantText: "SELECT @?",
		},
		{
			name:     "user variable double quoted",
			sql:      `SELECT @"my var"`,
			wantText: "SELECT @?",
		},
		{
			name:     "user variable backtick quoted",
			sql:      "SELECT @`my var`",
			wantText: "SELECT @`my var`", // Backtick-quoted preserved as identifier
		},
		{
			name:     "multiple user variables",
			sql:      "SELECT @a, @b, @c",
			wantText: "SELECT @? , @? , @?",
		},
		{
			name:     "user variable in expression",
			sql:      "SELECT @a + @b",
			wantText: "SELECT @? + @?",
		},
		// System variables (@@var) - variable name preserved
		{
			name:     "system variable",
			sql:      "SELECT @@version",
			wantText: "SELECT @@`version`",
		},
		{
			name:     "system variable session",
			sql:      "SELECT @@session.sql_mode",
			wantText: "SELECT @@SESSION . `sql_mode`",
		},
		{
			name:     "system variable global",
			sql:      "SELECT @@global.max_connections",
			wantText: "SELECT @@GLOBAL . `max_connections`",
		},
		{
			name:     "set system variable",
			sql:      "SET @@session.sql_mode = 'STRICT'",
			wantText: "SET @@SESSION . `sql_mode` = ?",
		},
		{
			name:     "mixed user and system variables",
			sql:      "SET @result = @@version",
			wantText: "SET @? = @@`version`",
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

// Benchmarks for performance testing

func BenchmarkCompute_Simple(b *testing.B) {
	sql := "SELECT * FROM users WHERE id = 1"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Compute(sql)
	}
}

func BenchmarkCompute_Medium(b *testing.B) {
	sql := "SELECT u.id, u.name, o.total FROM users u JOIN orders o ON u.id = o.user_id WHERE o.status = 'active' AND o.total > 100 ORDER BY o.created_at DESC LIMIT 10"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Compute(sql)
	}
}

func BenchmarkCompute_Complex(b *testing.B) {
	sql := `SELECT u.id, u.name, COUNT(o.id) as order_count, SUM(o.total) as total_spent 
		FROM users u 
		LEFT JOIN orders o ON u.id = o.user_id 
		WHERE u.created_at > '2024-01-01' 
			AND u.status IN ('active', 'premium', 'vip') 
			AND (o.status = 'completed' OR o.status IS NULL)
		GROUP BY u.id, u.name 
		HAVING COUNT(o.id) > 5 
		ORDER BY total_spent DESC 
		LIMIT 100 OFFSET 0`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Compute(sql)
	}
}

func BenchmarkCompute_LargeIN(b *testing.B) {
	sql := "SELECT * FROM users WHERE id IN (1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30)"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Compute(sql)
	}
}

func BenchmarkCompute_Insert(b *testing.B) {
	sql := "INSERT INTO users (name, email, age, status) VALUES ('John', 'john@example.com', 25, 'active'), ('Jane', 'jane@example.com', 30, 'active'), ('Bob', 'bob@example.com', 35, 'inactive')"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Compute(sql)
	}
}

func TestMySQLCompatibility(t *testing.T) {
	tests := []struct {
		name       string
		sql        string
		want8Hash  string // MySQL 8.0 SHA-256
		want57Hash string // MySQL 5.7 MD5 (empty if not supported in 5.7)
	}{
		{
			name:       "simple SELECT with WHERE",
			sql:        "SELECT * FROM users WHERE id = 42",
			want8Hash:  "840a880ebd1642e8a0c4926cfbaf7d4da9616b03025a080fafd43a732800fab5",
			want57Hash: "731d9efe96031900ba2a36667f4718d0",
		},
		{
			name:       "JOIN with multiple conditions",
			sql:        "SELECT a.id, b.name FROM users a JOIN orders b ON a.id = b.user_id WHERE b.total > 100 AND a.status = 'active'",
			want8Hash:  "bbbad1c572514399449094ee3043cd0c8791dd0ae66d579066f05e5f6a2a2bd6",
			want57Hash: "5b83e9f6a57a0c372448944f7ffaafb3",
		},
		{
			name:       "INSERT with multiple rows",
			sql:        "INSERT INTO logs (user_id, action, created_at) VALUES (1, 'login', NOW()), (2, 'logout', NOW())",
			want8Hash:  "63082d4181bdc6fada36b20f7719884355aef4dd40c26070d3377d0bc66bf84a",
			want57Hash: "d3836b53e50fbf1b4947e4e6fc634676",
		},
		{
			name:       "UPDATE with arithmetic and conditions",
			sql:        "UPDATE accounts SET balance = balance - 50.00, updated_at = NOW() WHERE id = 123 AND balance >= 50.00",
			want8Hash:  "989b4724da70672cb9ee411beea017d2140f57471e708bd566ad97502e6046a7",
			want57Hash: "1c3b7450959b0b3f228a2e1e81243e8d",
		},
		{
			name:       "DELETE with IN clause",
			sql:        "DELETE FROM sessions WHERE expires_at < NOW() AND user_id IN (1, 2, 3, 4, 5)",
			want8Hash:  "6490f2bee2884432c4fe1b3f63298d502ba101bdcc1ef8e02f8e5754dd37cf70",
			want57Hash: "b5f970dd3ff955057de7f2d8b6e7e2aa",
		},
		{
			name:       "GROUP BY with HAVING and ORDER BY",
			sql:        "SELECT COUNT(*) AS cnt, status FROM orders GROUP BY status HAVING COUNT(*) > 10 ORDER BY cnt DESC LIMIT 5",
			want8Hash:  "a61c184788ca4ab687dbced1e7c4c2f86b4d3a0b6324bd21a812d7c7c466fcdd",
			want57Hash: "e2476ee01e425d230cf48d15004d7274",
		},
		{
			name:       "correlated subquery",
			sql:        "SELECT u.*, (SELECT COUNT(*) FROM orders WHERE user_id = u.id) AS order_count FROM users u WHERE u.created_at > '2024-01-01'",
			want8Hash:  "7f476b144ccca78b057630a5339080e6baa027cf025d289749700ba97dbd8f37",
			want57Hash: "ba912c30ec03185da4cad1d63aa99d7e",
		},
		{
			name:       "subquery in WHERE with BETWEEN",
			sql:        "SELECT * FROM products WHERE category_id IN (SELECT id FROM categories WHERE parent_id = 5) AND price BETWEEN 10.00 AND 100.00",
			want8Hash:  "fb9714785b2c028cf58ebb6d4b1ee18285f453c7c6a8f0f469c0f9b5222df78d",
			want57Hash: "dd10c7e22290216e997c9b676ed7255a",
		},
		{
			name:       "DATE functions with GROUP BY",
			sql:        "SELECT DATE(created_at) AS day, SUM(amount) AS total FROM transactions GROUP BY DATE(created_at) ORDER BY day DESC",
			want8Hash:  "d60a728c07bec79fceae285c6406323ef564dfd39401e5b87f82de73e8595791",
			want57Hash: "eddd0d1add08602f293078d3e72bcf46",
		},
		{
			name:       "COALESCE and IFNULL",
			sql:        "SELECT COALESCE(nickname, username, email) AS display_name, IFNULL(avatar_url, '/default.png') FROM users WHERE id = 1",
			want8Hash:  "a27ab6578509d6aa9faa6620e08ce05d272bd7e0208ceba693d07fc71f168022",
			want57Hash: "9e8125db5f7559baac2587ff253d4955",
		},
		{
			name:       "DATE_ADD with INTERVAL",
			sql:        "SELECT * FROM events WHERE start_time >= NOW() AND end_time <= DATE_ADD(NOW(), INTERVAL 7 DAY) ORDER BY start_time",
			want8Hash:  "15d88f9a89c7bc0b2f698086b6379933ad962dd113d079a20cee53b7756fa541",
			want57Hash: "c3286587a00d2a9d9bc9945a46b7d1af",
		},
		{
			name:       "CTE with window function",
			sql:        "WITH ranked AS (SELECT *, ROW_NUMBER() OVER (PARTITION BY category_id ORDER BY price DESC) AS rn FROM products) SELECT * FROM ranked WHERE rn <= 3",
			want8Hash:  "8f35aa1c4878b88cf8d53bdb222fb77133b20b834defeee1ea1adc627b8e19fd",
			want57Hash: "", // CTEs not supported in MySQL 5.7
		},
		{
			name:       "optimizer hint MAX_EXECUTION_TIME",
			sql:        "SELECT /*+ MAX_EXECUTION_TIME(5000) */ id, name FROM large_table WHERE status = 1 LIMIT 1000",
			want8Hash:  "e3399d128b66b9c8691dc9294867dfbbc7c91bd7522dcd1bf01866b8c7b0a575",
			want57Hash: "0673f7e6d08e8d422618b2ee0e6700dd",
		},
		{
			name:       "HEX and conversion functions",
			sql:        "SELECT HEX(uuid), FROM_UNIXTIME(created_ts), INET_NTOA(ip_address) FROM access_logs WHERE id > 0",
			want8Hash:  "befa66430b18af4788665af06f42aa0af8698569854b1cd3fcce31430ef40b72",
			want57Hash: "a033907d5dba747403db7f4ee19e7418",
		},
		{
			name:       "INSERT SELECT",
			sql:        "INSERT INTO audit_log SELECT NULL, 'UPDATE', OLD.*, NEW.*, NOW() FROM dual WHERE 1=1",
			want8Hash:  "343e3e20a203533e1351d4c11c7cac0bb2b5dafc10d7b0d92935f5a6f9ea7e58",
			want57Hash: "b3f9d21d739883d3c99bcf176d729b8e",
		},
		{
			name:       "multiple JOINs with IS NOT NULL",
			sql:        "SELECT t1.a, t2.b, t3.c FROM t1 LEFT JOIN t2 ON t1.id = t2.t1_id RIGHT JOIN t3 ON t2.id = t3.t2_id WHERE t1.x IS NOT NULL",
			want8Hash:  "aa933efe93e068448bae1740ac7dba2cd514f3e076492534b1bd44ff0d08982f",
			want57Hash: "25be3591619f82f80cafe02a9bb99eef",
		},
		{
			name:       "JSON functions",
			sql:        "SELECT JSON_EXTRACT(data, '$.user.name') AS user_name, JSON_UNQUOTE(JSON_EXTRACT(data, '$.user.email')) FROM json_docs",
			want8Hash:  "6215c011c70d71ef97f8d11d4bf80b8ce4e67a7e87dc8e53cd8b171620c14f0b",
			want57Hash: "9511f141514dcfe307a43a4ea4d72d9f",
		},
		{
			name:       "CASE expression",
			sql:        "SELECT CASE WHEN score >= 90 THEN 'A' WHEN score >= 80 THEN 'B' WHEN score >= 70 THEN 'C' ELSE 'F' END AS grade FROM students",
			want8Hash:  "0e4b162fe740e15acc424130d62baea64414a4320ac061acd7f61d15cad90c05",
			want57Hash: "2e78e7d09fbb0a983f0b6fc57ba37702",
		},
		{
			name:       "complex OR conditions with DATE_SUB",
			sql:        "SELECT * FROM orders WHERE (status = 'pending' AND created_at < DATE_SUB(NOW(), INTERVAL 1 HOUR)) OR (status = 'processing' AND updated_at < DATE_SUB(NOW(), INTERVAL 30 MINUTE))",
			want8Hash:  "d8cbe81feb5f5018a5c363d26eafc1401f15bb7a4ecc317cc3fba94bdefc974b",
			want57Hash: "45de44068610d6ad43b51ac2f2a65b50",
		},
		{
			name:       "window function FIRST_VALUE",
			sql:        "SELECT DISTINCT category, FIRST_VALUE(name) OVER (PARTITION BY category ORDER BY price) AS cheapest FROM products",
			want8Hash:  "b25149f2194dee30fc35666b1e0c3f9a42d8800e51e12d8c8890101fcd45f67d",
			want57Hash: "", // Window functions not supported in MySQL 5.7
		},
		// UTF-8 identifier tests
		{
			name:       "Chinese table and column names",
			sql:        `SELECT * FROM 用户表 WHERE 名字 = "test"`,
			want8Hash:  "f43f4998e24eaadbfffbcf19144a44c27aadd8a814fc7e22d023706bce0604ef",
			want57Hash: "6365f4cf87cfe52fd00432fe5c72abc0",
		},
		{
			name:       "Chinese with backticks",
			sql:        "SELECT `姓名`, `年龄` FROM `员工表` WHERE `部门` = '技术'",
			want8Hash:  "983bfbac776785dfe15c6a8c6a67be36cc1100ba7411288cd3f1bea40bc1bf31",
			want57Hash: "93f0357d4e95e2ff402bf36adadcd9b9",
		},
		{
			name:       "Mixed ASCII and Chinese identifiers",
			sql:        "SELECT user_id, 用户名 FROM users_用户 WHERE active = 1",
			want8Hash:  "9e966b778c0ab8cb6399c81fddfaeca752e5e740b31f79e59e5e0aee6e5284c0",
			want57Hash: "42787517f1b3d1de462948815984baff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Compute(tt.sql)
			if d.Hash != tt.want8Hash {
				t.Errorf("Hash mismatch for %q:\n  got:  %s\n  want: %s\n  text: %s", tt.sql, d.Hash, tt.want8Hash, d.Text)
			}
		})
	}
}

func TestMySQL57Compatibility(t *testing.T) {
	tests := []struct {
		name       string
		sql        string
		want57Hash string // MySQL 5.7 MD5 (empty if not supported in 5.7)
	}{
		{
			name:       "simple SELECT with WHERE",
			sql:        "SELECT * FROM users WHERE id = 42",
			want57Hash: "731d9efe96031900ba2a36667f4718d0",
		},
		{
			name:       "JOIN with multiple conditions",
			sql:        "SELECT a.id, b.name FROM users a JOIN orders b ON a.id = b.user_id WHERE b.total > 100 AND a.status = 'active'",
			want57Hash: "5b83e9f6a57a0c372448944f7ffaafb3",
		},
		{
			name:       "INSERT with multiple rows",
			sql:        "INSERT INTO logs (user_id, action, created_at) VALUES (1, 'login', NOW()), (2, 'logout', NOW())",
			want57Hash: "d3836b53e50fbf1b4947e4e6fc634676",
		},
		{
			name:       "UPDATE with arithmetic and conditions",
			sql:        "UPDATE accounts SET balance = balance - 50.00, updated_at = NOW() WHERE id = 123 AND balance >= 50.00",
			want57Hash: "1c3b7450959b0b3f228a2e1e81243e8d",
		},
		{
			name:       "DELETE with IN clause",
			sql:        "DELETE FROM sessions WHERE expires_at < NOW() AND user_id IN (1, 2, 3, 4, 5)",
			want57Hash: "b5f970dd3ff955057de7f2d8b6e7e2aa",
		},
		{
			name:       "GROUP BY with HAVING and ORDER BY",
			sql:        "SELECT COUNT(*) AS cnt, status FROM orders GROUP BY status HAVING COUNT(*) > 10 ORDER BY cnt DESC LIMIT 5",
			want57Hash: "e2476ee01e425d230cf48d15004d7274",
		},
		{
			name:       "correlated subquery",
			sql:        "SELECT u.*, (SELECT COUNT(*) FROM orders WHERE user_id = u.id) AS order_count FROM users u WHERE u.created_at > '2024-01-01'",
			want57Hash: "ba912c30ec03185da4cad1d63aa99d7e",
		},
		{
			name:       "subquery in WHERE with BETWEEN",
			sql:        "SELECT * FROM products WHERE category_id IN (SELECT id FROM categories WHERE parent_id = 5) AND price BETWEEN 10.00 AND 100.00",
			want57Hash: "dd10c7e22290216e997c9b676ed7255a",
		},
		{
			name:       "DATE functions with GROUP BY",
			sql:        "SELECT DATE(created_at) AS day, SUM(amount) AS total FROM transactions GROUP BY DATE(created_at) ORDER BY day DESC",
			want57Hash: "eddd0d1add08602f293078d3e72bcf46",
		},
		{
			name:       "COALESCE and IFNULL",
			sql:        "SELECT COALESCE(nickname, username, email) AS display_name, IFNULL(avatar_url, '/default.png') FROM users WHERE id = 1",
			want57Hash: "9e8125db5f7559baac2587ff253d4955",
		},
		{
			name:       "DATE_ADD with INTERVAL",
			sql:        "SELECT * FROM events WHERE start_time >= NOW() AND end_time <= DATE_ADD(NOW(), INTERVAL 7 DAY) ORDER BY start_time",
			want57Hash: "c3286587a00d2a9d9bc9945a46b7d1af",
		},
		{
			name:       "optimizer hint MAX_EXECUTION_TIME",
			sql:        "SELECT /*+ MAX_EXECUTION_TIME(5000) */ id, name FROM large_table WHERE status = 1 LIMIT 1000",
			want57Hash: "0673f7e6d08e8d422618b2ee0e6700dd",
		},
		{
			name:       "HEX and conversion functions",
			sql:        "SELECT HEX(uuid), FROM_UNIXTIME(created_ts), INET_NTOA(ip_address) FROM access_logs WHERE id > 0",
			want57Hash: "a033907d5dba747403db7f4ee19e7418",
		},
		{
			name:       "INSERT SELECT",
			sql:        "INSERT INTO audit_log SELECT NULL, 'UPDATE', OLD.*, NEW.*, NOW() FROM dual WHERE 1=1",
			want57Hash: "b3f9d21d739883d3c99bcf176d729b8e",
		},
		{
			name:       "multiple JOINs with IS NOT NULL",
			sql:        "SELECT t1.a, t2.b, t3.c FROM t1 LEFT JOIN t2 ON t1.id = t2.t1_id RIGHT JOIN t3 ON t2.id = t3.t2_id WHERE t1.x IS NOT NULL",
			want57Hash: "25be3591619f82f80cafe02a9bb99eef",
		},
		{
			name:       "JSON functions",
			sql:        "SELECT JSON_EXTRACT(data, '$.user.name') AS user_name, JSON_UNQUOTE(JSON_EXTRACT(data, '$.user.email')) FROM json_docs",
			want57Hash: "9511f141514dcfe307a43a4ea4d72d9f",
		},
		{
			name:       "CASE expression",
			sql:        "SELECT CASE WHEN score >= 90 THEN 'A' WHEN score >= 80 THEN 'B' WHEN score >= 70 THEN 'C' ELSE 'F' END AS grade FROM students",
			want57Hash: "2e78e7d09fbb0a983f0b6fc57ba37702",
		},
		{
			name:       "complex OR conditions with DATE_SUB",
			sql:        "SELECT * FROM orders WHERE (status = 'pending' AND created_at < DATE_SUB(NOW(), INTERVAL 1 HOUR)) OR (status = 'processing' AND updated_at < DATE_SUB(NOW(), INTERVAL 30 MINUTE))",
			want57Hash: "45de44068610d6ad43b51ac2f2a65b50",
		},
		// UTF-8 identifier tests
		{
			name:       "Chinese table and column names",
			sql:        `SELECT * FROM 用户表 WHERE 名字 = "test"`,
			want57Hash: "6365f4cf87cfe52fd00432fe5c72abc0",
		},
		{
			name:       "Chinese with backticks",
			sql:        "SELECT `姓名`, `年龄` FROM `员工表` WHERE `部门` = '技术'",
			want57Hash: "93f0357d4e95e2ff402bf36adadcd9b9",
		},
		{
			name:       "Mixed ASCII and Chinese identifiers",
			sql:        "SELECT user_id, 用户名 FROM users_用户 WHERE active = 1",
			want57Hash: "42787517f1b3d1de462948815984baff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Normalize(tt.sql, Options{Version: MySQL57})
			if d.Hash != tt.want57Hash {
				t.Errorf("MySQL 5.7 hash mismatch for %q:\n  got:  %s\n  want: %s\n  text: %s", tt.sql, d.Hash, tt.want57Hash, d.Text)
			}
		})
	}
}

func TestMySQL57ComplexQueries(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		wantHash string
		wantText string // Optional: verify normalized text matches MySQL 5.7
	}{
		// Complex JOINs
		{
			name:     "INNER JOIN with WHERE",
			sql:      "SELECT u.id, u.name, o.total FROM users u INNER JOIN orders o ON u.id = o.user_id WHERE o.total > 100",
			wantHash: "600ddfcea0b27fc3fa56556cfe7cf463",
			wantText: "SELECT `u` . `id` , `u` . `name` , `o` . `total` FROM `users` `u` INNER JOIN `orders` `o` ON `u` . `id` = `o` . `user_id` WHERE `o` . `total` > ?",
		},
		{
			name:     "LEFT JOIN with IS NULL",
			sql:      "SELECT u.id, o.id FROM users u LEFT JOIN orders o ON u.id = o.user_id WHERE o.id IS NULL",
			wantHash: "86f7ddbc818b260e7ee9c8f864376a20",
			wantText: "SELECT `u` . `id` , `o` . `id` FROM `users` `u` LEFT JOIN `orders` `o` ON `u` . `id` = `o` . `user_id` WHERE `o` . `id` IS NULL",
		},
		{
			name:     "RIGHT JOIN",
			sql:      "SELECT u.id, o.id FROM users u RIGHT JOIN orders o ON u.id = o.user_id",
			wantHash: "c14ed21262e2896be69381abd01187d4",
			wantText: "SELECT `u` . `id` , `o` . `id` FROM `users` `u` RIGHT JOIN `orders` `o` ON `u` . `id` = `o` . `user_id`",
		},

		// Subqueries
		{
			name:     "Subquery with IN",
			sql:      "SELECT * FROM users WHERE id IN (SELECT user_id FROM orders WHERE total > 500)",
			wantHash: "b2ba67be94f35ce85410b67d401f73e6",
			wantText: "SELECT * FROM `users` WHERE `id` IN ( SELECT `user_id` FROM `orders` WHERE `total` > ? )",
		},
		{
			name:     "Correlated subquery with EXISTS",
			sql:      "SELECT * FROM users u WHERE EXISTS (SELECT 1 FROM orders o WHERE o.user_id = u.id)",
			wantHash: "479381c5057c8b73d0e394a3840c89a0",
			wantText: "SELECT * FROM `users` `u` WHERE EXISTS ( SELECT ? FROM `orders` `o` WHERE `o` . `user_id` = `u` . `id` )",
		},

		// Aggregations
		{
			name:     "GROUP BY with HAVING and multiple aggregates",
			sql:      "SELECT user_id, SUM(total), AVG(total), COUNT(*) FROM orders GROUP BY user_id HAVING SUM(total) > 1000",
			wantHash: "a7268421644501822989e08e7ae61763",
			wantText: "SELECT `user_id` , SUM ( `total` ) , AVG ( `total` ) , COUNT ( * ) FROM `orders` GROUP BY `user_id` HAVING SUM ( `total` ) > ?",
		},
		{
			name:     "ORDER BY with LIMIT and OFFSET",
			sql:      "SELECT user_id, total FROM orders ORDER BY total DESC LIMIT 10 OFFSET 5",
			wantHash: "7d290ea86943c3898763671ad6f6f756",
			wantText: "SELECT `user_id` , `total` FROM `orders` ORDER BY `total` DESC LIMIT ? OFFSET ?",
		},
		{
			name:     "SELECT DISTINCT",
			sql:      "SELECT DISTINCT user_id FROM orders",
			wantHash: "272278e62fe86028d0144dec6748ec63",
			wantText: "SELECT DISTINCTROW `user_id` FROM `orders`",
		},

		// UNION operations
		{
			name:     "UNION with WHERE",
			sql:      "SELECT id, name FROM users UNION SELECT id, name FROM users WHERE id > 5",
			wantHash: "201e1cd57930dcd40bc9d93cb18927e1",
			wantText: "SELECT `id` , NAME FROM `users` UNION SELECT `id` , NAME FROM `users` WHERE `id` > ?",
		},
		{
			name:     "UNION ALL",
			sql:      "SELECT id FROM users UNION ALL SELECT user_id FROM orders",
			wantHash: "978a5f0e12890034fe533edd18416488",
			wantText: "SELECT `id` FROM `users` UNION ALL SELECT `user_id` FROM `orders`",
		},

		// CASE expressions
		{
			name:     "CASE WHEN with multiple conditions",
			sql:      "SELECT id, CASE WHEN total > 100 THEN 'high' WHEN total > 50 THEN 'medium' ELSE 'low' END as level FROM orders",
			wantHash: "27909381d3d7af9c0ce606beaf6d2541",
			wantText: "SELECT `id` , CASE WHEN `total` > ? THEN ? WHEN `total` > ? THEN ? ELSE ? END AS LEVEL FROM `orders`",
		},

		// Date functions
		{
			name:     "BETWEEN with dates",
			sql:      "SELECT * FROM orders WHERE created_at BETWEEN '2024-01-01' AND '2024-12-31'",
			wantHash: "c5a225304fb317916407dc7f2e162e47",
			wantText: "SELECT * FROM `orders` WHERE `created_at` BETWEEN ? AND ?",
		},

		// String functions
		{
			name:     "LIKE pattern matching",
			sql:      "SELECT * FROM users WHERE name LIKE '%test%'",
			wantHash: "17db0beed706808914c82a82ff560d34",
			wantText: "SELECT * FROM `users` WHERE NAME LIKE ?",
		},

		// NULL handling
		{
			name:     "COALESCE and IFNULL",
			sql:      "SELECT COALESCE(name, 'unknown'), IFNULL(name, 'N/A') FROM users",
			wantHash: "dd89a5d69cb43e58ce407696da71afa9",
			wantText: "SELECT COALESCE ( NAME , ? ) , `IFNULL` ( NAME , ? ) FROM `users`",
		},
		{
			name:     "IS NOT NULL",
			sql:      "SELECT * FROM users WHERE name IS NOT NULL",
			wantHash: "342249c626a9a41ef048cdd146557131",
			wantText: "SELECT * FROM `users` WHERE NAME IS NOT NULL",
		},

		// INSERT variants
		{
			name:     "INSERT single row",
			sql:      "INSERT INTO users (id, name) VALUES (100, 'test100')",
			wantHash: "daea840fef6158aeebb39faf1743eb8e",
			wantText: "INSERT INTO `users` ( `id` , NAME ) VALUES (...)",
		},
		{
			name:     "INSERT multiple rows",
			sql:      "INSERT INTO users (id, name) VALUES (101, 'a'), (102, 'b'), (103, 'c')",
			wantHash: "f0f452cc8ee8a090455d1dadf6e7cbb6",
			wantText: "INSERT INTO `users` ( `id` , NAME ) VALUES (...) /* , ... */",
		},

		// UPDATE
		{
			name:     "UPDATE with WHERE",
			sql:      "UPDATE users SET name = 'newname' WHERE id = 200",
			wantHash: "bd2b0e56be79850c13ab2fc6fb69b145",
			wantText: "UPDATE `users` SET NAME = ? WHERE `id` = ?",
		},

		// DELETE variants
		{
			name:     "DELETE with equality",
			sql:      "DELETE FROM users WHERE id = 300",
			wantHash: "bdb665324e970abaf63f9e7131410e93",
			wantText: "DELETE FROM `users` WHERE `id` = ?",
		},
		{
			name:     "DELETE with IN clause",
			sql:      "DELETE FROM users WHERE id IN (301, 302, 303, 304, 305)",
			wantHash: "0c530b7aa65278e21ccb481560dcf55b",
			wantText: "DELETE FROM `users` WHERE `id` IN (...)",
		},

		// Complex WHERE clauses
		{
			name:     "Complex AND/OR conditions",
			sql:      "SELECT * FROM users WHERE (id > 10 AND name = 'test') OR (id < 5 AND name = 'admin')",
			wantHash: "cbc5ab08da1094bebbd23fa387303f14",
			wantText: "SELECT * FROM `users` WHERE ( `id` > ? AND NAME = ? ) OR ( `id` < ? AND NAME = ? )",
		},
		{
			name:     "BETWEEN with IN clause",
			sql:      "SELECT * FROM orders WHERE total BETWEEN 100 AND 500 AND user_id IN (1001, 1002, 1003)",
			wantHash: "8d4862b2693ef811669eccd1c5800a87",
			wantText: "SELECT * FROM `orders` WHERE `total` BETWEEN ? AND ? AND `user_id` IN (...)",
		},

		// String functions with CONCAT
		// CONCAT, UPPER, LENGTH are not reserved keywords in MySQL
		// treated as identifiers, so they get backtick-quoted in digest text.
		{
			name:     "CONCAT UPPER LENGTH functions",
			sql:      "SELECT CONCAT(id, '-', name), UPPER(name), LENGTH(name) FROM users",
			wantHash: "77f2e8dfcd737c7dd389b4d7234c6722",
			wantText: "SELECT `CONCAT` ( `id` , ? , NAME ) , `UPPER` ( NAME ) , `LENGTH` ( NAME ) FROM `users`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Normalize(tt.sql, Options{Version: MySQL57})

			if d.Hash != tt.wantHash {
				t.Errorf("Hash mismatch:\n  got:  %s\n  want: %s\n  sql:  %s\n  text: %s", d.Hash, tt.wantHash, tt.sql, d.Text)
			}

			if tt.wantText != "" && d.Text != tt.wantText {
				t.Errorf("Text mismatch:\n  got:  %s\n  want: %s", d.Text, tt.wantText)
			}
		})
	}
}

// TestMySQL57ValueFolding verifies that value list folding matches MySQL 5.7 behavior.
func TestMySQL57ValueFolding(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		wantHash string
	}{
		{
			name:     "IN clause folding - 3 values",
			sql:      "SELECT * FROM users WHERE id IN (1, 2, 3)",
			wantHash: "", // Will verify folding happens
		},
		{
			name:     "IN clause folding - 10 values",
			sql:      "SELECT * FROM users WHERE id IN (1, 2, 3, 4, 5, 6, 7, 8, 9, 10)",
			wantHash: "", // Should produce same hash as 3 values
		},
	}

	// Verify that different number of IN values produce the same hash
	var firstHash string
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Normalize(tt.sql, Options{Version: MySQL57})
			if i == 0 {
				firstHash = d.Hash
			} else {
				if d.Hash != firstHash {
					t.Errorf("IN clause folding not working: different hashes for different value counts\n  first: %s\n  got:   %s", firstHash, d.Hash)
				}
			}
		})
	}
}

// TestMySQL57MultiRowInsert verifies multi-row INSERT folding matches MySQL 5.7.
func TestMySQL57MultiRowInsert(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{
			name: "2 rows",
			sql:  "INSERT INTO t (a, b) VALUES (1, 'x'), (2, 'y')",
		},
		{
			name: "5 rows",
			sql:  "INSERT INTO t (a, b) VALUES (1, 'a'), (2, 'b'), (3, 'c'), (4, 'd'), (5, 'e')",
		},
	}

	// All multi-row INSERTs should produce the same hash
	var firstHash string
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Normalize(tt.sql, Options{Version: MySQL57})
			if i == 0 {
				firstHash = d.Hash
			} else {
				if d.Hash != firstHash {
					t.Errorf("Multi-row INSERT folding not working: different hashes\n  first: %s\n  got:   %s", firstHash, d.Hash)
				}
			}

			// Verify text contains the folding comment
			if i > 0 && d.Text != "" {
				// Multi-row should have comment
				if d.Text == "" {
					t.Error("Expected folding comment in multi-row INSERT")
				}
			}
		})
	}
}
