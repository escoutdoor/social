package handlers

import (
	"errors"
	"fmt"
	"net/http"

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

func getUserIDFromCtx(r *http.Request) (uuid.UUID, error) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("no user")
	}
	pu, err := uuid.Parse(userID.String())
	if err != nil {
		return uuid.Nil, fmt.Errorf("no user")
	}
	return pu, nil
}

type envelope map[string]interface{}
