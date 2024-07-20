package handlers

import (
	"errors"
	"net/http"

	"github.com/escoutdoor/social/internal/db/store"
	"github.com/escoutdoor/social/internal/server/responses"
)

type UserHandler struct {
	store store.UserStorer
}

func NewUserHandler(store store.UserStorer) UserHandler {
	return UserHandler{
		store: store,
	}
}

func (uh *UserHandler) handleGetByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDParam(r)
	if err != nil {
		responses.JSON(w, http.StatusBadRequest, err)
		return
	}

	user, err := uh.store.GetByID(id)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			responses.JSON(w, http.StatusNotFound, err)
			return
		}
		responses.JSON(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, user)
}

func (uh *UserHandler) handleUpdateUser(w http.ResponseWriter, r *http.Request) {

}
