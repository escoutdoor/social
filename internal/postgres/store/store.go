package store

import (
	"context"
	"database/sql"

	"github.com/escoutdoor/social/internal/config"
	"github.com/escoutdoor/social/internal/types"
	"github.com/google/uuid"
)

type Store struct {
	Auth    AuthStorer
	User    UserStorer
	Post    PostStorer
	Like    LikeStorer
	Comment CommentStorer
}

type AuthStorer interface {
	Create(ctx context.Context, input types.CreateUserReq) (uuid.UUID, error)
}

type UserStorer interface {
	GetByID(ctx context.Context, id uuid.UUID) (*types.User, error)
	GetByEmail(ctx context.Context, email string) (*types.User, error)
	Update(ctx context.Context, input types.User) (*types.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type PostStorer interface {
	Create(ctx context.Context, userID uuid.UUID, input types.CreatePostReq) (*types.Post, error)
	Update(ctx context.Context, postID uuid.UUID, input types.Post) (*types.Post, error)
	GetByID(ctx context.Context, id uuid.UUID) (*types.Post, error)
	GetAll(ctx context.Context) ([]types.Post, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type LikeStorer interface {
	IsPostLiked(ctx context.Context, postID uuid.UUID) (bool, error)
	IsCommentLiked(ctx context.Context, commentID uuid.UUID) (bool, error)
	LikePost(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error
	LikeComment(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) error
	RemoveLikeFromPost(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error
	RemoveLikeFromComment(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) error
}

type CommentStorer interface {
	Create(ctx context.Context, userID uuid.UUID, postID uuid.UUID, input types.CreateCommentReq) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*types.Comment, error)
	GetAll(ctx context.Context, postID uuid.UUID) ([]types.Comment, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

func NewStore(db *sql.DB, cfg *config.Config) *Store {
	return &Store{
		Auth:    NewAuthStore(db, cfg.SignKey),
		User:    NewUserStore(db),
		Post:    NewPostStore(db),
		Like:    NewLikeStore(db),
		Comment: NewCommentStore(db),
	}
}
