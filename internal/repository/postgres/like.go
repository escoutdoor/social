package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/escoutdoor/social/internal/repository/repoerrs"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type LikeRepository struct {
	db *sql.DB
}

func NewLikeRepository(db *sql.DB) *LikeRepository {
	return &LikeRepository{
		db: db,
	}
}

func (s *LikeRepository) LikePost(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error {
	stmt, err := s.db.PrepareContext(ctx, `
		INSERT INTO POST_LIKES(POST_ID, USER_ID) VALUES ($1, $2)
	`)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, postID, userID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23503" {
			return repoerrs.ErrPostNotFound
		}
		return err
	}
	if ra, _ := res.RowsAffected(); ra == 0 {
		return repoerrs.ErrLikeFailed
	}
	return nil
}

func (s *LikeRepository) RemoveLikeFromPost(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error {
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
		return repoerrs.ErrRemoveLikeFailed
	}
	return nil
}

func (s *LikeRepository) IsPostLiked(ctx context.Context, postID uuid.UUID) (bool, error) {
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

func (s *LikeRepository) LikeComment(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) error {
	stmt, err := s.db.PrepareContext(ctx, `
		INSERT INTO COMMENT_LIKES(COMMENT_ID, USER_ID) VALUES ($1, $2)
	`)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, commentID, userID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23503" {
			return repoerrs.ErrCommentNotFound
		}
		return err
	}
	if ra, _ := res.RowsAffected(); ra == 0 {
		return repoerrs.ErrLikeFailed
	}
	return nil
}

func (s *LikeRepository) RemoveLikeFromComment(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) error {
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
		return repoerrs.ErrRemoveLikeFailed
	}
	return nil
}

func (s *LikeRepository) IsCommentLiked(ctx context.Context, commentID uuid.UUID) (bool, error) {
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
