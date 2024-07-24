package store

import (
	"context"
	"database/sql"

	"github.com/escoutdoor/social/internal/types"
	"github.com/google/uuid"
)

type UserStorer interface {
	GetByID(ctx context.Context, id uuid.UUID) (*types.User, error)
	GetByEmail(ctx context.Context, email string) (*types.User, error)
	Update(ctx context.Context, id uuid.UUID, input types.UpdateUserReq) (*types.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (s *UserStore) GetByID(ctx context.Context, id uuid.UUID) (*types.User, error) {
	return nil, nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*types.User, error) {
	stmt, err := s.db.PrepareContext(ctx, `
		SELECT * FROM USERS WHERE EMAIL = $1
	`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryContext(ctx, email)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		return scanUser(rows)
	}
	return nil, ErrUserNotFound
}

func (s *UserStore) Update(ctx context.Context, id uuid.UUID, input types.UpdateUserReq) (*types.User, error) {
	return nil, nil
}

func (s *UserStore) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func scanUser(rows *sql.Rows) (*types.User, error) {
	var user types.User
	if err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.BirthDate,
		&user.Bio,
		&user.UpdatedAt,
		&user.CreatedAt,
	); err != nil {
		return nil, err
	}
	return &user, nil

}
