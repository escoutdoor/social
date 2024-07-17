package server

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/escoutdoor/social/internal/config"
	"github.com/escoutdoor/social/internal/db"
	"github.com/escoutdoor/social/internal/db/store"
	"github.com/escoutdoor/social/internal/server/handlers"
)

type Server struct {
	user handlers.UserHandler
}

func New(c config.Config) *http.Server {
	db, err := db.New(c.PostgresURL)
	if err != nil {
		slog.Error("new server db conn", "error", err)
		os.Exit(3)
	}

	userStore := store.NewUserStore(db)
	user := handlers.NewUserHandler(userStore)

	api := &Server{
		user: user,
	}
	server := &http.Server{
		Addr:    ":8080",
		Handler: api.NewRouter(),
	}

	return server
}
