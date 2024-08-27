package app

import (
	"fmt"
	"log/slog"

	"github.com/escoutdoor/social/internal/cache"
	"github.com/escoutdoor/social/internal/config"
	"github.com/escoutdoor/social/internal/httpserver"
	"github.com/escoutdoor/social/internal/repository"
	"github.com/escoutdoor/social/internal/repository/postgres"
	"github.com/escoutdoor/social/internal/s3"
	"github.com/escoutdoor/social/internal/service"
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
	repo := repository.New(db)

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

	services := service.NewServices(service.Opts{
		Repository: repo,
		Cache:      cache,
		S3:         s3,
		Validator:  validator,
		SignKey:    cfg.SignKey,
	})

	slog.Info("server is running", slog.Int("port", cfg.Port))
	s := httpserver.New(httpserver.Opts{
		Config:    cfg,
		Services:  services,
		Validator: validator,
	})
	if err := s.ListenAndServe(); err != nil {
		return fmt.Errorf("server encountered an error: %w", err)
	}
	slog.Info("shutting down..")
	return nil
}
