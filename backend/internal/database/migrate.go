package database

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func RunMigrations(db *sql.DB, migrationsDir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("run goose up: %w", err)
	}

	return nil
}
