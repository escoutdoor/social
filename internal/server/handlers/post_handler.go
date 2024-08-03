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

type PostHandler struct {
	store store.PostStorer
}

func NewPostHandler(store store.PostStorer) PostHandler {
	return PostHandler{
		store: store,
	}
}

func (h *PostHandler) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/", h.handleCreatePost)
	r.Put("/{id}", h.handleUpdatePost)
	r.Get("/{id}", h.handleGetByID)
	r.Delete("/{id}", h.handleDeletePost)
	return r
}

func (h *PostHandler) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	userIDCtx, err := getUserIDFromCtx(r)
	if err != nil {
		responses.UnauthorizedResponse(w, err)
		return
	}

	var input types.CreatePostReq
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responses.BadRequestResponse(w, ErrInvalidRequestBody)
		return
	}

	if err := validator.Validate(input); err != nil {
		responses.BadRequestResponse(w, err)
		return
	}

	id, err := h.store.Create(r.Context(), userIDCtx, input)
	if err != nil {
		slog.Error("PostHandler.handleCreatePost - PostStore.Create", "error", err)
		responses.InternalServerResponse(w, ErrInternalServerError)
		return
	}
	responses.JSON(w, http.StatusCreated, envelope{"id": id})
}

func (h *PostHandler) handleUpdatePost(w http.ResponseWriter, r *http.Request) {
	userIDCtx, err := getUserIDFromCtx(r)
	if err != nil {
		responses.UnauthorizedResponse(w, err)
		return
	}
	postID, err := getIDParam(r)
	if err != nil {
		responses.BadRequestResponse(w, err)
		return
	}

	p, err := h.store.GetByID(r.Context(), postID)
	if err != nil {
		if errors.Is(err, store.ErrPostNotFound) {
			responses.NotFoundResponse(w, err)
			return
		}
		slog.Error("PostHandler.handleUpdatePost - PostStore.GetByID", "error", err)
		responses.InternalServerResponse(w, ErrInternalServerError)
		return
	}
	if p.UserID != userIDCtx {
		responses.ForbiddenResponse(w, ErrForbidden)
		return
	}

	var input types.UpdatePostReq
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responses.BadRequestResponse(w, ErrInvalidRequestBody)
		return
	}
	if err := validator.Validate(input); err != nil {
		responses.BadRequestResponse(w, err)
		return
	}

	if input.Content != nil {
		p.Content = *input.Content
	}
	if input.PhotoURL != nil {
		p.PhotoURL = input.PhotoURL
	}
	post, err := h.store.Update(r.Context(), postID, *p)
	if err != nil {
		slog.Error("PostHandler.handleUpdatePost - PostStore.Update", "error", err)
		responses.InternalServerResponse(w, ErrInternalServerError)
		return
	}
	responses.JSON(w, http.StatusOK, post)
}

func (h *PostHandler) handleGetByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDParam(r)
	if err != nil {
		responses.BadRequestResponse(w, err)
		return
	}

	post, err := h.store.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrPostNotFound) {
			responses.NotFoundResponse(w, store.ErrPostNotFound)
			return
		}
		slog.Error("PostHandler.handleGetByID - PostStore.GetByID", "error", err)
		responses.InternalServerResponse(w, ErrInternalServerError)
		return
	}
	responses.JSON(w, http.StatusOK, post)
}

func (h *PostHandler) handleDeletePost(w http.ResponseWriter, r *http.Request) {
	userIDCtx, err := getUserIDFromCtx(r)
	if err != nil {
		responses.UnauthorizedResponse(w, err)
		return
	}
	postID, err := getIDParam(r)
	if err != nil {
		responses.BadRequestResponse(w, err)
		return
	}

	post, err := h.store.GetByID(r.Context(), postID)
	if err != nil {
		if errors.Is(err, store.ErrPostNotFound) {
			responses.NotFoundResponse(w, store.ErrPostNotFound)
			return
		}
		slog.Error("PostHandler.handleDeletePost - PostStore.GetByID", "error", err)
		responses.InternalServerResponse(w, ErrInternalServerError)
		return
	}
	if post.UserID != userIDCtx {
		responses.ForbiddenResponse(w, ErrForbidden)
		return
	}

	err = h.store.Delete(r.Context(), postID)
	if err != nil {
		slog.Error("PostHandler.handleDeletePost - PostStore.Delete", "error", err)
		responses.InternalServerResponse(w, ErrInternalServerError)
		return
	}
	responses.JSON(w, http.StatusOK, envelope{"message": "post successfully deleted"})
}
