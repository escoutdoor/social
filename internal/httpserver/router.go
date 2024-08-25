package httpserver

import (
	"fmt"
	"net/http"

	"github.com/escoutdoor/social/internal/httpserver/middlewares"
	"github.com/escoutdoor/social/internal/httpserver/responses"
	"github.com/escoutdoor/social/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) NewRouter(authSvc service.Auth, userSvc service.User) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.CleanPath)
	router.MethodNotAllowed(methodNotAllowed)

	authMiddleware := middlewares.NewAuthMiddleware(authSvc, userSvc)
	router.Route("/v1", func(r chi.Router) {
		r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
			responses.JSON(w, http.StatusOK, map[string]string{
				"status": "ok",
			})
		})
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

func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	responses.ErrorResponse(w, http.StatusMethodNotAllowed, fmt.Sprintf("method (%s) is not supported for this resource", r.Method))
}
