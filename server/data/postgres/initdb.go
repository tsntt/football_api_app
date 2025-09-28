package data

import (
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDB(host, user, password, dbname, sslmode string, port int) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		slog.Error("Failed to connect to database", slog.String("err", err.Error()))
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
