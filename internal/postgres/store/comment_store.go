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
	var id uuid.UUID
	stmt, err := s.db.PrepareContext(ctx, `
		INSERT INTO COMMENTS(CONTENT, PARENT_COMMENT_ID, POST_ID, USER_ID)
		VALUES($1, $2, $3, $4)
		RETURNING ID
	`)
	if err != nil {
		return uuid.Nil, err
	}

	err = stmt.QueryRow(input.Content, input.ParentCommentID, postID, userID).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
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

func (s *CommentStore) GetAll(ctx context.Context, postID uuid.UUID) ([]types.Comment, error) {
	var comments []types.Comment
	stmt, err := s.db.PrepareContext(ctx, `
		WITH RECURSIVE comments_cte AS (
			SELECT 
				id,
				content,
				user_id,
				post_id,
				parent_comment_id,
				created_at,
				updated_at,
				1 AS depth
			FROM COMMENTS 
			WHERE post_id = $1 AND parent_comment_id IS NULL
			UNION ALL
			SELECT 
				c.id,
				c.content,
				c.user_id,
				c.post_id,
				c.parent_comment_id,
				c.created_at,
				c.updated_at,
				cte.depth + 1
			FROM COMMENTS c
			JOIN comments_cte cte ON c.parent_comment_id = cte.id
		)
		SELECT
			id,
			content,
			user_id,
			post_id,
			parent_comment_id,
			created_at,
			updated_at,
			depth
		FROM comments_cte ORDER BY depth, created_at;
	`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryContext(ctx, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			c     types.Comment
			pcid  uuid.NullUUID
			depth int
		)
		if err := rows.Scan(
			&c.ID,
			&c.Content,
			&c.UserID,
			&c.PostID,
			&pcid,
			&c.CreatedAt,
			&c.UpdatedAt,
			&depth,
		); err != nil {
			return nil, err
		}

		if pcid.Valid {
			c.ParentCommentID = &pcid.UUID
		}
		comments = append(comments, c)
	}

	commentMap := make(map[uuid.UUID]types.Comment)
	var root []types.Comment
	for _, item := range comments {
		commentMap[item.ID] = item
		if item.ParentCommentID != nil {
			parent := commentMap[*item.ParentCommentID]
			parent.Replies = append(parent.Replies, item)
			continue
		}
		root = append(root, item)
	}
	return root, nil
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
		&comment.ParentCommentID,
		&comment.UpdatedAt,
		&comment.CreatedAt,
	)
	return &comment, err
}
