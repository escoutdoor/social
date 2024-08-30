package service

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/escoutdoor/social/internal/repository"
	"github.com/escoutdoor/social/internal/repository/repoerrs"
	"github.com/escoutdoor/social/internal/testutils"
	"github.com/escoutdoor/social/internal/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

type commentServiceSuite struct {
	suite.Suite
	container      testcontainers.Container
	redisContainer testcontainers.Container
	svc            Comment
	postSvc        Post
	authSvc        Auth
}

func (st *commentServiceSuite) SetupSuite() {
	container, db, err := testutils.NewPostgresContainer()
	st.Require().NoError(err, "failed to run postgres container")
	st.Require().NotEmpty(container, "expected to get postgres container")
	st.Require().NotEmpty(db, "expected to get db connection")

	redisContainer, c, err := testutils.NewRedisContainer()
	st.Require().NoError(err, "failed to run redis container")
	st.Require().NotEmpty(redisContainer, "expected to get redis container")
	st.Require().NotEmpty(c, "expected to get redis connection")

	repo := repository.New(db)

	st.container = container
	st.redisContainer = redisContainer
	st.svc = NewCommentService(repo.Comment, repo.Post)
	st.postSvc = NewPostService(repo.Post, c)
	st.authSvc = NewAuthService(repo.Auth, repo.User, signKey)
}

func (st *commentServiceSuite) TearDownSuite() {
	err := st.container.Terminate(context.Background())
	st.Require().NoError(err, "failed to terminate postgres container")

	err = st.redisContainer.Terminate(context.Background())
	st.Require().NoError(err, "failed to terminate redis container")
}

func (st *commentServiceSuite) TestCreateCommentExistingPost() {
	ctx := context.Background()

	signupIn := types.CreateUserReq{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Password:  randomPw(),
	}
	userID, err := st.authSvc.SignUp(ctx, signupIn)
	st.NoError(err, "failed to signup")
	st.NotEmpty(userID, "expected to get user id")

	postIn := types.CreatePostReq{
		Content:  gofakeit.Dessert(),
		PhotoURL: gofakeit.URL(),
	}
	post, err := st.postSvc.Create(ctx, userID, postIn)
	st.NoError(err, "failed to create post")
	st.NotEmpty(post, "expected to get post")

	commentIn := types.CreateCommentReq{
		Content: gofakeit.Comment(),
	}
	commentID, err := st.svc.Create(ctx, userID, post.ID, commentIn)
	st.NoError(err, "failed to create comment")
	st.NotEmpty(commentID, "expected to get comment id")
}

func (st *commentServiceSuite) TestCreateCommentNonExistingPost() {
	ctx := context.Background()

	signupIn := types.CreateUserReq{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Password:  randomPw(),
	}
	userID, err := st.authSvc.SignUp(ctx, signupIn)
	st.NoError(err, "failed to signup")
	st.NotEmpty(userID, "expected to get user id")

	commentIn := types.CreateCommentReq{
		Content: gofakeit.Comment(),
	}
	id, err := st.svc.Create(ctx, userID, uuid.New(), commentIn)
	st.Error(err, "expected to get error: post not found")
	st.ErrorIs(err, repoerrs.ErrPostNotFound, "expected to get post not found error")
	st.Empty(id, "expected to get no data")
}

func (st *commentServiceSuite) TestCreateCommentNonExistingUser() {
	ctx := context.Background()

	signupIn := types.CreateUserReq{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Password:  randomPw(),
	}
	userID, err := st.authSvc.SignUp(ctx, signupIn)
	st.NoError(err, "failed to signup")
	st.NotEmpty(userID, "expected to get user id")

	postIn := types.CreatePostReq{
		Content:  gofakeit.Dessert(),
		PhotoURL: gofakeit.URL(),
	}
	post, err := st.postSvc.Create(ctx, userID, postIn)
	st.NoError(err, "failed to create post")
	st.NotEmpty(post, "expected to get post")

	commentIn := types.CreateCommentReq{
		Content: gofakeit.Comment(),
	}
	commentID, err := st.svc.Create(ctx, uuid.New(), post.ID, commentIn)
	st.Error(err, "expected to get error: user not found")
	st.ErrorIs(err, repoerrs.ErrUserNotFound, "expected to get user not found error")
	st.Empty(commentID, "expected to get no data")
}

