package repository

import (
	"context"
	"database/sql"

	"github.com/escoutdoor/social/internal/config"
	"github.com/escoutdoor/social/internal/repository/postgres"
	"github.com/escoutdoor/social/internal/types"
	"github.com/google/uuid"
)

type Auth interface {
	Create(ctx context.Context, input types.CreateUserReq) (uuid.UUID, error)
}

type User interface {
	GetByID(ctx context.Context, id uuid.UUID) (*types.User, error)
	GetByEmail(ctx context.Context, email string) (*types.User, error)
	Update(ctx context.Context, input types.User) (*types.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Post interface {
	Create(ctx context.Context, userID uuid.UUID, input types.CreatePostReq) (*types.Post, error)
	Update(ctx context.Context, postID uuid.UUID, input types.Post) (*types.Post, error)
	GetByID(ctx context.Context, id uuid.UUID) (*types.Post, error)
	GetAll(ctx context.Context) ([]types.Post, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Like interface {
	IsPostLiked(ctx context.Context, postID uuid.UUID) (bool, error)
	IsCommentLiked(ctx context.Context, commentID uuid.UUID) (bool, error)
	LikePost(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error
	LikeComment(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) error
	RemoveLikeFromPost(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error
	RemoveLikeFromComment(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) error
}

type Comment interface {
	Create(ctx context.Context, userID uuid.UUID, postID uuid.UUID, input types.CreateCommentReq) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*types.Comment, error)
	GetAll(ctx context.Context, postID uuid.UUID) ([]types.Comment, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

func New(db *sql.DB, cfg *config.Config) *Repository {
	return &Repository{
		Auth:    postgres.NewAuthRepository(db, cfg.SignKey),
		User:    postgres.NewUserRepository(db),
		Post:    postgres.NewPostRepository(db),
		Like:    postgres.NewLikeRepository(db),
		Comment: postgres.NewCommentRepository(db),
	}
}

type Repository struct {
	Auth
	User
	Post
	Like
	Comment
}
