package server

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/escoutdoor/social/internal/config"
	"github.com/escoutdoor/social/internal/db/store"
	"github.com/escoutdoor/social/internal/s3"
	"github.com/escoutdoor/social/internal/server/handlers"
)

type Opts struct {
	Config    config.Config
	S3Storage *s3.MinIOClient
	DB        *sql.DB
}

func New(opts Opts) *http.Server {
	userStore := store.NewUserStore(opts.DB)
	user := handlers.NewUserHandler(userStore)

	authStore := store.NewAuthStore(opts.DB, opts.Config.JWTKey)
	auth := handlers.NewAuthHandler(authStore)

	postStore := store.NewPostStore(opts.DB)
	post := handlers.NewPostHandler(postStore)

	replyStore := store.NewReplyStore(opts.DB)
	reply := handlers.NewReplyHandler(replyStore, postStore)

	file := handlers.NewFileHandler(opts.S3Storage)
	api := &Server{
		user:  user,
		auth:  auth,
		post:  post,
		reply: reply,
		file:  file,
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", opts.Config.Port),
		Handler: api.NewRouter(authStore),
	}
	return server
}

type Server struct {
	user  handlers.UserHandler
	auth  handlers.AuthHandler
	post  handlers.PostHandler
	reply handlers.ReplyHandler
	file  handlers.FileHandler
}
