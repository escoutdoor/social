package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func New(connStr string) (*sql.DB, error) {
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("db connection error: %w", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := conn.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("db ping error: %w", err)
	}
	return conn, nil
}
