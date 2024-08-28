package service

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/escoutdoor/social/internal/repository/postgres"
	"github.com/escoutdoor/social/internal/testutils"
	"github.com/escoutdoor/social/internal/types"
	"github.com/escoutdoor/social/pkg/validator"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

type userServiceSuite struct {
	suite.Suite
	container testcontainers.Container
	svc       User
	authSvc   Auth
}

func (st *userServiceSuite) SetupTest() {
	container, db, err := testutils.NewPostgresContainer()
	st.Require().NoError(err, "failed to run container")
	st.Require().NotEmpty(container, "expected to get non-empty container")
	st.Require().NotEmpty(db, "expected to get non-empty db connection")

	repo := postgres.NewUserRepository(db)
	authRepo := postgres.NewAuthRepository(db)

	st.container = container
	st.svc = NewUserService(repo, validator.New())
	st.authSvc = NewAuthService(authRepo, repo, signKey)
}

func (st *userServiceSuite) TestGetByIDExistingUser() {
	ctx := context.Background()
	in := types.CreateUserReq{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Password:  randomPw(),
	}

	id, err := st.authSvc.SignUp(ctx, in)
	st.NoError(err, "failed to signup")
	st.NotEmpty(id, "id should be non-empty")

	u, err := st.svc.GetByID(ctx, id)
	st.NotEmpty(u, "user should be non-empty value")
	st.NoError(err, "failed to get user")
}

func (st *userServiceSuite) TestGetByIDNotFound() {
	ctx := context.Background()

	u, err := st.svc.GetByID(ctx, uuid.New())
	st.Error(err, "expected to get error: user not found")
	st.Empty(u, "user should be empty")
}

func (st *userServiceSuite) TestDeleteNotFound() {
	ctx := context.Background()

	err := st.svc.Delete(ctx, uuid.New())
	st.Error(err, "expected to get error: user not found")
}

func (st *userServiceSuite) TestDeleteExistingUser() {
	ctx := context.Background()
	in := types.CreateUserReq{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Password:  randomPw(),
	}

	id, err := st.authSvc.SignUp(ctx, in)
	st.NotEmpty(id, "id should be non-empty")
	st.NoError(err, "failed to signup")

	err = st.svc.Delete(ctx, id)
	st.NoError(err, "failed to delete user")
}

func (st *userServiceSuite) TestUpdateNotFound() {
	ctx := context.Background()

	dob := types.DOB(time.Now())
	user := types.User{
		ID:        uuid.New(),
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Password:  randomPw(),
		Bio:       strToPtr(gofakeit.Letter()),
		DOB:       &dob,
	}
	in := types.UpdateUserReq{
		FirstName: strToPtr(gofakeit.FirstName()),
		LastName:  strToPtr(gofakeit.LastName()),
		Email:     strToPtr(gofakeit.Email()),
		Password:  strToPtr(randomPw()),
		Bio:       strToPtr(gofakeit.Letter()),
	}

	_, err := st.svc.Update(ctx, user, in)
	st.Error(err, "expected to get error: user not found")
}

func (st *userServiceSuite) TestGetByIDUpdate() {
	ctx := context.Background()
	in := types.CreateUserReq{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Password:  randomPw(),
	}

	id, err := st.authSvc.SignUp(ctx, in)
	st.NoError(err, "failed to signup")
	st.NotEmpty(id, "id should be non-empty")

	user, err := st.svc.GetByID(ctx, id)
	st.NotEmpty(user, "user should be non-empty value")
	st.NoError(err, "failed to get user")

	updateIn := types.UpdateUserReq{
		FirstName: strToPtr(gofakeit.FirstName()),
		LastName:  strToPtr(gofakeit.LastName()),
		Email:     strToPtr(gofakeit.Email()),
	}
	u, err := st.svc.Update(ctx, *user, updateIn)
	st.NoError(err, "failed to update user")
	st.NotEmpty(u, "expected to get user")

	st.Equal(*updateIn.FirstName, u.FirstName, "user first name: expected %s, got %s", *updateIn.FirstName, u.FirstName)
	st.Equal(*updateIn.LastName, u.LastName, "user last name: expected %s, got %s", *updateIn.LastName, u.LastName)
	st.Equal(*updateIn.Email, u.Email, "user email: expected %s, got %s", *updateIn.Email, u.Email)
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(userServiceSuite))
}

func (st *userServiceSuite) AfterTest(suiteName, testName string) {
	st.container.Terminate(context.Background())
}

func strToPtr(s string) *string {
	return &s
}
