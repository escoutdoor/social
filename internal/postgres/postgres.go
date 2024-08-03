package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func New(connStr string) (*sql.DB, error) {
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("db connection error: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("db ping error: %w", err)
	}
	return conn, nil
}
