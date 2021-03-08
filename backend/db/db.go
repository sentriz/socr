package db

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v4"
)

//go:embed schema.pgsql
var schema string

type DB struct {
	*pgx.Conn
	*DBQuerier
}

func New(dsn string) (*DB, error) {
	db, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("opening engine: %w", err)
	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("pinging: %w", err)
	}

	if _, err := db.Exec(context.Background(), schema); err != nil {
		return nil, fmt.Errorf("executing schema: %w", err)
	}

	return &DB{
		Conn:      db,
		DBQuerier: NewQuerier(db),
	}, nil
}
