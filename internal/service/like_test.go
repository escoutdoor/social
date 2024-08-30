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

type likeServiceSuite struct {
	suite.Suite
	redisContainer testcontainers.Container
	container      testcontainers.Container
	svc            Like
	postSvc        Post
	commentSvc     Comment
	authSvc        Auth
}

func (st *likeServiceSuite) SetupSuite() {
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
	st.svc = NewLikeService(repo.Like, c)
	st.authSvc = NewAuthService(repo.Auth, repo.User, signKey)
	st.postSvc = NewPostService(repo.Post, c)
	st.commentSvc = NewCommentService(repo.Comment, repo.Post)
}

func (st *likeServiceSuite) TearDownSuite() {
	err := st.container.Terminate(context.Background())
	st.Require().NoError(err, "failed to terminate postgres container")

	err = st.redisContainer.Terminate(context.Background())
	st.Require().NoError(err, "failed to terminate redis container")
}

func (st *likeServiceSuite) TestLikeAndRemoveLikeFromPost() {
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
		Content:  gofakeit.Comment(),
		PhotoURL: gofakeit.URL(),
	}
	post, err := st.postSvc.Create(ctx, userID, postIn)
	st.NoError(err, "failed to create post")
	st.NotEmpty(post, "expected to get post")

	err = st.svc.LikePost(ctx, post.ID, userID)
	st.NoError(err, "failed to like post")

	err = st.svc.RemoveLikeFromPost(ctx, post.ID, userID)
	st.NoError(err, "failed to remove like from post")
}

func (st *likeServiceSuite) TestLikeNotExistingPost() {
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

	err = st.svc.LikePost(ctx, uuid.New(), userID)
	st.Error(err, "expected to get error: post not found")
	st.ErrorIs(err, repoerrs.ErrPostNotFound, "expected to get post not found error")
}

func (st *likeServiceSuite) TestRemoveLikeFromNotExistingPost() {
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

	err = st.svc.RemoveLikeFromPost(ctx, uuid.New(), userID)
	st.Error(err, "expected to get error: failed to remove like")
	st.ErrorIs(err, repoerrs.ErrRemoveLikeFailed, "expected to get failed to remove like error")
}

func (st *likeServiceSuite) TestLikeAlreadyLikedPost() {
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
		Content:  gofakeit.Comment(),
		PhotoURL: gofakeit.URL(),
	}
	post, err := st.postSvc.Create(ctx, userID, postIn)
	st.NoError(err, "failed to create post")
	st.NotEmpty(post, "expected to get post")

	err = st.svc.LikePost(ctx, post.ID, userID)
	st.NoError(err, "failed to like post")

	err = st.svc.LikePost(ctx, post.ID, userID)
	st.Error(err, "expected to get error: already liked by user")
	st.ErrorIs(err, ErrAlreadyLiked, "expected to get already like by user error")

}

func (st *likeServiceSuite) TestRemoveLikeFromNotLikedPost() {
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
		Content:  gofakeit.Comment(),
		PhotoURL: gofakeit.URL(),
	}
	post, err := st.postSvc.Create(ctx, userID, postIn)
	st.NoError(err, "failed to create post")
	st.NotEmpty(post, "expected to get post")

	err = st.svc.RemoveLikeFromPost(ctx, post.ID, userID)
	st.Error(err, "expected to get error: failed to remove like")
	st.ErrorIs(err, repoerrs.ErrRemoveLikeFailed, "expected to get failed to remove like error")
}

func (st *likeServiceSuite) TestLikeAndRemoveLikeFromComment() {
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
		Content:  gofakeit.Comment(),
		PhotoURL: gofakeit.URL(),
	}
	post, err := st.postSvc.Create(ctx, userID, postIn)
	st.NoError(err, "failed to create post")
	st.NotEmpty(post, "expected to get post")

	commentIn := types.CreateCommentReq{
		Content: gofakeit.Comment(),
	}
	commentID, err := st.commentSvc.Create(ctx, userID, post.ID, commentIn)
	st.NoError(err, "failed to create comment")
	st.NotEmpty(commentID, "expected to get comment id")

	err = st.svc.LikeComment(ctx, commentID, userID)
	st.NoError(err, "failed to like comment")

	err = st.svc.RemoveLikeFromComment(ctx, commentID, userID)
	st.NoError(err, "failed to remove like from comment")
}

func (st *likeServiceSuite) TestLikeNotExistingComment() {
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

	err = st.svc.LikeComment(ctx, uuid.New(), userID)
	st.Error(err, "expected to get error: comment not found")
	st.ErrorIs(err, repoerrs.ErrCommentNotFound, "expected to get comment not found error")
}

func (st *likeServiceSuite) TestRemoveLikeFromNotExistingComment() {
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

	err = st.svc.RemoveLikeFromComment(ctx, uuid.New(), userID)
	st.Error(err, "expected to get error: failed to remove like")
	st.ErrorIs(err, repoerrs.ErrRemoveLikeFailed, "expected to get failed to remove like error")
}

func (st *likeServiceSuite) TestLikeAlreadyLikedComment() {
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
		Content:  gofakeit.Comment(),
		PhotoURL: gofakeit.URL(),
	}
	post, err := st.postSvc.Create(ctx, userID, postIn)
	st.NoError(err, "failed to create post")
	st.NotEmpty(post, "expected to get post")

	commentIn := types.CreateCommentReq{
		Content: gofakeit.Comment(),
	}
	commentID, err := st.commentSvc.Create(ctx, userID, post.ID, commentIn)
	st.NoError(err, "failed to create comment")
	st.NotEmpty(commentID, "expected to get comment id")

	err = st.svc.LikeComment(ctx, commentID, userID)
	st.NoError(err, "failed to like comment")

	err = st.svc.LikeComment(ctx, commentID, userID)
	st.Error(err, "expected to get error: already like by user")
	st.ErrorIs(err, ErrAlreadyLiked, "expected to get already like by user error")
}

func (st *likeServiceSuite) TestRemoveLikeFromNotLikedComment() {
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
		Content:  gofakeit.Comment(),
		PhotoURL: gofakeit.URL(),
	}
	post, err := st.postSvc.Create(ctx, userID, postIn)
	st.NoError(err, "failed to create post")
	st.NotEmpty(post, "expected to get post")

	commentIn := types.CreateCommentReq{
		Content: gofakeit.Comment(),
	}
	commentID, err := st.commentSvc.Create(ctx, userID, post.ID, commentIn)
	st.NoError(err, "failed to create comment")
	st.NotEmpty(commentID, "expected to get comment id")

	err = st.svc.RemoveLikeFromComment(ctx, commentID, userID)
	st.Error(err, "expected to get error: failed to remove like")
	st.ErrorIs(err, repoerrs.ErrRemoveLikeFailed, "expected to get failed to remove like error")
}

func TestLikeService(t *testing.T) {
	suite.Run(t, new(likeServiceSuite))
}
