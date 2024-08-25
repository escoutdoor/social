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

type CommentHandler struct {
	svc       service.Comment
	validator *validator.Validator
}

func NewCommentHandler(svc service.Comment, v *validator.Validator) CommentHandler {
	return CommentHandler{
		svc:       svc,
		validator: v,
	}
}

func (h *CommentHandler) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/{id}", h.handleCreateComment)
	r.Get("/{id}", h.handleGetByID)
	r.Get("/all/{id}", h.handleGetAll)
	r.Delete("/{id}", h.handleDeleteComment)
	return r
}

func (h *CommentHandler) handleCreateComment(w http.ResponseWriter, r *http.Request) {
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

	var input types.CreateCommentReq
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responses.BadRequestResponse(w, ErrInvalidRequestBody)
		return
	}
	if err := h.validator.Validate(input); err != nil {
		responses.FailedValidationError(w, err)
		return
	}

	ctx := r.Context()
	id, err := h.svc.Create(ctx, user.ID, postID, input)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrPostNotFound):
			responses.NotFoundResponse(w, err)
			return
		case errors.Is(err, store.ErrCommentNotFound):
			responses.NotFoundResponse(w, err)
			return
		default:
			slog.Error("CommentHandler.handleCreateComment - CommentService.Create", "error", err)
			responses.InternalServerResponse(w, ErrInternalServer)
			return
		}
	}
	responses.JSON(w, http.StatusCreated, envelope{"id": id})
}

func (h *CommentHandler) handleGetByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDParam(r)
	if err != nil {
		responses.BadRequestResponse(w, err)
		return
	}

	ctx := r.Context()
	comment, err := h.svc.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, store.ErrCommentNotFound) {
			responses.NotFoundResponse(w, store.ErrCommentNotFound)
			return
		}
		slog.Error("CommentHandler.handleGetByID - CommentService.GetByID", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}
	responses.JSON(w, http.StatusOK, envelope{"comment": comment})
}

func (h *CommentHandler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	postID, err := getIDParam(r)
	if err != nil {
		responses.BadRequestResponse(w, err)
		return
	}

	ctx := r.Context()
	comments, err := h.svc.GetAll(ctx, postID)
	if err != nil {
		if errors.Is(err, store.ErrPostNotFound) {
			responses.NotFoundResponse(w, err)
			return
		}
		slog.Error("CommentHandler.handleGetAll - CommentService.GetAll", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}
	responses.JSON(w, http.StatusOK, envelope{"comments": comments})
}

func (h *CommentHandler) handleDeleteComment(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromCtx(r)
	if err != nil {
		responses.UnauthorizedResponse(w, err)
		return
	}
	commentID, err := getIDParam(r)
	if err != nil {
		responses.BadRequestResponse(w, err)
		return
	}

	ctx := r.Context()
	err = h.svc.Delete(ctx, commentID, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAccessDenied):
			responses.ForbiddenResponse(w, err)
			return
		case errors.Is(err, store.ErrCommentNotFound):
			responses.NotFoundResponse(w, err)
			return
		default:
			slog.Error("CommentHandler.handleDeleteComment - CommentService.Delete", "error", err)
			responses.InternalServerResponse(w, ErrInternalServer)
			return
		}
	}
	responses.JSON(w, http.StatusOK, envelope{"message": "comment successfully deleted"})
}
