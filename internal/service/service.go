package service

import (
	"context"
	"io"
	"mime/multipart"

	"github.com/escoutdoor/social/internal/cache"
	"github.com/escoutdoor/social/internal/repository"
	"github.com/escoutdoor/social/internal/s3"
	"github.com/escoutdoor/social/internal/types"
	"github.com/escoutdoor/social/pkg/validator"
	"github.com/google/uuid"
)

type Auth interface {
	ParseToken(jwtToken string) (uuid.UUID, error)
	SignIn(ctx context.Context, input types.LoginReq) (string, error)
	SignUp(ctx context.Context, input types.CreateUserReq) (uuid.UUID, error)
}

type User interface {
	GetByID(ctx context.Context, id uuid.UUID) (*types.User, error)
	Update(ctx context.Context, user types.User, input types.UpdateUserReq) (*types.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Post interface {
	Create(ctx context.Context, userID uuid.UUID, input types.CreatePostReq) (*types.Post, error)
	Update(ctx context.Context, postID uuid.UUID, userID uuid.UUID, input types.UpdatePostReq) (*types.Post, error)
	GetByID(ctx context.Context, id uuid.UUID) (*types.Post, error)
	GetAll(ctx context.Context) ([]types.Post, error)
	Delete(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error
}

type Comment interface {
	Create(ctx context.Context, userID uuid.UUID, postID uuid.UUID, input types.CreateCommentReq) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*types.Comment, error)
	GetAll(ctx context.Context, postID uuid.UUID) ([]types.Comment, error)
	Delete(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) error
}

type Like interface {
	LikePost(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error
	LikeComment(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) error
	RemoveLikeFromPost(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error
	RemoveLikeFromComment(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) error
}

type File interface {
	Create(ctx context.Context, src io.Reader, hdr *multipart.FileHeader) (string, error)
}

type Opts struct {
	Repository *repository.Repository
	Cache      cache.Repository
	S3         s3.Repository
	Validator  *validator.Validator

	SignKey string
}

func NewServices(opts Opts) *Services {
	return &Services{
		Auth:    NewAuthService(opts.Repository.Auth, opts.Repository.User, opts.SignKey),
		User:    NewUserService(opts.Repository.User, opts.Validator),
		Post:    NewPostService(opts.Repository.Post, opts.Cache),
		Comment: NewCommentService(opts.Repository.Comment, opts.Repository.Post),
		Like:    NewLikeService(opts.Repository.Like, opts.Cache),
		File:    NewFileService(opts.S3),
	}
}

type Services struct {
	Auth
	User
	Post
	Comment
	Like
	File
}
