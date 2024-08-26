package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/escoutdoor/social/internal/repository"
	"github.com/escoutdoor/social/internal/repository/repoerrs"
	"github.com/escoutdoor/social/internal/types"
	"github.com/escoutdoor/social/pkg/hasher"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService struct {
	repo     repository.Auth
	userRepo repository.User
	signKey  string
}

func NewAuthService(repo repository.Auth, userRepo repository.User, signKey string) *AuthService {
	return &AuthService{
		repo:     repo,
		userRepo: userRepo,
		signKey:  signKey,
	}
}

func (s *AuthService) SignUp(ctx context.Context, input types.CreateUserReq) (uuid.UUID, error) {
	var err error
	input.Password, err = hasher.HashPw(input.Password)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to hash password: %w", err)
	}

	return s.repo.Create(ctx, input)
}

func (s *AuthService) SignIn(ctx context.Context, input types.LoginReq) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, repoerrs.ErrUserNotFound) {
			return "", ErrInvalidEmailOrPw
		}
		return "", err
	}

	if ok := hasher.ComparePw(input.Password, user.Password); !ok {
		return "", ErrInvalidEmailOrPw
	}
	token, err := s.generateToken(ctx, user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

type AuthTokenClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *AuthService) generateToken(_ context.Context, userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &AuthTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: userID,
	})

	tokenStr, err := token.SignedString([]byte(s.signKey))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (s *AuthService) ParseToken(jwtToken string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(jwtToken, &AuthTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(s.signKey), nil
	})
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			return uuid.Nil, fmt.Errorf("it doesn't look like a token")
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return uuid.Nil, fmt.Errorf("invalid token signature")
		case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
			return uuid.Nil, fmt.Errorf("token is either expired or not active yet")
		default:
			return uuid.Nil, ErrInvalidToken
		}
	}
	claims, ok := token.Claims.(*AuthTokenClaims)
	if !ok {
		return uuid.Nil, err
	}
	return claims.UserID, nil
}
