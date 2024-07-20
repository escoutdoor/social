package server

import "github.com/go-chi/chi"

func (s *Server) NewRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Route("/v1", func(r chi.Router) {

	})
	return router
}
