package suite

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/escoutdoor/social/internal/repository/postgres"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	pgmodule "github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	dbName = "test"
	dbUser = "test-user"
	dbPw   = "test-pw"
)

type Suite struct {
	Container testcontainers.Container
	DB        *sql.DB
}

func New() (*Suite, error) {
	ctx := context.Background()
	container, err := pgmodule.Run(
		ctx,
		"docker.io/pgmodule:16-alpine",
		pgmodule.WithDatabase(dbName),
		pgmodule.WithUsername(dbUser),
		pgmodule.WithPassword(dbPw),
		pgmodule.BasicWaitStrategies(),
		pgmodule.WithSQLDriver("postgres"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to run postgres container: %w", err)
	}

	p, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, fmt.Errorf("failed to get container external port: %w", err)
	}

	slog.Info("pgmodule container ready and running", "port", p.Port())

	dbAddr := fmt.Sprintf("localhost:%s", p.Port())
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUser, dbPw, dbAddr, dbName)
	db, err := postgres.New(dsn)
	if err != nil {
		return nil, err
	}

	if err := dbMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}

	return &Suite{
		Container: container,
		DB:        db,
	}, nil
}

func dbMigrate(db *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "../../migrations"); err != nil {
		return err
	}
	return nil
}
