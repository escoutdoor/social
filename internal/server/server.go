package server

import (
	"fmt"
	"net/http"

	"github.com/escoutdoor/social/internal/cache"
	"github.com/escoutdoor/social/internal/config"
	"github.com/escoutdoor/social/internal/postgres/store"
	"github.com/escoutdoor/social/internal/s3"
	"github.com/escoutdoor/social/internal/server/handlers"
	"github.com/escoutdoor/social/pkg/validator"
)

type Opts struct {
	Config    config.Config
	Store     *store.Store
	S3Store   *s3.MinIOClient
	Cache     *cache.Cache
	Validator *validator.Validator
}

func New(opts Opts) *http.Server {
	user := handlers.NewUserHandler(opts.Store.User, opts.Validator)
	auth := handlers.NewAuthHandler(opts.Store.Auth, opts.Validator)
	post := handlers.NewPostHandler(opts.Store.Post, opts.Cache, opts.Validator)
	like := handlers.NewLikeHandler(opts.Store.Like, opts.Store.Post, opts.Store.Comment)
	comment := handlers.NewCommentHandler(opts.Store.Comment, opts.Store.Post, opts.Validator)
	file := handlers.NewFileHandler(opts.S3Store)

	api := &Server{
		user:    user,
		auth:    auth,
		post:    post,
		like:    like,
		comment: comment,
		file:    file,
	}
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", opts.Config.Port),
		Handler: api.NewRouter(opts.Store.Auth, opts.Store.User),
	}
	return server
}

type Server struct {
	user    handlers.UserHandler
	auth    handlers.AuthHandler
	post    handlers.PostHandler
	like    handlers.LikeHandler
	comment handlers.CommentHandler
	file    handlers.FileHandler
}
