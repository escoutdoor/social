package test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/escoutdoor/social/internal/repository/postgres"
	"github.com/escoutdoor/social/internal/service"
	"github.com/escoutdoor/social/internal/types"
	"github.com/escoutdoor/social/test/suite"
	"github.com/stretchr/testify/require"
)

var (
	signKey = "test"
)

func TestSignUp(t *testing.T) {
	st, err := suite.New()
	require.NoError(t, err)

	ctx := context.Background()
	defer st.Container.Terminate(ctx)

	authRepo := postgres.NewAuthRepository(st.DB)
	userRepo := postgres.NewUserRepository(st.DB)
	svc := service.NewAuthService(authRepo, userRepo, signKey)

	in := types.CreateUserReq{
		FirstName: gofakeit.FirstName(),
		LastName:  gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Password:  randomPw(),
	}
	id, err := svc.SignUp(ctx, in)
	require.NoError(t, err)
	require.NotEmpty(t, id)
}

func randomPw() string {
	return gofakeit.Password(true, true, true, true, false, 6)
}
