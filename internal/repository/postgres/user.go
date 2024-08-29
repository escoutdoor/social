package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/escoutdoor/social/internal/repository/repoerrs"
	"github.com/escoutdoor/social/internal/types"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (s *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*types.User, error) {
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
	return nil, repoerrs.ErrUserNotFound
}

func (s *UserRepository) GetByEmail(ctx context.Context, email string) (*types.User, error) {
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
	return nil, repoerrs.ErrUserNotFound
}

func (s *UserRepository) Update(ctx context.Context, input types.User) (*types.User, error) {
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

	var dob time.Time
	if input.DOB != nil {
		dob = time.Time(*input.DOB)
	}

	args := []interface{}{
		input.FirstName,
		input.LastName,
		input.Email,
		input.Password,
		dob,
		input.Bio,
		input.AvatarURL,
		input.ID,
	}
	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		var errPq *pq.Error
		if errors.As(err, &errPq) && errPq.Code == "23505" {
			return nil, repoerrs.ErrEmailAlreadyExists
		}
		return nil, err
	}
	return s.GetByID(ctx, input.ID)
}

func (s *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	stmt, err := s.db.PrepareContext(ctx, `
		DELETE FROM USERS WHERE ID = $1
	`)
	if err != nil {
		return err
	}

	result, err := stmt.ExecContext(ctx, id)
	if v, _ := result.RowsAffected(); v == 0 {
		return repoerrs.ErrUserNotFound
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
