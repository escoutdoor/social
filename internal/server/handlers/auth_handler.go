package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/escoutdoor/social/internal/postgres/store"
	"github.com/escoutdoor/social/internal/server/responses"
	"github.com/escoutdoor/social/internal/types"
	"github.com/escoutdoor/social/pkg/validator"
	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	store     store.AuthStorer
	validator *validator.Validator
}

func NewAuthHandler(s store.AuthStorer, v *validator.Validator) AuthHandler {
	return AuthHandler{
		store:     s,
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
		responses.BadRequestResponse(w, err)
		return
	}

	id, err := h.store.SignUp(r.Context(), input)
	if err != nil {
		if errors.Is(err, store.ErrUserAlreadyExists) {
			responses.BadRequestResponse(w, err)
			return
		}
		slog.Error("AuthHandler.handleSignUp - AuthStore.SignUp", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}

	token, err := h.store.GenerateToken(r.Context(), id)
	if err != nil {
		slog.Error("AuthHandler.handleSignUp - AuthStore.GenerateToken", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}
	responses.JSON(w, http.StatusOK, envelope{
		"user_id": id,
		"token":   token,
	})
}

func (h *AuthHandler) handleSignIn(w http.ResponseWriter, r *http.Request) {
	var input types.LoginReq
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responses.BadRequestResponse(w, ErrInvalidRequestBody)
		return
	}

	if err := h.validator.Validate(input); err != nil {
		responses.BadRequestResponse(w, err)
		return
	}

	user, err := h.store.SignIn(r.Context(), input)
	if err != nil {
		if errors.Is(err, store.ErrInvalidEmailOrPw) {
			responses.BadRequestResponse(w, err)
			return
		}
		slog.Error("AuthHandler.handleSignIn - AuthStore.SignIn", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}

	token, err := h.store.GenerateToken(r.Context(), user.ID)
	if err != nil {
		slog.Error("AuthHandler.handleSignIn - AuthStore.GenerateToken", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}
	responses.JSON(w, http.StatusOK, envelope{
		"user":  user,
		"token": token,
	})
}
