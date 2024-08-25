package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type LikeStore struct {
	db *sql.DB
}

func NewLikeStore(db *sql.DB) *LikeStore {
	return &LikeStore{
		db: db,
	}
}

func (s *LikeStore) LikePost(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error {
	stmt, err := s.db.PrepareContext(ctx, `
		INSERT INTO POST_LIKES(POST_ID, USER_ID) VALUES ($1, $2)
	`)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return ErrPostNotFound
		}
		return err
	}
	res, err := stmt.ExecContext(ctx, postID, userID)
	if err != nil {
		return err
	}
	if ra, _ := res.RowsAffected(); ra == 0 {
		return ErrLikeFailed
	}
	return nil
}

func (s *LikeStore) RemoveLikeFromPost(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error {
	stmt, err := s.db.PrepareContext(ctx, `
		DELETE FROM POST_LIKES WHERE POST_ID = $1 AND USER_ID = $2
	`)
	if err != nil {
		return err
	}
	res, err := stmt.ExecContext(ctx, postID, userID)
	if err != nil {
		return err
	}
	if ra, _ := res.RowsAffected(); ra == 0 {
		return ErrRemoveLikeFailed
	}
	return nil
}

func (s *LikeStore) IsPostLiked(ctx context.Context, postID uuid.UUID) (bool, error) {
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

func (s *LikeStore) LikeComment(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) error {
	stmt, err := s.db.PrepareContext(ctx, `
		INSERT INTO COMMENT_LIKES(COMMENT_ID, USER_ID) VALUES ($1, $2)
	`)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return ErrCommentNotFound
		}
		return err
	}
	res, err := stmt.ExecContext(ctx, commentID, userID)
	if err != nil {
		return err
	}
	if ra, _ := res.RowsAffected(); ra == 0 {
		return ErrLikeFailed
	}
	return nil
}

func (s *LikeStore) RemoveLikeFromComment(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) error {
	stmt, err := s.db.PrepareContext(ctx, `
		DELETE FROM COMMENT_LIKES WHERE COMMENT_ID = $1 AND USER_ID = $2
	`)
	if err != nil {
		return err
	}
	res, err := stmt.ExecContext(ctx, commentID, userID)
	if err != nil {
		return err
	}
	if ra, _ := res.RowsAffected(); ra == 0 {
		return ErrRemoveLikeFailed
	}
	return nil
}

func (s *LikeStore) IsCommentLiked(ctx context.Context, commentID uuid.UUID) (bool, error) {
	stmt, err := s.db.PrepareContext(ctx, `SELECT COUNT(*) FROM COMMENT_LIKES WHERE COMMENT_ID = $1`)
	if err != nil {
		return false, err
	}
	var count int
	if err = stmt.QueryRowContext(ctx, commentID).Scan(&count); err != nil {
		return false, err
	}
	return count != 0, nil
}
