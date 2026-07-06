package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	Host, User, Password, Name string
	Port                       int
	SSLMode                    string // "disable", "require", etc.
}

// The current SQL dialect that is being used
const dialect string = "postgres"

func Connect(cfg Config) (*sql.DB, error) {
	sslmode := cfg.SSLMode
	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, sslmode,
	)

	// Using database/sql with pgx via stdlib adapter.
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Pool settings (tune as needed).
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
