package db

import (
	"database/sql"
	_ "embed"
	"fmt"
)

type Conn struct {
	Conn *sql.DB
	*Queries
}

func NewConn(connStr string) (*Conn, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("opening engine: %w", err)
	}

	//go:embed sql/schema.sql
	var schema string

	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("executing schema: %w", err)
	}

	return &Conn{
		Conn:    db,
		Queries: New(db),
	}, nil
}
