package service

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/escoutdoor/social/internal/repository/postgres"
	"github.com/escoutdoor/social/internal/testutils"
	"github.com/escoutdoor/social/internal/types"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

var (
	signKey = "test"
)

type authServiceSuite struct {
	suite.Suite
	container testcontainers.Container
	svc       Auth
}

func (st *authServiceSuite) SetupTest() {
	container, db, err := testutils.NewPostgresContainer()
	st.Require().NoError(err, "failed to run container")
	st.Require().NotEmpty(container, "expected to get non-empty container")
	st.Require().NotEmpty(db, "expected to get non-empty db connection")

	authRepo := postgres.NewAuthRepository(db)
	userRepo := postgres.NewUserRepository(db)

	st.container = container
	st.svc = NewAuthService(authRepo, userRepo, signKey)
}

func (st *authServiceSuite) TestSignUp() {
	ctx := context.Background()
	in := types.CreateUserReq{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Password:  randomPw(),
	}

	id, err := st.svc.SignUp(ctx, in)
	st.NoError(err, "failed to signup")
	st.NotEmpty(id, "id should be non-empty")
}

func (st *authServiceSuite) TestSignUpSignInParseToken() {
	ctx := context.Background()

	registerIn := types.CreateUserReq{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Password:  randomPw(),
	}
	id, err := st.svc.SignUp(ctx, registerIn)
	st.NoError(err, "failed to signup")
	st.NotEmpty(id, "id should be non-empty")

	loginIn := types.LoginReq{
		Email:    registerIn.Email,
		Password: registerIn.Password,
	}
	token, err := st.svc.SignIn(ctx, loginIn)
	st.NoError(err, "failed to signin")
	st.NotEmpty(token, "expected to get non-empty value")

	tokenID, err := st.svc.ParseToken(token)
	st.NoError(err, "failed to parse token")
	st.Equal(id, tokenID, "id's not equal")
}

func (st *authServiceSuite) TestSignUpWithExistingEmail() {
	ctx := context.Background()

	in := types.CreateUserReq{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Password:  randomPw(),
	}
	id, err := st.svc.SignUp(ctx, in)
	st.NoError(err, "failed to signup")
	st.NotEmpty(id, "id should be non-empty")

	id, err = st.svc.SignUp(ctx, in)
	st.Error(err, "expected to get error: user already exists")
	st.Empty(id, "id should be empty")
}

func (st *authServiceSuite) TestSignInWithFakeEmail() {
	ctx := context.Background()
	in := types.LoginReq{
		Email:    gofakeit.Email(),
		Password: randomPw(),
	}

	token, err := st.svc.SignIn(ctx, in)
	st.Error(err, "expected to get error: invalid email or password")
	st.Empty(token, "expected to get empty value")
}

func TestAuthService(t *testing.T) {
	suite.Run(t, new(authServiceSuite))
}

func randomPw() string {
	return gofakeit.Password(true, true, true, true, false, 6)
}
