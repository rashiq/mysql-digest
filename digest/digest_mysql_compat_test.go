package digest

import (
	"testing"
)

// TestMySQLCompatibility tests that our digest computation matches MySQL 8.0's STATEMENT_DIGEST() function.
// These hashes were generated using MySQL 8.0.34 with:
//
//	SELECT STATEMENT_DIGEST('...');
//
// MySQL 5.7 hashes were obtained from performance_schema.events_statements_summary_by_digest
// Note: MySQL 5.7 uses MD5 (32 hex chars), MySQL 8.0 uses SHA-256 (64 hex chars)
// Note: MySQL 5.7 doesn't support CTEs or window functions, those tests have empty want57Hash
func TestMySQLCompatibility(t *testing.T) {
	tests := []struct {
		name       string
		sql        string
		wantHash   string // MySQL 8.0 SHA-256
		want57Hash string // MySQL 5.7 MD5 (empty if not supported in 5.7)
	}{
		{
			name:       "simple SELECT with WHERE",
			sql:        "SELECT * FROM users WHERE id = 42",
			wantHash:   "840a880ebd1642e8a0c4926cfbaf7d4da9616b03025a080fafd43a732800fab5",
			want57Hash: "731d9efe96031900ba2a36667f4718d0",
		},
		{
			name:       "JOIN with multiple conditions",
			sql:        "SELECT a.id, b.name FROM users a JOIN orders b ON a.id = b.user_id WHERE b.total > 100 AND a.status = 'active'",
			wantHash:   "bbbad1c572514399449094ee3043cd0c8791dd0ae66d579066f05e5f6a2a2bd6",
			want57Hash: "5b83e9f6a57a0c372448944f7ffaafb3",
		},
		{
			name:       "INSERT with multiple rows",
			sql:        "INSERT INTO logs (user_id, action, created_at) VALUES (1, 'login', NOW()), (2, 'logout', NOW())",
			wantHash:   "63082d4181bdc6fada36b20f7719884355aef4dd40c26070d3377d0bc66bf84a",
			want57Hash: "d3836b53e50fbf1b4947e4e6fc634676",
		},
		{
			name:       "UPDATE with arithmetic and conditions",
			sql:        "UPDATE accounts SET balance = balance - 50.00, updated_at = NOW() WHERE id = 123 AND balance >= 50.00",
			wantHash:   "989b4724da70672cb9ee411beea017d2140f57471e708bd566ad97502e6046a7",
			want57Hash: "1c3b7450959b0b3f228a2e1e81243e8d",
		},
		{
			name:       "DELETE with IN clause",
			sql:        "DELETE FROM sessions WHERE expires_at < NOW() AND user_id IN (1, 2, 3, 4, 5)",
			wantHash:   "6490f2bee2884432c4fe1b3f63298d502ba101bdcc1ef8e02f8e5754dd37cf70",
			want57Hash: "b5f970dd3ff955057de7f2d8b6e7e2aa",
		},
		{
			name:       "GROUP BY with HAVING and ORDER BY",
			sql:        "SELECT COUNT(*) AS cnt, status FROM orders GROUP BY status HAVING COUNT(*) > 10 ORDER BY cnt DESC LIMIT 5",
			wantHash:   "a61c184788ca4ab687dbced1e7c4c2f86b4d3a0b6324bd21a812d7c7c466fcdd",
			want57Hash: "e2476ee01e425d230cf48d15004d7274",
		},
		{
			name:       "correlated subquery",
			sql:        "SELECT u.*, (SELECT COUNT(*) FROM orders WHERE user_id = u.id) AS order_count FROM users u WHERE u.created_at > '2024-01-01'",
			wantHash:   "7f476b144ccca78b057630a5339080e6baa027cf025d289749700ba97dbd8f37",
			want57Hash: "ba912c30ec03185da4cad1d63aa99d7e",
		},
		{
			name:       "subquery in WHERE with BETWEEN",
			sql:        "SELECT * FROM products WHERE category_id IN (SELECT id FROM categories WHERE parent_id = 5) AND price BETWEEN 10.00 AND 100.00",
			wantHash:   "fb9714785b2c028cf58ebb6d4b1ee18285f453c7c6a8f0f469c0f9b5222df78d",
			want57Hash: "dd10c7e22290216e997c9b676ed7255a",
		},
		{
			name:       "DATE functions with GROUP BY",
			sql:        "SELECT DATE(created_at) AS day, SUM(amount) AS total FROM transactions GROUP BY DATE(created_at) ORDER BY day DESC",
			wantHash:   "d60a728c07bec79fceae285c6406323ef564dfd39401e5b87f82de73e8595791",
			want57Hash: "eddd0d1add08602f293078d3e72bcf46",
		},
		{
			name:       "COALESCE and IFNULL",
			sql:        "SELECT COALESCE(nickname, username, email) AS display_name, IFNULL(avatar_url, '/default.png') FROM users WHERE id = 1",
			wantHash:   "a27ab6578509d6aa9faa6620e08ce05d272bd7e0208ceba693d07fc71f168022",
			want57Hash: "9e8125db5f7559baac2587ff253d4955",
		},
		{
			name:       "DATE_ADD with INTERVAL",
			sql:        "SELECT * FROM events WHERE start_time >= NOW() AND end_time <= DATE_ADD(NOW(), INTERVAL 7 DAY) ORDER BY start_time",
			wantHash:   "15d88f9a89c7bc0b2f698086b6379933ad962dd113d079a20cee53b7756fa541",
			want57Hash: "c3286587a00d2a9d9bc9945a46b7d1af",
		},
		{
			name:       "CTE with window function",
			sql:        "WITH ranked AS (SELECT *, ROW_NUMBER() OVER (PARTITION BY category_id ORDER BY price DESC) AS rn FROM products) SELECT * FROM ranked WHERE rn <= 3",
			wantHash:   "8f35aa1c4878b88cf8d53bdb222fb77133b20b834defeee1ea1adc627b8e19fd",
			want57Hash: "", // CTEs not supported in MySQL 5.7
		},
		{
			name:       "optimizer hint MAX_EXECUTION_TIME",
			sql:        "SELECT /*+ MAX_EXECUTION_TIME(5000) */ id, name FROM large_table WHERE status = 1 LIMIT 1000",
			wantHash:   "e3399d128b66b9c8691dc9294867dfbbc7c91bd7522dcd1bf01866b8c7b0a575",
			want57Hash: "0673f7e6d08e8d422618b2ee0e6700dd",
		},
		{
			name:       "HEX and conversion functions",
			sql:        "SELECT HEX(uuid), FROM_UNIXTIME(created_ts), INET_NTOA(ip_address) FROM access_logs WHERE id > 0",
			wantHash:   "befa66430b18af4788665af06f42aa0af8698569854b1cd3fcce31430ef40b72",
			want57Hash: "a033907d5dba747403db7f4ee19e7418",
		},
		{
			name:       "INSERT SELECT",
			sql:        "INSERT INTO audit_log SELECT NULL, 'UPDATE', OLD.*, NEW.*, NOW() FROM dual WHERE 1=1",
			wantHash:   "343e3e20a203533e1351d4c11c7cac0bb2b5dafc10d7b0d92935f5a6f9ea7e58",
			want57Hash: "b3f9d21d739883d3c99bcf176d729b8e",
		},
		{
			name:       "multiple JOINs with IS NOT NULL",
			sql:        "SELECT t1.a, t2.b, t3.c FROM t1 LEFT JOIN t2 ON t1.id = t2.t1_id RIGHT JOIN t3 ON t2.id = t3.t2_id WHERE t1.x IS NOT NULL",
			wantHash:   "aa933efe93e068448bae1740ac7dba2cd514f3e076492534b1bd44ff0d08982f",
			want57Hash: "25be3591619f82f80cafe02a9bb99eef",
		},
		{
			name:       "JSON functions",
			sql:        "SELECT JSON_EXTRACT(data, '$.user.name') AS user_name, JSON_UNQUOTE(JSON_EXTRACT(data, '$.user.email')) FROM json_docs",
			wantHash:   "6215c011c70d71ef97f8d11d4bf80b8ce4e67a7e87dc8e53cd8b171620c14f0b",
			want57Hash: "9511f141514dcfe307a43a4ea4d72d9f",
		},
		{
			name:       "CASE expression",
			sql:        "SELECT CASE WHEN score >= 90 THEN 'A' WHEN score >= 80 THEN 'B' WHEN score >= 70 THEN 'C' ELSE 'F' END AS grade FROM students",
			wantHash:   "0e4b162fe740e15acc424130d62baea64414a4320ac061acd7f61d15cad90c05",
			want57Hash: "2e78e7d09fbb0a983f0b6fc57ba37702",
		},
		{
			name:       "complex OR conditions with DATE_SUB",
			sql:        "SELECT * FROM orders WHERE (status = 'pending' AND created_at < DATE_SUB(NOW(), INTERVAL 1 HOUR)) OR (status = 'processing' AND updated_at < DATE_SUB(NOW(), INTERVAL 30 MINUTE))",
			wantHash:   "d8cbe81feb5f5018a5c363d26eafc1401f15bb7a4ecc317cc3fba94bdefc974b",
			want57Hash: "45de44068610d6ad43b51ac2f2a65b50",
		},
		{
			name:       "window function FIRST_VALUE",
			sql:        "SELECT DISTINCT category, FIRST_VALUE(name) OVER (PARTITION BY category ORDER BY price) AS cheapest FROM products",
			wantHash:   "b25149f2194dee30fc35666b1e0c3f9a42d8800e51e12d8c8890101fcd45f67d",
			want57Hash: "", // Window functions not supported in MySQL 5.7
		},
		// UTF-8 identifier tests
		{
			name:       "Chinese table and column names",
			sql:        `SELECT * FROM 用户表 WHERE 名字 = "test"`,
			wantHash:   "f43f4998e24eaadbfffbcf19144a44c27aadd8a814fc7e22d023706bce0604ef",
			want57Hash: "6365f4cf87cfe52fd00432fe5c72abc0",
		},
		{
			name:       "Chinese with backticks",
			sql:        "SELECT `姓名`, `年龄` FROM `员工表` WHERE `部门` = '技术'",
			wantHash:   "983bfbac776785dfe15c6a8c6a67be36cc1100ba7411288cd3f1bea40bc1bf31",
			want57Hash: "93f0357d4e95e2ff402bf36adadcd9b9",
		},
		{
			name:       "Mixed ASCII and Chinese identifiers",
			sql:        "SELECT user_id, 用户名 FROM users_用户 WHERE active = 1",
			wantHash:   "9e966b778c0ab8cb6399c81fddfaeca752e5e740b31f79e59e5e0aee6e5284c0",
			want57Hash: "42787517f1b3d1de462948815984baff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Compute(tt.sql)
			if d.Hash != tt.wantHash {
				t.Errorf("Hash mismatch for %q:\n  got:  %s\n  want: %s\n  text: %s", tt.sql, d.Hash, tt.wantHash, d.Text)
			}
		})
	}
}

// TestMySQL57Compatibility tests that our digest computation matches MySQL 5.7's digest algorithm.
// MySQL 5.7 uses MD5 hashing instead of SHA-256.
// Note: MySQL 5.7 doesn't support CTEs or window functions, those tests are skipped.
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
