package db

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

//go:embed schema.pgsql
var schema string

type DB struct {
	*pgxpool.Pool
	*DBQuerier
}

func New(dsn string) (*DB, error) {
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("create and connect pool: %w", err)
	}

	if _, err := pool.Exec(context.Background(), schema); err != nil {
		return nil, fmt.Errorf("executing schema: %w", err)
	}

	return &DB{
		Pool:      pool,
		DBQuerier: NewQuerier(pool),
	}, nil
}
