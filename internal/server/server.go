package server

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/escoutdoor/social/internal/config"
	"github.com/escoutdoor/social/internal/db/store"
	"github.com/escoutdoor/social/internal/server/handlers"
)

type Opts struct {
	Config config.Config
	DB     *sql.DB
}

func New(opts Opts) *http.Server {
	userStore := store.NewUserStore(opts.DB)
	user := handlers.NewUserHandler(userStore)

	authStore := store.NewAuthStore(opts.DB)
	auth := handlers.NewAuthHandler(authStore)
	api := &Server{
		user: user,
		auth: auth,
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", opts.Config.Port),
		Handler: api.NewRouter(),
	}
	return server
}

type Server struct {
	user handlers.UserHandler
	auth handlers.AuthHandler
}
