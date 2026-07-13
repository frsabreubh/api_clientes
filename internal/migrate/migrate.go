package migrate

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"sort"
	"strings"
)

//go:embed sql/*.sql
var migrationsFS embed.FS

// Run applies all pending up-migrations in version order.
func Run(db *sql.DB) error {
	// Ensure the tracking table exists.
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version    VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("creating schema_migrations: %w", err)
	}

	// Collect already-applied versions.
	rows, err := db.Query(`SELECT version FROM schema_migrations ORDER BY version`)
	if err != nil {
		return fmt.Errorf("reading schema_migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return err
		}
		applied[v] = true
	}
	if err := rows.Err(); err != nil {
		return err
	}

	// List and sort up-migration files.
	entries, err := fs.ReadDir(migrationsFS, "sql")
	if err != nil {
		return fmt.Errorf("reading migration directory: %w", err)
	}

	var upFiles []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".up.sql") {
			upFiles = append(upFiles, e.Name())
		}
	}
	sort.Strings(upFiles)

	// Apply pending migrations.
	for _, name := range upFiles {
		version := strings.TrimSuffix(name, ".up.sql")
		if applied[version] {
			continue
		}

		content, err := fs.ReadFile(migrationsFS, "sql/"+name)
		if err != nil {
			return fmt.Errorf("reading %s: %w", name, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("applying %s: %w", name, err)
		}

		if _, err := db.Exec(
			`INSERT INTO schema_migrations (version) VALUES ($1)`, version,
		); err != nil {
			return fmt.Errorf("recording %s: %w", name, err)
		}

		log.Printf("Migration applied: %s", version)
	}

	return nil
}
