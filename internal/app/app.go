package app

import (
	"log/slog"
	"os"

	"github.com/escoutdoor/social/internal/config"
	"github.com/escoutdoor/social/internal/db"
	"github.com/escoutdoor/social/internal/server"
	"github.com/escoutdoor/social/pkg/logger"
)

func Run() {
	logger.SetupLogger()
	cfg, err := config.New()
	if err != nil {
		slog.Error("failed to initialize config", "error", err)
		os.Exit(1)
	}

	db, err := db.New(cfg.PostgresURL)
	if err != nil {
		slog.Error("new server db connection", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	slog.Info("server is running", slog.Int("port", cfg.Port))
	s := server.New(server.Opts{Config: cfg, DB: db})
	if err := s.ListenAndServe(); err != nil {
		slog.Error("server encountered an error", "error", err)
		os.Exit(1)
	}
	slog.Info("shutting down..")
}
