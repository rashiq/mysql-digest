# mysql-digest

A library for computing unique fingerprints for MySQL queries, matching the one in MySQL's Performance Schema. 

It reimplements MySQL's sql lexer to accurately normalize queries.



## Installation

### Library

```bash
go get github.com/rashiq/mysql-digest
```

### CLI

```bash
go install github.com/rashiq/mysql-digest/cmd/mysql-digest@latest
```

## Usage

### Library

```go
package main

import (
    "fmt"
    digest "github.com/rashiq/mysql-digest"
)

func main() {
    result, _ := digest.Compute("SELECT * FROM users WHERE id = 123")
    fmt.Println(result.Hash) // SHA-256 hash
    fmt.Println(result.Text) // SELECT * FROM `users` WHERE `id` = ?

    // With options
    result, _ = digest.Compute("SELECT * FROM users WHERE id = 123", digest.Options{
        Version: digest.MySQL57, // Produces MD5 hash
        SQLMode: digest.MODE_ANSI_QUOTES,
    })

    d := digest.NewDigester(digest.Options{Version: digest.MySQL84})
    d.Digest("SELECT * FROM t WHERE id = 1")
}
```

### CLI

```bash
mysql-digest "SELECT * FROM users WHERE id = 123"

# From file
mysql-digest -f query.sql

# From stdin
echo "SELECT * FROM users WHERE id = 123" | mysql-digest

# Output formats
mysql-digest "SELECT 1" --json
mysql-digest "SELECT 1" --hash-only
mysql-digest "SELECT 1" --text-only
```

**Example output:**

```
DIGEST: 840a880ebd1642e8a0c4926cfbaf7d4da9616b03025a080fafd43a732800fab5
DIGEST_TEXT: SELECT * FROM `users` WHERE `id` = ?
```

## Behavior

- The default version is MySQL 8.0. Set `Options.Version` to select 5.7, 8.4, or 9.0.
- The version comment thresholds are 5.7.0, 8.0.0, 8.4.0, and 9.0.0. Patch-level thresholds are not configurable.
- Unsupported versions and SQL modes return an error. Always check the error before you use a result.
- Parameter markers (`?`) and literal values produce the same digest.
- A `Digester` has no mutable query state. You can reuse it across calls and goroutines.
- Digest text is display output. Collapsed lists such as `IN (...)` cannot reproduce the original hash.
- `MaxLength` limits display bytes, including `...`, without a partial UTF-8 character. Nonpositive values disable this limit.
- `MaxLength` does not change the hash or limit input processing. MySQL token-buffer truncation is not implemented.

## License

MIT License - see [LICENSE](LICENSE) file.
