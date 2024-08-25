package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/escoutdoor/social/internal/httpserver/responses"
	"github.com/escoutdoor/social/internal/postgres/store"
	"github.com/escoutdoor/social/internal/service"
	"github.com/escoutdoor/social/internal/types"
	"github.com/escoutdoor/social/pkg/validator"
	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	svc       service.Auth
	validator *validator.Validator
}

func NewAuthHandler(svc service.Auth, v *validator.Validator) AuthHandler {
	return AuthHandler{
		svc:       svc,
		validator: v,
	}
}

func (h *AuthHandler) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/sign-up", h.handleSignUp)
	r.Post("/sign-in", h.handleSignIn)
	return r
}

func (h *AuthHandler) handleSignUp(w http.ResponseWriter, r *http.Request) {
	var input types.CreateUserReq
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responses.BadRequestResponse(w, ErrInvalidRequestBody)
		return
	}

	if err := h.validator.Validate(input); err != nil {
		responses.FailedValidationError(w, err)
		return
	}

	ctx := r.Context()
	id, err := h.svc.SignUp(ctx, input)
	if err != nil {
		if errors.Is(err, store.ErrUserAlreadyExists) {
			responses.BadRequestResponse(w, err)
			return
		}
		slog.Error("AuthHandler.handleSignUp - AuthService.SignUp", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}

	responses.JSON(w, http.StatusOK, envelope{
		"user_id": id,
	})
}

func (h *AuthHandler) handleSignIn(w http.ResponseWriter, r *http.Request) {
	var input types.LoginReq
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responses.BadRequestResponse(w, ErrInvalidRequestBody)
		return
	}

	if err := h.validator.Validate(input); err != nil {
		responses.FailedValidationError(w, err)
		return
	}

	ctx := r.Context()
	token, err := h.svc.SignIn(ctx, input)
	if err != nil {
		if errors.Is(err, service.ErrInvalidEmailOrPw) {
			responses.BadRequestResponse(w, err)
			return
		}
		slog.Error("AuthHandler.handleSignIn - AuthService.SignIn", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}
	responses.JSON(w, http.StatusOK, envelope{
		"token": token,
	})
}
