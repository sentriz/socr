package db

import (
	ssql "database/sql"
	"fmt"

	"go.senan.xyz/socr/backend/sql"
)

type Conn struct {
	Conn *ssql.DB
	*Queries
}

func NewConn(dsn string) (*Conn, error) {
	db, err := ssql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening engine: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("pinging: %w", err)
	}

	if _, err := db.Exec(sql.Schema); err != nil {
		return nil, fmt.Errorf("executing schema: %w", err)
	}

	return &Conn{
		Conn:    db,
		Queries: New(db),
	}, nil
}
