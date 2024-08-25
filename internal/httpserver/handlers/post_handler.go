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

type PostHandler struct {
	svc       service.Post
	validator *validator.Validator
}

func NewPostHandler(svc service.Post, v *validator.Validator) PostHandler {
	return PostHandler{
		svc:       svc,
		validator: v,
	}
}

func (h *PostHandler) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/", h.handleCreatePost)
	r.Get("/", h.handleGetAll)
	r.Get("/{id}", h.handleGetByID)
	r.Patch("/{id}", h.handleUpdatePost)
	r.Delete("/{id}", h.handleDeletePost)

	return r
}

func (h *PostHandler) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromCtx(r)
	if err != nil {
		responses.UnauthorizedResponse(w, err)
		return
	}

	var input types.CreatePostReq
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responses.BadRequestResponse(w, ErrInvalidRequestBody)
		return
	}
	if err := h.validator.Validate(input); err != nil {
		responses.FailedValidationError(w, err)
		return
	}

	ctx := r.Context()
	post, err := h.svc.Create(ctx, user.ID, input)
	if err != nil {
		slog.Error("PostHandler.handleCreatePost - PostService.Create", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}
	responses.JSON(w, http.StatusCreated, envelope{"post": post})
}

func (h *PostHandler) handleUpdatePost(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromCtx(r)
	if err != nil {
		responses.UnauthorizedResponse(w, err)
		return
	}
	postID, err := getIDParam(r)
	if err != nil {
		responses.BadRequestResponse(w, err)
		return
	}

	var input types.UpdatePostReq
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responses.BadRequestResponse(w, ErrInvalidRequestBody)
		return
	}
	if err := h.validator.Validate(input); err != nil {
		responses.FailedValidationError(w, err)
		return
	}

	ctx := r.Context()
	post, err := h.svc.Update(ctx, postID, user.ID, input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAccessDenied):
			responses.ForbiddenResponse(w, err)
			return
		case errors.Is(err, store.ErrPostNotFound):
			responses.NotFoundResponse(w, err)
			return
		default:
			slog.Error("PostHandler.handleUpdatePost - PostService.Update", "error", err)
			responses.InternalServerResponse(w, ErrInternalServer)
			return
		}
	}
	responses.JSON(w, http.StatusOK, envelope{"post": post})
}

func (h *PostHandler) handleGetByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDParam(r)
	if err != nil {
		responses.BadRequestResponse(w, err)
		return
	}

	ctx := r.Context()
	post, err := h.svc.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, store.ErrPostNotFound) {
			responses.NotFoundResponse(w, err)
			return
		}
		slog.Error("PostHandler.handleGetByID - PostService.GetByID", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}
	responses.JSON(w, http.StatusOK, envelope{"post": post})
}

func (h *PostHandler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	posts, err := h.svc.GetAll(ctx)
	if err != nil {
		slog.Error("PostHandler.handleGetAll - PostService.GetAll", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}
	responses.JSON(w, http.StatusOK, envelope{"posts": posts})
}

func (h *PostHandler) handleDeletePost(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromCtx(r)
	if err != nil {
		responses.UnauthorizedResponse(w, err)
		return
	}
	postID, err := getIDParam(r)
	if err != nil {
		responses.BadRequestResponse(w, err)
		return
	}

	ctx := r.Context()
	err = h.svc.Delete(ctx, postID, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAccessDenied):
			responses.ForbiddenResponse(w, err)
			return
		case errors.Is(err, store.ErrPostNotFound):
			responses.NotFoundResponse(w, err)
			return
		default:
			slog.Error("PostHandler.handleDeletePost - PostService.Delete", "error", err)
			responses.InternalServerResponse(w, ErrInternalServer)
			return
		}
	}
	responses.JSON(w, http.StatusOK, envelope{"message": "post successfully deleted"})
}
