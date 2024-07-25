package middlewares

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/escoutdoor/social/internal/db/store"
	"github.com/escoutdoor/social/internal/server/responses"
)

type AuthMiddleware struct {
	authStore store.AuthStorer
}

func NewAuthMiddleware(s store.AuthStorer) *AuthMiddleware {
	return &AuthMiddleware{
		authStore: s,
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

		id, err := m.authStore.ParseToken(jwtToken)
		if err != nil {
			slog.Error("AuthMiddleware: failed to parse token", "error", err.Error())
			responses.UnauthorizedResponse(w, store.ErrInvalidToken)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", id)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
