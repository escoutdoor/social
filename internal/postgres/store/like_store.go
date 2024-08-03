package store

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type LikeStore struct {
	db *sql.DB
}

func NewLikeStore(db *sql.DB) *LikeStore {
	return &LikeStore{
		db: db,
	}
}

func (s *LikeStore) Like(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error {
	stmt, err := s.db.PrepareContext(ctx, `
		INSERT INTO LIKES(POST_ID, USER_ID) VALUES ($1, $2)
	`)
	if err != nil {
		return err
	}
	res, err := stmt.ExecContext(ctx, postID, userID)
	if err != nil {
		return err
	}
	if ra, _ := res.RowsAffected(); ra == 0 {
		return ErrFailedToLike
	}
	return nil
}

func (s *LikeStore) RemoveLike(ctx context.Context, postID, userID uuid.UUID) error {
	stmt, err := s.db.PrepareContext(ctx, `
		DELETE FROM LIKES WHERE POST_ID = $1 AND USER_ID = $2
	`)
	if err != nil {
		return err
	}
	res, err := stmt.ExecContext(ctx, postID, userID)
	if err != nil {
		return err
	}
	if ra, _ := res.RowsAffected(); ra == 0 {
		return ErrFailedToRemoveLike
	}
	return nil
}

func (s *LikeStore) IsLiked(ctx context.Context, postID uuid.UUID) (bool, error) {
	stmt, err := s.db.PrepareContext(ctx, `SELECT COUNT(*) FROM POST_LIKES WHERE POST_ID = $1`)
	if err != nil {
		return false, err
	}
	var count int
	if err = stmt.QueryRowContext(ctx, postID).Scan(&count); err != nil {
		return false, err
	}
	return count != 0, nil
}
