package store

import (
	"context"
	"database/sql"

	"github.com/escoutdoor/social/internal/config"
	"github.com/escoutdoor/social/internal/types"
	"github.com/google/uuid"
)

type Store struct {
	Auth  AuthStorer
	User  UserStorer
	Post  PostStorer
	Reply ReplyStorer
}

type AuthStorer interface {
	SignUp(ctx context.Context, input types.CreateUserReq) (uuid.UUID, error)
	SignIn(ctx context.Context, input types.LoginReq) (*types.User, error)
	GenerateToken(ctx context.Context, userID uuid.UUID) (string, error)
	ParseToken(jwtToken string) (uuid.UUID, error)
}

type UserStorer interface {
	GetByID(ctx context.Context, id uuid.UUID) (*types.User, error)
	GetByEmail(ctx context.Context, email string) (*types.User, error)
	Update(ctx context.Context, id uuid.UUID, input types.User) (*types.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type PostStorer interface {
	Create(ctx context.Context, userID uuid.UUID, input types.CreatePostReq) (uuid.UUID, error)
	Update(ctx context.Context, postID uuid.UUID, input types.Post) (*types.Post, error)
	GetByID(ctx context.Context, id uuid.UUID) (*types.Post, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type ReplyStorer interface {
	Create(ctx context.Context, userID uuid.UUID, postID uuid.UUID, input types.CreateReplyReq) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*types.Reply, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

func NewStore(db *sql.DB, cfg config.Config) *Store {
	return &Store{
		Auth:  NewAuthStore(db, cfg.JWTKey),
		User:  NewUserStore(db),
		Post:  NewPostStore(db),
		Reply: NewReplyStore(db),
	}
}
