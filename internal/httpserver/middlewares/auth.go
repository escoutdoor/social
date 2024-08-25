package middlewares

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/escoutdoor/social/internal/httpserver/responses"
	"github.com/escoutdoor/social/internal/postgres/store"
	"github.com/escoutdoor/social/internal/service"
)

const UserCtxKey string = "user"

type AuthMiddleware struct {
	authSvc service.Auth
	userSvc service.User
}

func NewAuthMiddleware(authSvc service.Auth, userSvc service.User) *AuthMiddleware {
	return &AuthMiddleware{
		authSvc: authSvc,
		userSvc: userSvc,
	}
}

func (m *AuthMiddleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtToken := r.Header.Get("Authorization")
		if len(jwtToken) == 0 {
			slog.Error("AuthMiddleware: missing bearer token", "error", ErrInvalidAuthorizationHeader.Error())
			responses.UnauthorizedResponse(w, ErrInvalidAuthorizationHeader)
			return
		}
		jwtToken = jwtToken[len("Bearer "):]
		id, err := m.authSvc.ParseToken(jwtToken)
		if err != nil {
			slog.Error("AuthMiddleware: failed to parse token", "error", err.Error())
			responses.UnauthorizedResponse(w, fmt.Errorf("failed to parse token: %w", err))
			return
		}

		user, err := m.userSvc.GetByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, store.ErrUserNotFound) {
				responses.UnauthorizedResponse(w, err)
				return
			}
			slog.Error("AuthMiddleware - UserStore.GetByID", "error", err.Error())
			responses.InternalServerResponse(w, fmt.Errorf("failed to authorize"))
			return
		}

		ctx := context.WithValue(r.Context(), UserCtxKey, user)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
