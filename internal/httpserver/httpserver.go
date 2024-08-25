package httpserver

import (
	"fmt"
	"net/http"

	"github.com/escoutdoor/social/internal/config"
	"github.com/escoutdoor/social/internal/httpserver/handlers"
	"github.com/escoutdoor/social/internal/service"
	"github.com/escoutdoor/social/pkg/validator"
)

type Opts struct {
	Config    *config.Config
	Services  *service.Services
	Validator *validator.Validator
}

func New(opts Opts) *http.Server {
	user := handlers.NewUserHandler(opts.Services.User, opts.Validator)
	auth := handlers.NewAuthHandler(opts.Services.Auth, opts.Validator)
	post := handlers.NewPostHandler(opts.Services.Post, opts.Validator)
	like := handlers.NewLikeHandler(opts.Services.Like)
	comment := handlers.NewCommentHandler(opts.Services.Comment, opts.Validator)
	file := handlers.NewFileHandler(opts.Services.File)

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
		Handler: api.NewRouter(opts.Services.Auth, opts.Services.User),
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
