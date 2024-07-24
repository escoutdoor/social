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
	r.Get("/{id}", h.handleGetByID)
	r.Post("/", h.handleCreatePost)
	r.Delete("/{id}", h.handleDeletePost)

	return r
}

func (h *PostHandler) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromCtx(r)
	if err != nil {
		responses.JSON(w, http.StatusUnauthorized, err)
		return
	}

	var input types.CreatePostReq
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responses.BadRequestResponse(w, ErrInvalidRequestBody)
		return
	}

	if err := validation.Validate(input); err != nil {
		responses.BadRequestResponse(w, err)
		return
	}

	id, err := h.store.Create(r.Context(), user.ID, input)
	if err != nil {
		slog.Error("PostHandler.handleCreatePost - PostStore.Create", "error", err)
		responses.InternalServerResponse(w, ErrIntervalServerError)
		return
	}
	responses.JSON(w, http.StatusCreated, envelope{"id": id})
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
		responses.InternalServerResponse(w, ErrIntervalServerError)
		return
	}
	responses.JSON(w, http.StatusOK, post)
}

func (h *PostHandler) handleDeletePost(w http.ResponseWriter, r *http.Request) {
	id, err := getIDParam(r)
	if err != nil {
		responses.BadRequestResponse(w, err)
		return
	}

	err = h.store.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, store.ErrPostNotFound) {
			responses.NotFoundResponse(w, store.ErrPostNotFound)
			return
		}
		slog.Error("PostHandler.handleDeletePost - PostStore.Delete", "error", err)
		responses.InternalServerResponse(w, ErrIntervalServerError)
		return
	}
	responses.JSON(w, http.StatusOK, envelope{"message": "post successfully deleted"})
}
