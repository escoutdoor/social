package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/escoutdoor/social/internal/repository/repoerrs"
	"github.com/escoutdoor/social/internal/types"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

func (s *CommentRepository) Create(ctx context.Context, userID uuid.UUID, postID uuid.UUID, input types.CreateCommentReq) (uuid.UUID, error) {
	var id uuid.UUID
	stmt, err := s.db.PrepareContext(ctx, `
		INSERT INTO COMMENTS(CONTENT, PARENT_COMMENT_ID, POST_ID, USER_ID)
		VALUES($1, $2, $3, $4)
		RETURNING ID
	`)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			switch pqErr.Constraint {
			case "comments_post_id_fkey":
				return id, repoerrs.ErrPostNotFound
			case "comments_parent_comment_id_fkey":
				return id, repoerrs.ErrCommentNotFound
			}
		}
		return id, err
	}

	args := []interface{}{input.Content, input.ParentCommentID, postID, userID}
	err = stmt.QueryRowContext(ctx, args...).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (s *CommentRepository) GetByID(ctx context.Context, id uuid.UUID) (*types.Comment, error) {
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
	return nil, repoerrs.ErrCommentNotFound
}

func (s *CommentRepository) GetAll(ctx context.Context, postID uuid.UUID) ([]types.Comment, error) {
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

	var comments []types.Comment
	commentMap := make(map[uuid.UUID]*types.Comment)

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

		commentMap[c.ID] = &c
		if c.ParentCommentID != nil {
			parent := commentMap[*c.ParentCommentID]
			parent.Replies = append(parent.Replies, c)
		} else {
			comments = append(comments, c)
		}
	}
	return comments, nil
}

func (s *CommentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	stmt, err := s.db.PrepareContext(ctx, `
		DELETE FROM COMMENTS WHERE ID = $1
	`)
	if err != nil {
		return err
	}

	result, err := stmt.ExecContext(ctx, id)
	if v, _ := result.RowsAffected(); v == 0 {
		return repoerrs.ErrCommentNotFound
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
