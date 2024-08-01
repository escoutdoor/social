package store

import (
	"context"
	"database/sql"

	"github.com/escoutdoor/social/internal/types"
	"github.com/google/uuid"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (s *UserStore) GetByID(ctx context.Context, id uuid.UUID) (*types.User, error) {
	stmt, err := s.db.PrepareContext(ctx, `
		SELECT * FROM USERS WHERE ID = $1
	`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		return scanUser(rows)
	}
	return nil, ErrUserNotFound
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

func (s *UserStore) Update(ctx context.Context, id uuid.UUID, input types.User) (*types.User, error) {
	stmt, err := s.db.PrepareContext(ctx, `
		UPDATE USERS SET
			FIRST_NAME = $1,
			LAST_NAME = $2,
			EMAIL = $3,
			PASSWORD = $4,
			DATE_OF_BIRTH = $5,
			BIO = $6,
			AVATAR_URL = $7
		WHERE ID = $8
	`)
	if err != nil {
		return nil, err
	}

	_, err = stmt.ExecContext(ctx,
		input.FirstName,
		input.LastName,
		input.Email,
		input.Password,
		input.DOB,
		input.Bio,
		input.AvatarURL,
		id,
	)
	if err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

func (s *UserStore) Delete(ctx context.Context, id uuid.UUID) error {
	stmt, err := s.db.PrepareContext(ctx, `
		DELETE FROM USERS WHERE ID = $1
	`)
	if err != nil {
		return err
	}

	result, err := stmt.ExecContext(ctx, id)
	if v, _ := result.RowsAffected(); v == 0 {
		return ErrUserNotFound
	}
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
		&user.DOB,
		&user.Bio,
		&user.AvatarURL,
		&user.UpdatedAt,
		&user.CreatedAt,
	); err != nil {
		return nil, err
	}
	return &user, nil

}
