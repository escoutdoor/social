package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/escoutdoor/social/internal/types"
	"github.com/escoutdoor/social/pkg/hasher"
	"github.com/golang-jwt/jwt"
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
		slog.Error("hasher.HashPw", "error", err)
		return uuid.Nil, err
	}

	var id uuid.UUID
	err = stmt.QueryRowContext(ctx,
		input.FirstName,
		input.LastName,
		input.Email,
		input.Password,
	).Scan(&id)
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
	jwt.StandardClaims
}

func (s *AuthStore) GenerateToken(ctx context.Context, userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &AuthTokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			IssuedAt:  time.Now().Unix(),
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
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*AuthTokenClaims)
	if !ok {
		return uuid.Nil, ErrInvalidToken
	}
	return claims.UserID, nil
}
