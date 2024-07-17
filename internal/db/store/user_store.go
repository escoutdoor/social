package store

import (
	"database/sql"

	"github.com/escoutdoor/social/internal/types"
)

type UserStorer interface {
	GetByID(id int) (*types.User, error)
}

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (s *UserStore) GetByID(id int) (*types.User, error) {
	return nil, nil
}
