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
	store        store.LikeStorer
	postStore    store.PostStorer
	commentStore store.CommentStorer
}

func NewLikeHandler(s store.LikeStorer, postStore store.PostStorer, commentStore store.CommentStorer) LikeHandler {
	return LikeHandler{
		store:        s,
		postStore:    postStore,
		commentStore: commentStore,
	}
}

func (h *LikeHandler) Router() *chi.Mux {
	r := chi.NewRouter()
	// likes/posts/id
	r.Route("/posts", func(r chi.Router) {
		r.Post("/{id}", h.handleLikePost)
		r.Delete("/{id}", h.handleRemoveLikeFromPost)
	})
	r.Route("/comments", func(r chi.Router) {
		r.Post("/{id}", h.handleLikeComment)
		r.Delete("/{id}", h.handleRemoveLikeFromComment)
	})
	return r
}

func (h *LikeHandler) handleLikePost(w http.ResponseWriter, r *http.Request) {
	postID, err := getIDParam(r)
	if err != nil {
		responses.BadRequestResponse(w, err)
		return
	}
	user, err := getUserFromCtx(r)
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
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}

	isLiked, err := h.store.IsPostLiked(r.Context(), postID)
	if err != nil {
		slog.Error("LikeHandler.handleLikePost - LikeStore.IsPostLiked", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}
	if isLiked {
		responses.BadRequestResponse(w, store.ErrAlreadyLiked)
		return
	}

	if err := h.store.LikePost(r.Context(), postID, user.ID); err != nil {
		slog.Error("LikeHandler.handleLikePost - LikeStore.LikePost", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
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
	user, err := getUserFromCtx(r)
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
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}

	if err := h.store.RemoveLikeFromPost(r.Context(), postID, user.ID); err != nil {
		if errors.Is(err, store.ErrRemoveLikeFailed) {
			responses.BadRequestResponse(w, err)
			return
		}
		slog.Error("LikeHandler.handleRemoveLikeFromPost - LikeStore.RemoveLikeFromPost", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}
	responses.JSON(w, http.StatusOK, envelope{"message": "you removed the like from the post"})
}

func (h *LikeHandler) handleLikeComment(w http.ResponseWriter, r *http.Request) {
	commentID, err := getIDParam(r)
	if err != nil {
		responses.BadRequestResponse(w, err)
		return
	}
	user, err := getUserFromCtx(r)
	if err != nil {
		responses.UnauthorizedResponse(w, err)
		return
	}

	if _, err = h.commentStore.GetByID(r.Context(), commentID); err != nil {
		if errors.Is(err, store.ErrCommentNotFound) {
			responses.NotFoundResponse(w, err)
			return
		}
		slog.Error("LikeHandler.handleLikeComment - CommentStore.GetByID", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}

	isLiked, err := h.store.IsCommentLiked(r.Context(), commentID)
	if err != nil {
		slog.Error("LikeHandler.handleLikeComment - LikeStore.IsCommentLiked", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}
	if isLiked {
		responses.BadRequestResponse(w, store.ErrAlreadyLiked)
		return
	}

	if err := h.store.LikeComment(r.Context(), commentID, user.ID); err != nil {
		slog.Error("LikeHandler.handleLikeComment - LikeStore.LikeComment", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}
	responses.JSON(w, http.StatusOK, envelope{"message": "you just liked the comment"})
}

func (h *LikeHandler) handleRemoveLikeFromComment(w http.ResponseWriter, r *http.Request) {
	commentID, err := getIDParam(r)
	if err != nil {
		responses.BadRequestResponse(w, err)
		return
	}
	user, err := getUserFromCtx(r)
	if err != nil {
		responses.UnauthorizedResponse(w, err)
		return
	}

	if _, err := h.commentStore.GetByID(r.Context(), commentID); err != nil {
		if errors.Is(err, store.ErrCommentNotFound) {
			responses.NotFoundResponse(w, err)
			return
		}
		slog.Error("LikeHandler.handleRemoveLikeFromComment - CommentStore.GetByID", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}

	if err := h.store.RemoveLikeFromComment(r.Context(), commentID, user.ID); err != nil {
		if errors.Is(err, store.ErrRemoveLikeFailed) {
			responses.BadRequestResponse(w, err)
			return
		}
		slog.Error("LikeHandler.handleRemoveLikeFromComment - LikeStore.RemoveLikeFromComment", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}
	responses.JSON(w, http.StatusOK, envelope{"message": "you removed the like from the comment"})
}
