package server

import (
	"fmt"
	"net/http"

	"github.com/escoutdoor/social/internal/postgres/store"
	"github.com/escoutdoor/social/internal/server/middlewares"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (s *Server) NewRouter(authStore store.AuthStorer, userStore store.UserStorer) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.StripSlashes)

	router.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "server is ok")
	})

	authMiddleware := middlewares.NewAuthMiddleware(authStore, userStore)
	router.Route("/v1", func(r chi.Router) {
		r.Mount("/auth", s.auth.Router())
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Auth)
			r.Mount("/users", s.user.Router())
			r.Mount("/posts", s.post.Router())
			r.Mount("/likes", s.like.Router())
			r.Mount("/comments", s.comment.Router())
			r.Mount("/files", s.file.Router())
		})
	})
	return router
}
