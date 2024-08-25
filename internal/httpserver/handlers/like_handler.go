package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/escoutdoor/social/internal/httpserver/responses"
	"github.com/escoutdoor/social/internal/postgres/store"
	"github.com/escoutdoor/social/internal/service"
	"github.com/go-chi/chi/v5"
)

type LikeHandler struct {
	svc service.Like
}

func NewLikeHandler(svc service.Like) LikeHandler {
	return LikeHandler{
		svc: svc,
	}
}

func (h *LikeHandler) Router() *chi.Mux {
	r := chi.NewRouter()
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

	ctx := r.Context()
	if err := h.svc.LikePost(ctx, postID, user.ID); err != nil {
		switch {
		case errors.Is(err, service.ErrAlreadyLiked):
			responses.BadRequestResponse(w, err)
			return
		case errors.Is(err, store.ErrPostNotFound):
			responses.NotFoundResponse(w, err)
			return
		default:
			slog.Error("LikeHandler.handleLikePost - LikeService.LikePost", "error", err)
			responses.InternalServerResponse(w, ErrInternalServer)
			return
		}
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

	ctx := r.Context()
	if err := h.svc.RemoveLikeFromPost(ctx, postID, user.ID); err != nil {
		if errors.Is(err, store.ErrRemoveLikeFailed) {
			responses.BadRequestResponse(w, err)
			return
		}
		slog.Error("LikeHandler.handleRemoveLikeFromPost - LikeService.RemoveLikeFromPost", "error", err)
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

	ctx := r.Context()
	if err := h.svc.LikeComment(ctx, commentID, user.ID); err != nil {
		switch {
		case errors.Is(err, service.ErrAlreadyLiked):
			responses.BadRequestResponse(w, err)
			return
		case errors.Is(err, store.ErrCommentNotFound):
			responses.NotFoundResponse(w, err)
			return
		default:
			slog.Error("LikeHandler.handleLikeComment - LikeService.LikeComment", "error", err)
			responses.InternalServerResponse(w, ErrInternalServer)
			return
		}
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

	ctx := r.Context()
	if err := h.svc.RemoveLikeFromComment(ctx, commentID, user.ID); err != nil {
		if errors.Is(err, store.ErrRemoveLikeFailed) {
			responses.BadRequestResponse(w, err)
			return
		}
		slog.Error("LikeHandler.handleRemoveLikeFromComment - LikeService.RemoveLikeFromComment", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}
	responses.JSON(w, http.StatusOK, envelope{"message": "you removed the like from the comment"})
}
