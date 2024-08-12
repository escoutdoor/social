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
	"github.com/go-chi/chi"
)

type CommentHandler struct {
	store     store.CommentStorer
	postStore store.PostStorer
	validator *validator.Validator
}

func NewCommentHandler(s store.CommentStorer, ps store.PostStorer, v *validator.Validator) CommentHandler {
	return CommentHandler{
		store:     s,
		postStore: ps,
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

	if _, err := h.postStore.GetByID(r.Context(), postID); err != nil {
		if errors.Is(err, store.ErrPostNotFound) {
			responses.NotFoundResponse(w, err)
			return
		}
		slog.Error("CommentHandler.handleCreateComment - PostStore.GetByID", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}

	var input types.CreateCommentReq
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responses.BadRequestResponse(w, ErrInvalidRequestBody)
		return
	}
	if err := h.validator.Validate(input); err != nil {
		responses.BadRequestResponse(w, err)
		return
	}

	if input.ParentCommentID != nil {
		if _, err := h.store.GetByID(r.Context(), *input.ParentCommentID); err != nil {
			if errors.Is(err, store.ErrCommentNotFound) {
				responses.NotFoundResponse(w, err)
				return
			}
			slog.Error("CommentHandler.handleCreateComment - CommentStore.GetByID", "error", err)
			responses.InternalServerResponse(w, ErrInternalServer)
			return
		}
	}
	id, err := h.store.Create(r.Context(), user.ID, postID, input)
	if err != nil {
		slog.Error("CommentHandler.handleCreateComment - CommentStore.Create", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}
	responses.JSON(w, http.StatusCreated, envelope{"id": id})
}

func (h *CommentHandler) handleGetByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDParam(r)
	if err != nil {
		responses.BadRequestResponse(w, err)
		return
	}

	comment, err := h.store.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrCommentNotFound) {
			responses.NotFoundResponse(w, store.ErrCommentNotFound)
			return
		}
		slog.Error("CommentHandler.handleGetByID - CommentStore.GetByID", "error", err)
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

	if _, err := h.postStore.GetByID(r.Context(), postID); err != nil {
		if errors.Is(err, store.ErrPostNotFound) {
			responses.NotFoundResponse(w, err)
			return
		}
		slog.Error("CommentHandler.handleGetAll - PostStore.GetByID", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}

	comments, err := h.store.GetAll(r.Context(), postID)
	if err != nil {
		slog.Error("CommentHandler.handleGetAll - CommentStore.GetAll", "error", err)
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

	comment, err := h.store.GetByID(r.Context(), commentID)
	if err != nil {
		if errors.Is(err, store.ErrCommentNotFound) {
			responses.NotFoundResponse(w, err)
			return
		}
		slog.Error("CommentHandler.handleDeleteComment - CommentStore.GetByID", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}
	if comment.UserID != user.ID {
		responses.ForbiddenResponse(w, ErrAccessDenied)
		return
	}

	err = h.store.Delete(r.Context(), commentID)
	if err != nil {
		slog.Error("CommentHandler.handleDeleteComment - CommentStore.Delete", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}
	responses.JSON(w, http.StatusOK, envelope{"message": "comment successfully deleted"})
}
