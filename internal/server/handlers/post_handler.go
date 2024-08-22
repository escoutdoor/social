package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/escoutdoor/social/internal/cache"
	"github.com/escoutdoor/social/internal/postgres/store"
	"github.com/escoutdoor/social/internal/server/responses"
	"github.com/escoutdoor/social/internal/types"
	"github.com/escoutdoor/social/pkg/validator"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type PostHandler struct {
	store     store.PostStorer
	cache     *cache.Cache
	validator *validator.Validator
}

func NewPostHandler(store store.PostStorer, cache *cache.Cache, v *validator.Validator) PostHandler {
	return PostHandler{
		store:     store,
		cache:     cache,
		validator: v,
	}
}

func (h *PostHandler) Router() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/", h.handleCreatePost)
	r.Get("/", h.handleGetAll)
	r.Get("/{id}", h.handleGetByID)
	r.Put("/{id}", h.handleUpdatePost)
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
		responses.BadRequestResponse(w, err)
		return
	}

	ctx := r.Context()
	post, err := h.store.Create(ctx, user.ID, input)
	if err != nil {
		slog.Error("PostHandler.handleCreatePost - PostStore.Create", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}

	key := generatePostKey(post.ID)
	if err := h.cache.Set(ctx, key, post, time.Minute*1).Err(); err != nil {
		slog.Error("PostHandler.handleCreatePost - Rdb.Set", "error", err)
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

	ctx := r.Context()
	key := generatePostKey(postID)
	p, err := h.cache.GetPost(ctx, key)
	if errors.Is(err, redis.Nil) {
		p, err = h.store.GetByID(ctx, postID)
		if err != nil {
			if errors.Is(err, store.ErrPostNotFound) {
				responses.NotFoundResponse(w, err)
				return
			}
			slog.Error("PostHandler.handleUpdatePost - PostStore.GetByID", "error", err)
			responses.InternalServerResponse(w, ErrInternalServer)
			return
		}
	}
	if err != nil {
		slog.Error("PostHandler.handleUpdatePost - Cache.GetPost", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}
	if p.UserID != user.ID {
		responses.ForbiddenResponse(w, ErrAccessDenied)
		return
	}

	var input types.UpdatePostReq
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		responses.BadRequestResponse(w, ErrInvalidRequestBody)
		return
	}
	if err := h.validator.Validate(input); err != nil {
		responses.BadRequestResponse(w, err)
		return
	}

	if input.Content != nil {
		p.Content = *input.Content
	}
	if input.PhotoURL != nil {
		p.PhotoURL = input.PhotoURL
	}
	post, err := h.store.Update(ctx, postID, *p)
	if err != nil {
		slog.Error("PostHandler.handleUpdatePost - PostStore.Update", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}

	if err := h.cache.Set(ctx, key, post, time.Minute*1).Err(); err != nil {
		slog.Error("PostHandler.handleUpdatePost - Cache.Set", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
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
	key := generatePostKey(id)
	post, err := h.cache.GetPost(ctx, key)
	if errors.Is(err, redis.Nil) {
		post, err = h.store.GetByID(ctx, id)
		if err != nil {
			if errors.Is(err, store.ErrPostNotFound) {
				responses.NotFoundResponse(w, store.ErrPostNotFound)
				return
			}
			slog.Error("PostHandler.handleGetByID - PostStore.GetByID", "error", err)
			responses.InternalServerResponse(w, ErrInternalServer)
			return
		}
		if err := h.cache.Set(ctx, key, post, time.Minute*1).Err(); err != nil {
			slog.Error("PostHandler.handleGetByID - Cache.Set", "error", err)
			responses.InternalServerResponse(w, ErrInternalServer)
			return
		}
	}
	if err != nil {
		slog.Error("PostHandler.handleGetByID - Cache.GetPost", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}
	responses.JSON(w, http.StatusOK, envelope{"post": post})
}

func (h *PostHandler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	var posts types.Posts
	key := "posts"

	ctx := r.Context()
	posts, err := h.cache.GetPosts(ctx, key)
	if errors.Is(err, redis.Nil) {
		posts, err = h.store.GetAll(ctx)
		if err != nil {
			slog.Error("PostHandler.handleGetAll - PostStore.GetAll", "error", err)
			responses.InternalServerResponse(w, ErrInternalServer)
			return
		}
		if err := h.cache.Set(ctx, key, posts, time.Minute*1).Err(); err != nil {
			slog.Error("PostHandler.handleGetAll - Cache.Set", "error", err)
			responses.InternalServerResponse(w, ErrInternalServer)
			return
		}
	}
	if err != nil {
		slog.Error("PostHandler.handleGetAll - Cache.GetAllPosts", "error", err)
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
	var post *types.Post
	key := generatePostKey(postID)
	post, err = h.cache.GetPost(ctx, key)
	if errors.Is(err, redis.Nil) {
		post, err = h.store.GetByID(ctx, postID)
		if err != nil {
			if errors.Is(err, store.ErrPostNotFound) {
				responses.NotFoundResponse(w, store.ErrPostNotFound)
				return
			}
			slog.Error("PostHandler.handleDeletePost - PostStore.GetByID", "error", err)
			responses.InternalServerResponse(w, ErrInternalServer)
			return
		}
	}
	if err != nil {
		slog.Error("PostHandler.handleDeletePost - Cache.GetPost", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}
	if post.UserID != user.ID {
		responses.ForbiddenResponse(w, ErrAccessDenied)
		return
	}

	err = h.store.Delete(ctx, postID)
	if err != nil {
		slog.Error("PostHandler.handleDeletePost - PostStore.Delete", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}
	if err := h.cache.Del(ctx, key).Err(); err != nil {
		slog.Error("PostHandler.handleDeletePost - Cache.Del", "error", err)
		responses.InternalServerResponse(w, ErrInternalServer)
		return
	}
	responses.JSON(w, http.StatusOK, envelope{"message": "post successfully deleted"})
}

func generatePostKey(id uuid.UUID) string {
	return fmt.Sprintf("post%s", id)
}
