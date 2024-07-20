package store

import (
	"database/sql"

	"github.com/escoutdoor/social/internal/types"
)

type AuthStore struct {
	db *sql.DB
}

type AuthStorer interface {
	SignUp()
}

func NewAuthStore(db *sql.DB) *AuthStore {
	return &AuthStore{
		db: db,
	}
}

func (s *AuthStore) SignUp(input types.User) (int, error) {

	return 0, nil
}