func (st *commentServiceSuite) TestCreateCommentOnComment() {
	ctx := context.Background()

	signupIn := types.CreateUserReq{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Password:  randomPw(),
	}
	userID, err := st.authSvc.SignUp(ctx, signupIn)
	st.NoError(err, "failed to signup")
	st.NotEmpty(userID, "expected to get user id")

	postIn := types.CreatePostReq{
		Content:  gofakeit.Dessert(),
		PhotoURL: gofakeit.URL(),
	}
	post, err := st.postSvc.Create(ctx, userID, postIn)
	st.NoError(err, "failed to create post")
	st.NotEmpty(post, "expected to get post")

	commentIn := types.CreateCommentReq{
		Content: gofakeit.Comment(),
	}
	commentID, err := st.svc.Create(ctx, userID, post.ID, commentIn)
	st.NoError(err, "failed to create parent comment")
	st.NotEmpty(commentID, "expected to get comment id")

	commentIDPtr := &commentID
	commentIn = types.CreateCommentReq{
		Content:         gofakeit.Comment(),
		ParentCommentID: commentIDPtr,
	}
	commentOnCommentID, err := st.svc.Create(ctx, userID, post.ID, commentIn)
	st.NoError(err, "failed to create comment")
	st.NotEmpty(commentOnCommentID, "expected to get comment id")
}

func (st *commentServiceSuite) TestGetByIDNotFound() {
	ctx := context.Background()

	comment, err := st.svc.GetByID(ctx, uuid.New())
	st.Error(err, "expected to get error: comment not found")
	st.ErrorIs(err, repoerrs.ErrCommentNotFound, "expected to get comment not found error")
	st.Empty(comment, "expected to get no data")
}

func (st *commentServiceSuite) TestGetByID() {
	ctx := context.Background()

	signupIn := types.CreateUserReq{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Password:  randomPw(),
	}
	userID, err := st.authSvc.SignUp(ctx, signupIn)
	st.NoError(err, "failed to signup")
	st.NotEmpty(userID, "expected to get user id")

	postIn := types.CreatePostReq{
		Content:  gofakeit.Dessert(),
		PhotoURL: gofakeit.URL(),
	}
	post, err := st.postSvc.Create(ctx, userID, postIn)
	st.NoError(err, "failed to create post")
	st.NotEmpty(post, "expected to get post")

	commentIn := types.CreateCommentReq{
		Content: gofakeit.Comment(),
	}
	commentID, err := st.svc.Create(ctx, userID, post.ID, commentIn)
	st.NoError(err, "failed to create comment")
	st.NotEmpty(commentID, "expected to get comment id")

	comment, err := st.svc.GetByID(ctx, commentID)
	st.NoError(err, "failed to get comment")
	st.NotEmpty(comment, "expected to get comment")
}

func (st *commentServiceSuite) TestDeleteComment() {
	ctx := context.Background()

	signupIn := types.CreateUserReq{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Password:  randomPw(),
	}
	userID, err := st.authSvc.SignUp(ctx, signupIn)
	st.NoError(err, "failed to signup")
	st.NotEmpty(userID, "expected to get user id")

	postIn := types.CreatePostReq{
		Content:  gofakeit.Dessert(),
		PhotoURL: gofakeit.URL(),
	}
	post, err := st.postSvc.Create(ctx, userID, postIn)
	st.NoError(err, "failed to create post")
	st.NotEmpty(post, "expected to get post")

	commentIn := types.CreateCommentReq{
		Content: gofakeit.Comment(),
	}
	commentID, err := st.svc.Create(ctx, userID, post.ID, commentIn)
	st.NoError(err, "failed to create comment")
	st.NotEmpty(commentID, "expected to get comment id")

	err = st.svc.Delete(ctx, commentID, userID)
	st.NoError(err, "failed to delete comment")
}

func (st *commentServiceSuite) TestDeleteNotExistingComment() {
	ctx := context.Background()

	err := st.svc.Delete(ctx, uuid.New(), uuid.New())
	st.Error(err, "expected to get error: comment not found")
	st.ErrorIs(err, repoerrs.ErrCommentNotFound, "expected to get comment not found error")
}

func TestCommentService(t *testing.T) {
	suite.Run(t, new(commentServiceSuite))
}
