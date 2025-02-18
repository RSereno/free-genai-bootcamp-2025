package testutils

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

// SetupTestDB creates a new in-memory SQLite database for testing
func SetupTestDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file:testdb?mode=memory&cache=shared")
	if err != nil {
		return nil, fmt.Errorf("failed to open test database: %w", err)
	}

	// Run migrations
	migrationsPath := filepath.Join("..", "db", "migrations")
	if err := runMigrations(db, migrationsPath); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

func runMigrations(db *sql.DB, path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			content, err := os.ReadFile(filepath.Join(path, file.Name()))
			if err != nil {
				return err
			}
			if _, err := db.Exec(string(content)); err != nil {
				return fmt.Errorf("migration %s failed: %w", file.Name(), err)
			}
		}
	}
	return nil
}
