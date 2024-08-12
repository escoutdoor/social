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
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	store     store.UserStorer
	validator *validator.Validator
}

func NewUserHandler(store store.UserStorer, v *validator.Validator) UserHandler {
	return UserHandler{
		store:     store,
		validator: v,
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
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}
	responses.JSON(w, http.StatusOK, envelope{"user": user})
}

func (h *UserHandler) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromCtx(r)
	if err != nil {
		responses.UnauthorizedResponse(w, err)
		return
	}

	var input types.UpdateUserReq
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responses.BadRequestResponse(w, ErrInvalidRequestBody)
		return
	}
	if err := h.validator.Validate(input); err != nil {
		responses.BadRequestResponse(w, err)
		return
	}

	if input.FirstName != nil {
		user.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		user.LastName = *input.LastName
	}
	if input.Email != nil {
		ok, err := h.store.GetByEmail(r.Context(), *input.Email)
		switch {
		case ok != nil:
			responses.BadRequestResponse(w, store.ErrEmailAlreadyExists)
			return
		case err != nil && !errors.Is(err, store.ErrUserNotFound):
			slog.Error("UserHandler.handleUpdateUser - UserStore.GetByEmail", "error", err)
			responses.InternalServerResponse(w, ErrInternalServer)
			return
		}
		user.Email = *input.Email
	}
	if input.Password != nil && !hasher.ComparePw(*input.Password, user.Password) {
		user.Password, err = hasher.HashPw(*input.Password)
		if err != nil {
			slog.Error("UserHandler.handleUpdateUser - hasher.HashPw", "error", err)
			responses.InternalServerResponse(w, ErrInternalServer)
			return
		}
	}
	if input.DOB != nil {
		dob, err := h.validator.ValidateDate(*input.DOB)
		if err != nil {
			responses.BadRequestResponse(w, err)
			return
		}
		user.DOB = &dob
	}
	if input.Bio != nil {
		user.Bio = input.Bio
	}
	if input.AvatarURL != nil {
		user.AvatarURL = input.AvatarURL
	}

	uu, err := h.store.Update(r.Context(), user.ID, *user)
	if err != nil {
		slog.Error("UserHandler.handleUpdateUser - UserStore.Update", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}
	responses.JSON(w, http.StatusOK, envelope{"user": uu})
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
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}
	responses.JSON(w, http.StatusOK, envelope{"message": "user successfully deleted"})
}
