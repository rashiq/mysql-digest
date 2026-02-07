# mysql-digest

A library for computing MySQL query digests, matching the one in MySQL's Performance Schema. 

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
    "log"
    digest "github.com/rashiq/mysql-digest"
)

func main() {
    result, err := digest.Compute("SELECT * FROM users WHERE id = 123")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(result.Hash) // SHA-256 hash
    fmt.Println(result.Text) // SELECT * FROM `users` WHERE `id` = ?

    // With options
    result, err = digest.Compute("SELECT * FROM users WHERE id = 123", digest.Options{
        Version: digest.MySQL57, // Produces MD5 hash
        SQLMode: digest.MODE_ANSI_QUOTES,
    })
    if err != nil {
        log.Fatal(err)
    }

    d := digest.NewDigester(digest.Options{Version: digest.MySQL84})
    r1, _ := d.Digest("SELECT * FROM t WHERE id = 1")
    r2, _ := d.Digest("INSERT INTO t VALUES (1, 2)")
    fmt.Println(r1.Hash, r2.Hash)
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

## License

MIT License - see [LICENSE](LICENSE) file.
