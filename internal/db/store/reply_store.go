package store

import (
	"context"
	"database/sql"

	"github.com/escoutdoor/social/internal/types"
	"github.com/google/uuid"
)

type ReplyStore struct {
	db *sql.DB
}

type ReplyStorer interface {
	Create(ctx context.Context, userID uuid.UUID, postID uuid.UUID, input types.CreateReplyReq) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*types.Reply, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

func NewReplyStore(db *sql.DB) *ReplyStore {
	return &ReplyStore{
		db: db,
	}
}

func (s *ReplyStore) Create(ctx context.Context, userID uuid.UUID, postID uuid.UUID, input types.CreateReplyReq) (uuid.UUID, error) {
	stmt, err := s.db.PrepareContext(ctx, `
		INSERT INTO REPLIES(TEXT, POST_ID, USER_ID)
		VALUES($1, $2, $3)
		RETURNING ID
	`)
	if err != nil {
		return uuid.Nil, err
	}

	var replyID uuid.UUID
	err = stmt.QueryRow(input.Text, postID, userID).Scan(&replyID)
	if err != nil {
		return uuid.Nil, err
	}
	return replyID, nil
}

func (s *ReplyStore) GetByID(ctx context.Context, id uuid.UUID) (*types.Reply, error) {
	stmt, err := s.db.PrepareContext(ctx, `
		SELECT * FROM REPLIES WHERE ID = $1 
	`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		return scanReply(rows)
	}
	return nil, ErrReplyNotFound
}

func (s *ReplyStore) Delete(ctx context.Context, id uuid.UUID) error {
	stmt, err := s.db.PrepareContext(ctx, `
		DELETE FROM REPLIES WHERE ID = $1
	`)
	if err != nil {
		return err
	}

	result, err := stmt.ExecContext(ctx, id)
	if v, _ := result.RowsAffected(); v == 0 {
		return ErrReplyNotFound
	}
	return nil
}

func scanReply(rows *sql.Rows) (*types.Reply, error) {
	var reply types.Reply
	err := rows.Scan(
		&reply.ID,
		&reply.Text,
		&reply.UserID,
		&reply.PostID,
		&reply.UpdatedAt,
		&reply.CreatedAt,
	)
	return &reply, err
}
