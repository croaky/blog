# go / sqlite

I use [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite)
for SQLite in Go web servers.
It's a pure Go implementation with no CGo dependencies,
eliminating cross-compilation issues.

## Setup

I configure SQLite with Write-Ahead Logging (WAL) mode
and a single connection to prevent lock contention:

```go
import (
    "database/sql"
    _ "modernc.org/sqlite"
)

db, err := sql.Open("sqlite", "app.db?_busy_timeout=5000&_journal_mode=WAL&_synchronous=NORMAL&_foreign_keys=ON")
if err != nil {
    log.Fatal(err)
}
defer db.Close()

db.SetMaxOpenConns(1)
db.SetMaxIdleConns(1)
db.SetConnMaxLifetime(0)
```

Connection string parameters:

- `_busy_timeout=5000`: Wait 5 seconds for locks before timing out
- `_journal_mode=WAL`: Enable Write-Ahead Logging for better concurrency
- `_synchronous=NORMAL`: Balance durability and performance
- `_foreign_keys=ON`: Enable foreign key constraints

I set performance pragmas on initialization:

```go
func initDB(db *sql.DB) error {
    pragmas := []string{
        "PRAGMA temp_store = memory",   // Store temp tables in memory
        "PRAGMA mmap_size = 268435456", // 256MB memory-mapped I/O
        "PRAGMA cache_size = 10000",    // Cache size in pages
    }

    for _, pragma := range pragmas {
        if _, err := db.Exec(pragma); err != nil {
            return fmt.Errorf("failed to set pragma %s: %w", pragma, err)
        }
    }
    return nil
}
```

## Basic operations

I use parameterized queries with `?` placeholders:

```go
// Insert
result, err := db.Exec("INSERT INTO notes (content) VALUES (?)", content)
if err != nil {
    return err
}
id, err := result.LastInsertId()

// Query
rows, err := db.Query("SELECT id, content FROM notes WHERE id = ?", id)
if err != nil {
    return err
}
defer rows.Close()

// Update
_, err = db.Exec("UPDATE notes SET content = ? WHERE id = ?", newContent, id)

// Delete
_, err = db.Exec("DELETE FROM notes WHERE id = ?", id)
```

## Testing

I use in-memory databases for tests:

```go
func setupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("sqlite", ":memory:")
    if err != nil {
        t.Fatalf("Failed to create test database: %v", err)
    }

    if err := initDB(db); err != nil {
        t.Fatalf("Failed to initialize test database: %v", err)
    }

    return db
}
```

In-memory databases are isolated per connection and
automatically cleaned up when closed.
