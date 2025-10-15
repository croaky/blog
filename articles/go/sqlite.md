# go / sqlite

I use SQLite in single-process Go web servers
when I have modest traffic requirements
and operational simplicity matters more than other factors.

## Setup

I use [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) as
my database driver.
It is a pure Go implementation with no CGo dependencies,
reducing cross-compilation issues.

```bash
go mod init server
go get modernc.org/sqlite
```

I configure SQLite for single-process access to avoid
`database is locked (5) (SQLITE_BUSY)` errors,
following David Crawshaw's
[one process programming notes](https://crawshaw.io/blog/one-process-programming-notes).

```go
import (
    "database/sql"
    "fmt"
    "log"
    "net/http"

    _ "modernc.org/sqlite"
)

func initDB(db *sql.DB) error {
    pragmas := []string{
        "PRAGMA temp_store = memory",    // Store temp tables in memory
        "PRAGMA mmap_size  = 268435456", // 256MB memory-mapped I/O
        "PRAGMA cache_size = 10000",     // Cache size in pages
    }

    for _, pragma := range pragmas {
        if _, err := db.Exec(pragma); err != nil {
            return fmt.Errorf("failed to set pragma %s: %w", pragma, err)
        }
    }

    return nil
}

type Server struct {
    db *sql.DB
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
    var result int
    if err := s.db.QueryRow("SELECT 1").Scan(&result); err != nil {
        http.Error(w, "Database error", 500)
        return
    }
    w.Write([]byte("OK"))
}

func main() {
    conn := "app.db?" +
        "_busy_timeout=5000&" +  // Avoid immediate lock failures in concurrent access
        "_journal_mode=WAL&" +   // Better concurrency than default rollback journal
        "_synchronous=NORMAL&" + // Faster writes while maintaining crash safety
        "_foreign_keys=ON"       // Enforce referential integrity
    db, err := sql.Open("sqlite", conn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // For a single-process web server,
    // limit connections to prevent lock contention
    db.SetMaxOpenConns(1)
    db.SetMaxIdleConns(1)
    db.SetConnMaxLifetime(0) // Keep connections alive

    if err := initDB(db); err != nil {
        log.Fatal("Failed to initialize database:", err)
    }

    server := &Server{db: db}
    http.HandleFunc("/health", server.health)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Testing

I use in-memory databases to isolate each
test case.

```go
import (
    "database/sql"
    "net/http"
    "net/http/httptest"
    "testing"

    _ "modernc.org/sqlite"
)

func initTestDB(t *testing.T) (*sql.DB, *Server) {
    t.Helper()

    db, err := sql.Open("sqlite", ":memory:?_foreign_keys=ON")
    if err != nil {
        t.Fatalf("Failed to create test database: %v", err)
    }

    if err := initDB(db); err != nil {
        db.Close()
        t.Fatalf("Failed to initialize test database: %v", err)
    }

    server := &Server{db: db}
    return db, server
}

func TestHealthCheck(t *testing.T) {
    db, server := initTestDB(t)
    defer db.Close()

    req, err := http.NewRequest("GET", "/health", nil)
    if err != nil {
        t.Fatalf("Failed to create request: %v", err)
    }

    rr := httptest.NewRecorder()
    http.HandlerFunc(server.health).ServeHTTP(rr, req)

    if rr.Code != 200 {
        t.Errorf("Expected status 200, got %d", rr.Code)
    }

    if rr.Body.String() != "OK" {
        t.Errorf("Expected body 'OK', got '%s'", rr.Body.String())
    }
}
```
