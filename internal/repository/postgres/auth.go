package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/escoutdoor/social/internal/repository/repoerrs"
	"github.com/escoutdoor/social/internal/types"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB, jwtKey string) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (s *AuthRepository) Create(ctx context.Context, input types.CreateUserReq) (uuid.UUID, error) {
	var id uuid.UUID
	stmt, err := s.db.PrepareContext(ctx, `
		INSERT INTO USERS(FIRST_NAME, LAST_NAME, EMAIL, PASSWORD)
		VALUES ($1, $2, $3, $4)
		RETURNING ID
	`)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return id, repoerrs.ErrUserAlreadyExists
		}
		return id, err
	}

	args := []interface{}{input.FirstName, input.LastName, input.Email, input.Password}
	err = stmt.QueryRowContext(ctx, args...).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}
