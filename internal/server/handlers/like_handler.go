package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/escoutdoor/social/internal/postgres/store"
	"github.com/escoutdoor/social/internal/server/responses"
	"github.com/go-chi/chi"
)

type LikeHandler struct {
	store     store.LikeStorer
	postStore store.PostStorer
}

func NewLikeHandler(s store.LikeStorer, postStore store.PostStorer) LikeHandler {
	return LikeHandler{
		store:     s,
		postStore: postStore,
	}
}

func (h *LikeHandler) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/{id}", h.handleLikePost)
	r.Delete("/{id}", h.handleRemoveLikeFromPost)
	return r
}

func (h *LikeHandler) handleLikePost(w http.ResponseWriter, r *http.Request) {
	postID, err := getIDParam(r)
	if err != nil {
		responses.BadRequestResponse(w, err)
		return
	}
	userID, err := getUserIDFromCtx(r)
	if err != nil {
		responses.UnauthorizedResponse(w, err)
		return
	}

	if _, err = h.postStore.GetByID(r.Context(), postID); err != nil {
		if errors.Is(err, store.ErrPostNotFound) {
			responses.NotFoundResponse(w, store.ErrPostNotFound)
			return
		}
		slog.Error("LikeHandler.handleLikePost - PostStore.GetByID", "error", err)
		responses.InternalServerResponse(w, ErrInternalServerError)
		return
	}

	isLiked, err := h.store.IsLiked(r.Context(), postID)
	if err != nil {
		slog.Error("LikeHandler.handleLikePost - LikeStore.IsLiked", "error", err)
		responses.InternalServerResponse(w, ErrInternalServerError)
		return
	}
	if isLiked {
		responses.BadRequestResponse(w, store.ErrPostAlreadyLiked)
		return
	}

	if err := h.store.Like(r.Context(), postID, userID); err != nil {
		slog.Error("LikeHandler.handleLikePost - LikeStore.Like", "error", err)
		responses.InternalServerResponse(w, ErrInternalServerError)
		return
	}
	responses.JSON(w, http.StatusOK, envelope{"message": "you just liked the post"})
}

func (h *LikeHandler) handleRemoveLikeFromPost(w http.ResponseWriter, r *http.Request) {
	postID, err := getIDParam(r)
	if err != nil {
		responses.BadRequestResponse(w, err)
		return
	}
	userID, err := getUserIDFromCtx(r)
	if err != nil {
		responses.UnauthorizedResponse(w, err)
		return
	}

	if _, err := h.postStore.GetByID(r.Context(), postID); err != nil {
		if errors.Is(err, store.ErrPostNotFound) {
			responses.NotFoundResponse(w, err)
			return
		}
		slog.Error("LikeHandler.handleRemoveLikeFromPost - PostStore.GetByID", "error", err)
		responses.InternalServerResponse(w, ErrInternalServerError)
		return
	}

	if err := h.store.RemoveLike(r.Context(), postID, userID); err != nil {
		if errors.Is(err, store.ErrFailedToRemoveLike) {
			responses.BadRequestResponse(w, err)
			return
		}
		slog.Error("LikeHandler.handleRemoveLikeFromPost - LikeStore.RemoveLike", "error", err)
		responses.InternalServerResponse(w, ErrInternalServerError)
		return
	}
	responses.JSON(w, http.StatusOK, envelope{"message": "you removed the like from the post"})
}
