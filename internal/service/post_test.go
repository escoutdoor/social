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

type postServiceSuite struct {
	suite.Suite
	container      testcontainers.Container
	redisContainer testcontainers.Container
	svc            Post
	authSvc        Auth
}

func (st *postServiceSuite) SetupTest() {
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
	st.svc = NewPostService(repo.Post, c)
	st.authSvc = NewAuthService(repo.Auth, repo.User, signKey)
}

func (st *postServiceSuite) TestGetByIDExistingPost() {
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
	post, err := st.svc.Create(ctx, userID, postIn)
	st.NoError(err, "failed to create post")
	st.NotEmpty(post, "expected to get post")

	p, err := st.svc.GetByID(ctx, post.ID)
	st.NoError(err, "failed to get post")
	st.NotEmpty(p, "expected to get post")
}

func (st *postServiceSuite) TestGetByIDNotFound() {
	ctx := context.Background()

	p, err := st.svc.GetByID(ctx, uuid.New())
	st.Error(err, "expected to get error: post not found")
	st.ErrorIs(err, repoerrs.ErrPostNotFound, "expected to get post not found error")
	st.Empty(p, "expected to get no data")
}

func (st *postServiceSuite) TestDeleteNotFound() {
	ctx := context.Background()

	err := st.svc.Delete(ctx, uuid.New(), uuid.New())
	st.Error(err, "expected to get error: post not found")
	st.ErrorIs(err, repoerrs.ErrPostNotFound, "expected to get post not found error")
}

func (st *postServiceSuite) TestDeleteExistingPost() {
	ctx := context.Background()

	in := types.CreateUserReq{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Password:  randomPw(),
	}
	userID, err := st.authSvc.SignUp(ctx, in)
	st.NoError(err, "failed to signup")
	st.NotEmpty(userID, "expected to get user id")

	postIn := types.CreatePostReq{
		Content:  gofakeit.Dessert(),
		PhotoURL: gofakeit.URL(),
	}
	post, err := st.svc.Create(ctx, userID, postIn)
	st.NoError(err, "failed to create post")
	st.NotEmpty(post, "expected to get post")

	err = st.svc.Delete(ctx, post.ID, userID)
	st.NoError(err, "failed to delete post")
}

func (st *postServiceSuite) TestUpdateExistingPost() {
	ctx := context.Background()

	in := types.CreateUserReq{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Password:  randomPw(),
	}
	userID, err := st.authSvc.SignUp(ctx, in)
	st.NoError(err, "failed to signup")
	st.NotEmpty(userID, "expected to get user id")

	postIn := types.CreatePostReq{
		Content:  gofakeit.Dessert(),
		PhotoURL: gofakeit.URL(),
	}
	post, err := st.svc.Create(ctx, userID, postIn)
	st.NoError(err, "failed to create post")
	st.NotEmpty(post, "expected to get post")

	updateIn := types.UpdatePostReq{
		Content:  strToPtr(gofakeit.CarModel()),
		PhotoURL: strToPtr(gofakeit.URL()),
	}
	updatedPost, err := st.svc.Update(ctx, post.ID, userID, updateIn)
	st.NoError(err, "failed to update post")
	st.NotEmpty(updatedPost, "expected to get post")

	st.Equal(*updateIn.Content, updatedPost.Content, "post content: expected %s, got %s", *updateIn.Content, updatedPost.Content)
	if updatedPost.PhotoURL != nil {
		st.Equal(*updateIn.PhotoURL, *updatedPost.PhotoURL, "post photo url: expected %s, got %s", *updateIn.PhotoURL, *updatedPost.PhotoURL)
	} else {
		st.Fail("updatedPost photo url is nil, but expected a value")
	}
}

func (st *postServiceSuite) TestUpdateNotFoundPost() {
	ctx := context.Background()

	updateIn := types.UpdatePostReq{
		Content:  strToPtr(gofakeit.CarModel()),
		PhotoURL: strToPtr(gofakeit.URL()),
	}
	updatedPost, err := st.svc.Update(ctx, uuid.New(), uuid.New(), updateIn)
	st.Error(err, "expected to get error: post not found")
	st.ErrorIs(err, repoerrs.ErrPostNotFound, "expected to get post not found error")
	st.Empty(updatedPost, "expected to get no post data")
}

func TestPostService(t *testing.T) {
	suite.Run(t, new(postServiceSuite))
}
