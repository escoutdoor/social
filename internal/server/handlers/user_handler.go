package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/escoutdoor/social/internal/db/store"
	"github.com/escoutdoor/social/internal/server/responses"
	"github.com/escoutdoor/social/internal/types"
	"github.com/escoutdoor/social/pkg/validation"
	"github.com/go-chi/chi"
)

type UserHandler struct {
	store store.UserStorer
}

func NewUserHandler(store store.UserStorer) UserHandler {
	return UserHandler{
		store: store,
	}
}

func (h *UserHandler) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Put("/", h.handleUpdateUser)
	r.Delete("/", h.handleDeleteUser)
	r.Get("/{id}", h.handleGetByID)

	return r
}

func (h *UserHandler) handleGetByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDParam(r)
	if err != nil {
		responses.BadRequestResponse(w, err)
		return
	}

	user, err := h.store.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			responses.NotFoundResponse(w, err)
			return
		}
		slog.Error("UserHandler.handleGetByID - UserStore.GetByID", "error", err)
		responses.InternalServerResponse(w, ErrIntervalServerError)
		return
	}
	responses.JSON(w, http.StatusOK, user)
}

func (h *UserHandler) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	userFromCtx, err := getUserFromCtx(r)
	if err != nil {
		responses.UnauthorizedResponse(w, err)
		return
	}

	var input types.UpdateUserReq
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responses.BadRequestResponse(w, ErrInvalidRequestBody)
		return
	}

	if err := validation.Validate(input); err != nil {
		responses.BadRequestResponse(w, err)
		return
	}

	user, err := h.store.Update(r.Context(), userFromCtx.ID, input)
	if err != nil {
		slog.Error("UserHandler.handleUpdateUser - UserStore.Update", "error", err)
		responses.InternalServerResponse(w, ErrIntervalServerError)
		return
	}
	responses.JSON(w, http.StatusOK, user)
}

func (h *UserHandler) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromCtx(r)
	if err != nil {
		responses.UnauthorizedResponse(w, err)
		return
	}

	err = h.store.Delete(r.Context(), user.ID)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			responses.NotFoundResponse(w, err)
			return
		}
		slog.Error("UserHandler.Delete - UserStore.Delete", "error", err)
		responses.InternalServerResponse(w, ErrIntervalServerError)
		return
	}
	responses.JSON(w, http.StatusOK, envelope{"message": "user successfully deleted"})
}
