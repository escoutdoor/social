package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/escoutdoor/social/internal/postgres/store"
	"github.com/escoutdoor/social/internal/server/responses"
	"github.com/escoutdoor/social/internal/types"
	"github.com/escoutdoor/social/pkg/hasher"
	"github.com/escoutdoor/social/pkg/validator"
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
		responses.InternalServerResponse(w, ErrInternalServerError)
		return
	}
	responses.JSON(w, http.StatusOK, user)
}

func (h *UserHandler) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	userIDCtx, err := getUserIDFromCtx(r)
	if err != nil {
		responses.UnauthorizedResponse(w, err)
		return
	}

	u, err := h.store.GetByID(r.Context(), userIDCtx)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			responses.NotFoundResponse(w, store.ErrUserNotFound)
			return
		}
		slog.Error("UserHandler.handleUpdateUser - UserStore.GetByID", "error", err)
		responses.InternalServerResponse(w, ErrInternalServerError)
		return
	}

	var input types.UpdateUserReq
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responses.BadRequestResponse(w, ErrInvalidRequestBody)
		return
	}
	if err := validator.Validate(input); err != nil {
		responses.BadRequestResponse(w, err)
		return
	}

	if input.FirstName != nil {
		u.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		u.LastName = *input.LastName
	}
	if input.Email != nil {
		ok, err := h.store.GetByEmail(r.Context(), *input.Email)
		switch {
		case ok != nil:
			responses.BadRequestResponse(w, store.ErrEmailAlreadyExists)
			return
		case err != nil && !errors.Is(err, store.ErrUserNotFound):
			slog.Error("UserHandler.handleUpdateUser - UserStore.GetByEmail", "error", err)
			responses.InternalServerResponse(w, ErrInternalServerError)
			return
		}
		u.Email = *input.Email
	}
	if input.Password != nil && !hasher.ComparePw(*input.Password, u.Password) {
		u.Password, err = hasher.HashPw(*input.Password)
		if err != nil {
			slog.Error("UserHandler.handleUpdateUser - hasher.HashPw", "error", err)
			responses.InternalServerResponse(w, ErrInternalServerError)
			return
		}
	}
	if input.DOB != nil {
		dob, err := validator.ValidateDate(*input.DOB)
		if err != nil {
			responses.BadRequestResponse(w, err)
			return
		}
		u.DOB = &dob
	}
	if input.Bio != nil {
		u.Bio = input.Bio
	}
	if input.AvatarURL != nil {
		u.AvatarURL = input.AvatarURL
	}

	user, err := h.store.Update(r.Context(), userIDCtx, *u)
	if err != nil {
		slog.Error("UserHandler.handleUpdateUser - UserStore.Update", "error", err)
		responses.InternalServerResponse(w, ErrInternalServerError)
		return
	}
	responses.JSON(w, http.StatusOK, user)
}

func (h *UserHandler) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	userIDCtx, err := getUserIDFromCtx(r)
	if err != nil {
		responses.UnauthorizedResponse(w, err)
		return
	}

	err = h.store.Delete(r.Context(), userIDCtx)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			responses.NotFoundResponse(w, err)
			return
		}
		slog.Error("UserHandler.Delete - UserStore.Delete", "error", err)
		responses.InternalServerResponse(w, ErrInternalServerError)
		return
	}
	responses.JSON(w, http.StatusOK, envelope{"message": "user successfully deleted"})
}
