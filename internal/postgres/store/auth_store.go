package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/escoutdoor/social/internal/types"
	"github.com/escoutdoor/social/pkg/hasher"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthStore struct {
	db        *sql.DB
	userStore UserStorer
	jwtKey    string
}

func NewAuthStore(db *sql.DB, jwtKey string) *AuthStore {
	return &AuthStore{
		db:        db,
		userStore: NewUserStore(db),
		jwtKey:    jwtKey,
	}
}

func (s *AuthStore) SignUp(ctx context.Context, input types.CreateUserReq) (uuid.UUID, error) {
	_, err := s.userStore.GetByEmail(ctx, input.Email)
	switch {
	case err != nil && !errors.Is(err, ErrUserNotFound):
		return uuid.Nil, err
	case err == nil:
		return uuid.Nil, ErrUserAlreadyExists
	}
	stmt, err := s.db.PrepareContext(ctx, `
		INSERT INTO USERS(FIRST_NAME, LAST_NAME, EMAIL, PASSWORD)
		VALUES ($1, $2, $3, $4)
		RETURNING ID
	`)
	if err != nil {
		return uuid.Nil, err
	}

	input.Password, err = hasher.HashPw(input.Password)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to hash password: %w", err)
	}

	args := []interface{}{input.FirstName, input.LastName, input.Email, input.Password}
	var id uuid.UUID
	err = stmt.QueryRowContext(ctx, args...).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (s *AuthStore) SignIn(ctx context.Context, input types.LoginReq) (*types.User, error) {
	user, err := s.userStore.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrInvalidEmailOrPw
		}
		return nil, err
	}

	if ok := hasher.ComparePw(input.Password, user.Password); !ok {
		return nil, ErrInvalidEmailOrPw
	}
	return user, nil
}

type AuthTokenClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *AuthStore) GenerateToken(ctx context.Context, userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &AuthTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: userID,
	})

	tokenStr, err := token.SignedString([]byte(s.jwtKey))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (s *AuthStore) ParseToken(jwtToken string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(jwtToken, &AuthTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(s.jwtKey), nil
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
