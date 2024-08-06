package middlewares

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/escoutdoor/social/internal/postgres/store"
	"github.com/escoutdoor/social/internal/server/handlers"
	"github.com/escoutdoor/social/internal/server/responses"
)

const userIDCtx string = "user_id"

type AuthMiddleware struct {
	authStore store.AuthStorer
	userStore store.UserStorer
}

func NewAuthMiddleware(authStore store.AuthStorer, userStore store.UserStorer) *AuthMiddleware {
	return &AuthMiddleware{
		authStore: authStore,
		userStore: userStore,
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

		if _, err := m.userStore.GetByID(r.Context(), id); err != nil {
			if errors.Is(err, store.ErrUserNotFound) {
				responses.UnauthorizedResponse(w, store.ErrInvalidToken)
				return
			}
			slog.Error("AuthMiddleware - UserStore.GetByID", "error", err.Error())
			responses.InternalServerResponse(w, handlers.ErrInternalServer)
			return
		}

		ctx := context.WithValue(r.Context(), userIDCtx, id)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
