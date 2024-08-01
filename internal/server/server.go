package server

import (
	"fmt"
	"net/http"

	"github.com/escoutdoor/social/internal/config"
	"github.com/escoutdoor/social/internal/postgres/store"
	"github.com/escoutdoor/social/internal/s3"
	"github.com/escoutdoor/social/internal/server/handlers"
)

type Opts struct {
	Config  config.Config
	S3Store *s3.MinIOClient
	Store   *store.Store
}

func New(opts Opts) *http.Server {
	user := handlers.NewUserHandler(opts.Store.User)
	auth := handlers.NewAuthHandler(opts.Store.Auth)
	post := handlers.NewPostHandler(opts.Store.Post)
	reply := handlers.NewReplyHandler(opts.Store.Reply, opts.Store.Post)
	file := handlers.NewFileHandler(opts.S3Store)
	api := &Server{
		user:  user,
		auth:  auth,
		post:  post,
		reply: reply,
		file:  file,
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", opts.Config.Port),
		Handler: api.NewRouter(opts.Store.Auth),
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
