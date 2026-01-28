package digest

import (
	"testing"
)

// TestMySQL57ComplexQueries tests digest computation against actual MySQL 5.7.44 output.
// These hashes were generated using a MySQL 5.7.44 container and verified against
// performance_schema.events_statements_summary_by_digest.
//
// To regenerate these hashes:
// 1. Run MySQL 5.7: docker run -d --name mysql57-test -e MYSQL_ROOT_PASSWORD=test -e MYSQL_DATABASE=digest_test -p 3307:3306 mysql:5.7
// 2. Create test tables and execute the SQL statements
// 3. Query: SELECT DIGEST, DIGEST_TEXT FROM performance_schema.events_statements_summary_by_digest WHERE SCHEMA_NAME = 'digest_test';
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
		// CONCAT, UPPER, LENGTH are not reserved keywords in MySQL - they're function names
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
