package app

import (
	"log/slog"
	"os"

	"github.com/escoutdoor/social/internal/cache"
	"github.com/escoutdoor/social/internal/config"
	"github.com/escoutdoor/social/internal/postgres"
	"github.com/escoutdoor/social/internal/postgres/store"
	"github.com/escoutdoor/social/internal/s3"
	"github.com/escoutdoor/social/internal/server"
	"github.com/escoutdoor/social/pkg/logger"
	"github.com/escoutdoor/social/pkg/validator"
)

func Run() {
	logger.SetupLogger()
	cfg, err := config.New()
	if err != nil {
		slog.Error("failed to initialize config", "error", err)
		os.Exit(1)
	}

	db, err := postgres.New(cfg.PostgresURL)
	if err != nil {
		slog.Error("failed to connect to postgres", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	slog.Info("successfully connected to postgres")

	store := store.NewStore(db, cfg)

	s3, err := s3.New(cfg)
	if err != nil {
		slog.Error("failed to connect to s3", "error", err)
		os.Exit(1)
	}
	slog.Info("successfully connected to s3")

	cache, err := cache.New(cfg)
	if err != nil {
		slog.Error("failed to connect to redis", "error", err)
		os.Exit(1)
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
		slog.Error("server encountered an error", "error", err)
		os.Exit(1)
	}
	slog.Info("shutting down..")
}
