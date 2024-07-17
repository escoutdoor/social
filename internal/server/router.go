package server

import "github.com/go-chi/chi"

func (s *Server) NewRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Route("/api", func(r chi.Router) {

	})
	return router
}
