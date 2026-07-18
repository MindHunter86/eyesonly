package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"time"
)

//go:embed migrations/*/*.sql
var migrationFiles embed.FS

func Migrate(ctx context.Context, database *sql.DB, driver string) (e error) {
	driver = normalizeMigrationDriver(driver)
	if _, e := database.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (name VARCHAR(255) PRIMARY KEY, applied_at VARCHAR(64) NOT NULL)`); e != nil {
		return e
	}

	root := "migrations/" + driver
	var entries []fs.DirEntry
	if entries, e = fs.ReadDir(migrationFiles, root); e != nil {
		return
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)

	for _, name := range names {
		migrationName := driver + "/" + name

		var applied bool
		if applied, e = migrationApplied(ctx, database, migrationName); e != nil {
			return
		} else if applied {
			continue
		}

		var body []byte
		if body, e = migrationFiles.ReadFile(root + "/" + name); e != nil {
			return
		}

		var tx *sql.Tx
		if tx, e = database.BeginTx(ctx, nil); e != nil {
			return
		}
		if _, e = tx.ExecContext(ctx, string(body)); e != nil {
			_ = tx.Rollback()
			return fmt.Errorf("apply migration %s: %w", migrationName, e)
		}
		if _, e = tx.ExecContext(ctx, `INSERT INTO schema_migrations(name, applied_at) VALUES (?, ?)`, migrationName, time.Now().UTC().Format(time.RFC3339)); e != nil {
			_ = tx.Rollback()
			return
		}
		if e = tx.Commit(); e != nil {
			return
		}
	}

	return nil
}

func migrationApplied(ctx context.Context, database *sql.DB, name string) (_ bool, e error) {
	var existing string
	if e = database.QueryRowContext(ctx, `SELECT name FROM schema_migrations WHERE name = ?`, name).Scan(&existing); e == sql.ErrNoRows {
		return false, nil
	}

	return e == nil, e
}

func normalizeMigrationDriver(driver string) string {
	if driver == "sqlite3" {
		return "sqlite"
	}
	if driver == "mysql" {
		return "mysql"
	}
	return "sqlite"
}
