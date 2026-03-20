package db

import (
	"database/sql"
	"fmt"
	"os"
)

func RunMigrations(db *sql.DB) error {
	sqlBytes, err := os.ReadFile("migrations/001_init.sql")
	if err != nil {
		return fmt.Errorf("failed to read sql file: %w", err)
	}

	if _, err := db.Exec(string(sqlBytes)); err != nil {
		return fmt.Errorf("failed to run sql: %w", err)
	}
	return nil
}
