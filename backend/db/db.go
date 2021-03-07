package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"

	"go.senan.xyz/socr/backend/sql"
)

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

	if _, err := db.Exec(context.Background(), sql.Schema); err != nil {
		return nil, fmt.Errorf("executing schema: %w", err)
	}

	return &DB{
		Conn:      db,
		DBQuerier: NewQuerier(db),
	}, nil
}
