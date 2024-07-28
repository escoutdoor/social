package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/escoutdoor/social/internal/db/store"
	"github.com/escoutdoor/social/internal/server/responses"
	"github.com/escoutdoor/social/internal/types"
	"github.com/escoutdoor/social/pkg/validation"
	"github.com/go-chi/chi"
)

type ReplyHandler struct {
	store     store.ReplyStorer
	postStore store.PostStorer
}

func NewReplyHandler(s store.ReplyStorer, ps store.PostStorer) ReplyHandler {
	return ReplyHandler{
		store:     s,
		postStore: ps,
	}
}

func (h *ReplyHandler) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/{id}", h.handleCreateReply)
	r.Get("/{id}", h.handleGetByID)
	r.Delete("/{id}", h.handleDeleteReply)
	return r
}

func (h *ReplyHandler) handleCreateReply(w http.ResponseWriter, r *http.Request) {
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

	if _, err := h.postStore.GetByID(r.Context(), postID); err != nil {
		if errors.Is(err, store.ErrPostNotFound) {
			responses.NotFoundResponse(w, err)
			return
		}
		slog.Error("ReplyHandler.handleCreateReply - PostStore.GetByID", "error", err)
		responses.InternalServerResponse(w, ErrIntervalServerError)
		return
	}

	var input types.CreateReplyReq
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responses.BadRequestResponse(w, ErrInvalidRequestBody)
		return
	}
	if err := validation.Validate(input); err != nil {
		responses.BadRequestResponse(w, err)
		return
	}

	id, err := h.store.Create(r.Context(), userIDCtx, postID, input)
	if err != nil {
		slog.Error("ReplyHandler.handleCreateReply - ReplyStore.Create", "error", err)
		responses.InternalServerResponse(w, ErrIntervalServerError)
		return
	}
	responses.JSON(w, http.StatusCreated, envelope{"id": id})
}

func (h *ReplyHandler) handleGetByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDParam(r)
	if err != nil {
		responses.BadRequestResponse(w, err)
		return
	}

	reply, err := h.store.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrReplyNotFound) {
			responses.NotFoundResponse(w, store.ErrReplyNotFound)
			return
		}
		slog.Error("ReplyHandler.handleGetByID - ReplyStore.GetByID", "error", err)
		responses.InternalServerResponse(w, ErrIntervalServerError)
		return
	}
	responses.JSON(w, http.StatusOK, reply)
}

func (h *ReplyHandler) handleDeleteReply(w http.ResponseWriter, r *http.Request) {
	userIDCtx, err := getUserIDFromCtx(r)
	if err != nil {
		responses.UnauthorizedResponse(w, err)
		return
	}
	replyID, err := getIDParam(r)
	if err != nil {
		responses.BadRequestResponse(w, err)
		return
	}

	reply, err := h.store.GetByID(r.Context(), replyID)
	if err != nil {
		if errors.Is(err, store.ErrReplyNotFound) {
			responses.NotFoundResponse(w, err)
			return
		}
		slog.Error("ReplyHandler.handleDeleteReply - ReplyStore.GetByID", "error", err)
		responses.InternalServerResponse(w, ErrIntervalServerError)
		return
	}
	if reply.UserID != userIDCtx {
		responses.ForbiddenResponse(w, ErrForbidden)
		return
	}

	err = h.store.Delete(r.Context(), replyID)
	if err != nil {
		slog.Error("ReplyHandler.handleDeleteReply - ReplyStore.Delete", "error", err)
		responses.InternalServerResponse(w, ErrIntervalServerError)
		return
	}
	responses.JSON(w, http.StatusOK, envelope{"message": "reply successfully deleted"})
}
