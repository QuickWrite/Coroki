package db

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed migration/*.sql
var embedMigrations embed.FS

func Migrate(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(dialect); err != nil {
		return err
	}

	if err := goose.Up(db, "migration"); err != nil {
		return err
	}

	return nil
}
