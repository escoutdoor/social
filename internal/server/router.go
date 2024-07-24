package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func (s *Server) NewRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello")
	})

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/auth", s.auth.Router())
		r.Mount("/users", s.user.Router())
		r.Mount("/posts", s.post.Router())
	})
	return router
}
