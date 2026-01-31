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
    digest "github.com/rashiq/mysql-digest"
)

func main() {
    // Simple usage
    result := digest.Compute("SELECT * FROM users WHERE id = 123")
    fmt.Println(result.Hash) // SHA-256 hash
    fmt.Println(result.Text) // SELECT * FROM `users` WHERE `id` = ?

    // With options
    result = digest.Compute("SELECT * FROM users WHERE id = 123", digest.Options{
        Version:   digest.MySQL57,  // Produces MD5 hash
        SQLMode:   digest.AnsiQuotes, // Enable ANSI_QUOTES mode
    })
}
```

### CLI

```bash
# From command line argument
mysql-digest -sql "SELECT * FROM users WHERE id = 123"

# From file
mysql-digest -file query.sql

# From stdin
echo "SELECT * FROM users WHERE id = 123" | mysql-digest

# Output formats
mysql-digest -sql "SELECT 1" -json
mysql-digest -sql "SELECT 1" -hash-only
mysql-digest -sql "SELECT 1" -text-only

# Debug mode (show lexer tokens)
mysql-digest -sql "SELECT 1" -debug
```

**Example output:**

```
DIGEST: 5b79e33af9e9...
DIGEST_TEXT: SELECT * FROM `users` WHERE `id` = ?
```

## License

MIT License - see [LICENSE](LICENSE) file.
