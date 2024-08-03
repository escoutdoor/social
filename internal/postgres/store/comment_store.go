package store

import (
	"context"
	"database/sql"

	"github.com/escoutdoor/social/internal/types"
	"github.com/google/uuid"
)

type CommentStore struct {
	db *sql.DB
}

func NewCommentStore(db *sql.DB) *CommentStore {
	return &CommentStore{
		db: db,
	}
}

func (s *CommentStore) Create(ctx context.Context, userID uuid.UUID, postID uuid.UUID, input types.CreateCommentReq) (uuid.UUID, error) {
	stmt, err := s.db.PrepareContext(ctx, `
		INSERT INTO COMMENTS(CONTENT, POST_ID, USER_ID)
		VALUES($1, $2, $3)
		RETURNING ID
	`)
	if err != nil {
		return uuid.Nil, err
	}

	var commentID uuid.UUID
	err = stmt.QueryRow(input.Content, postID, userID).Scan(&commentID)
	if err != nil {
		return uuid.Nil, err
	}
	return commentID, nil
}

func (s *CommentStore) GetByID(ctx context.Context, id uuid.UUID) (*types.Comment, error) {
	stmt, err := s.db.PrepareContext(ctx, `
		SELECT * FROM COMMENTS WHERE ID = $1 
	`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		return scanComment(rows)
	}
	return nil, ErrCommentNotFound
}

func (s *CommentStore) Delete(ctx context.Context, id uuid.UUID) error {
	stmt, err := s.db.PrepareContext(ctx, `
		DELETE FROM COMMENTS WHERE ID = $1
	`)
	if err != nil {
		return err
	}

	result, err := stmt.ExecContext(ctx, id)
	if v, _ := result.RowsAffected(); v == 0 {
		return ErrCommentNotFound
	}
	return nil
}

func scanComment(rows *sql.Rows) (*types.Comment, error) {
	var comment types.Comment
	err := rows.Scan(
		&comment.ID,
		&comment.Content,
		&comment.UserID,
		&comment.PostID,
		&comment.UpdatedAt,
		&comment.CreatedAt,
	)
	return &comment, err
}
