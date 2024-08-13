package app

import (
	"fmt"
	"log/slog"

	"github.com/escoutdoor/social/internal/cache"
	"github.com/escoutdoor/social/internal/config"
	"github.com/escoutdoor/social/internal/postgres"
	"github.com/escoutdoor/social/internal/postgres/store"
	"github.com/escoutdoor/social/internal/s3"
	"github.com/escoutdoor/social/internal/server"
	"github.com/escoutdoor/social/pkg/logger"
	"github.com/escoutdoor/social/pkg/validator"
)

func Run() error {
	logger.SetupLogger()
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}

	db, err := postgres.New(cfg.PostgresURL)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	defer db.Close()
	slog.Info("successfully connected to postgres")

	store := store.NewStore(db, cfg)

	s3, err := s3.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to s3: %w", err)
	}
	slog.Info("successfully connected to s3")

	cache, err := cache.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}
	slog.Info("successfully connected to redis")

	validator := validator.New()
	slog.Info("server is running", slog.Int("port", cfg.Port))
	s := server.New(server.Opts{
		Config:    cfg,
		Store:     store,
		S3Store:   s3,
		Cache:     cache,
		Validator: validator,
	})
	if err := s.ListenAndServe(); err != nil {
		return fmt.Errorf("server encountered an error: %w", err)
	}
	slog.Info("shutting down..")
	return nil
}
