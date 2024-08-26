package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/escoutdoor/social/internal/httpserver/responses"
	"github.com/escoutdoor/social/internal/repository/repoerrs"
	"github.com/escoutdoor/social/internal/service"
	"github.com/escoutdoor/social/internal/types"
	"github.com/escoutdoor/social/pkg/validator"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	svc       service.User
	validator *validator.Validator
}

func NewUserHandler(svc service.User, v *validator.Validator) UserHandler {
	return UserHandler{
		svc:       svc,
		validator: v,
	}
}

func (h *UserHandler) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Patch("/", h.handleUpdateUser)
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

	ctx := r.Context()
	user, err := h.svc.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repoerrs.ErrUserNotFound) {
			responses.NotFoundResponse(w, err)
			return
		}
		slog.Error("UserHandler.handleGetByID - UserService.GetByID", "error", err)
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
		responses.FailedValidationError(w, err)
		return
	}

	ctx := r.Context()
	uu, err := h.svc.Update(ctx, *user, input)
	if err != nil {
		if errors.Is(err, validator.ErrInvalidDateFormat) {
			responses.BadRequestResponse(w, err)
			return
		}

		slog.Error("UserHandler.handleUpdateUser - UserService.Update", "error", err)
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

	ctx := r.Context()
	err = h.svc.Delete(ctx, user.ID)
	if err != nil {
		if errors.Is(err, repoerrs.ErrUserNotFound) {
			responses.NotFoundResponse(w, err)
			return
		}
		slog.Error("UserHandler.Delete - UserService.Delete", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}
	responses.JSON(w, http.StatusOK, envelope{"message": "user successfully deleted"})
}
