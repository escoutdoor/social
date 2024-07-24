package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/escoutdoor/social/internal/types"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func getIDParam(r *http.Request) (uuid.UUID, error) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return uuid.Nil, errors.New("invalid id parameter")
	}
	return id, nil
}

func getUserFromCtx(r *http.Request) (types.User, error) {
	user, ok := r.Context().Value("user").(types.User)
	if !ok {
		return user, fmt.Errorf("no user")
	}
	return user, nil
}

type envelope map[string]interface{}
