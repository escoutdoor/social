package app

import (
	"github.com/escoutdoor/social/internal/config"
	"github.com/escoutdoor/social/internal/server"
	"github.com/escoutdoor/social/pkg/logger"
)

func New() {
	logger.SetupLogger()
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	s := server.New(cfg)
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}
