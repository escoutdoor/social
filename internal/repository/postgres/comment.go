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
		return id, err
	}
	defer stmt.Close()

	args := []interface{}{input.Content, input.ParentCommentID, postID, userID}
	err = stmt.QueryRowContext(ctx, args...).Scan(&id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23503" {
			switch pqErr.Constraint {
			case "comments_post_id_fkey":
				return id, repoerrs.ErrPostNotFound
			case "comments_parent_comment_id_fkey":
				return id, repoerrs.ErrCommentNotFound
			case "comments_user_id_fkey":
				return id, repoerrs.ErrUserNotFound
			}
		}
		return id, err
	}
	return id, nil
}

func (s *CommentRepository) GetByID(ctx context.Context, id uuid.UUID) (*types.Comment, error) {
	stmt, err := s.db.PrepareContext(ctx, `
		SELECT 
			c.ID,
			c.CONTENT,
			c.USER_ID,
			c.POST_ID,
			c.PARENT_COMMENT_ID,
			COUNT(l.ID) as LIKES,
			c.UPDATED_AT,
			c.CREATED_AT
		FROM COMMENTS c
		LEFT JOIN COMMENT_LIKES l ON c.ID = l.COMMENT_ID
		WHERE c.ID = $1
		GROUP BY c.ID
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var comment types.Comment
	err = stmt.QueryRowContext(ctx, id).Scan(
		&comment.ID,
		&comment.Content,
		&comment.UserID,
		&comment.PostID,
		&comment.ParentCommentID,
		&comment.Likes,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repoerrs.ErrCommentNotFound
		}
		return nil, err
	}
	return &comment, err
}

func (s *CommentRepository) GetAll(ctx context.Context, postID uuid.UUID) ([]types.Comment, error) {
	stmt, err := s.db.PrepareContext(ctx, `
		SELECT 
			cm.ID,
			cm.CONTENT,
			cm.USER_ID,
			cm.POST_ID,
			cm.PARENT_COMMENT_ID,
			COUNT(cl.ID) AS LIKES,
			cm.CREATED_AT,
			cm.UPDATED_AT
		FROM COMMENTS cm
		LEFT JOIN COMMENT_LIKES cl ON cm.ID = cl.COMMENT_ID
		WHERE cm.POST_ID = $1
		GROUP BY cm.ID
		ORDER BY LIKES, CREATED_AT
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		commentsMap = make(map[uuid.UUID]*types.Comment)
		comments    []types.Comment
	)
	for rows.Next() {
		var (
			c    types.Comment
			pcid uuid.NullUUID
		)
		if err := rows.Scan(
			&c.ID,
			&c.Content,
			&c.UserID,
			&c.PostID,
			&pcid,
			&c.Likes,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if pcid.Valid {
			c.ParentCommentID = &pcid.UUID
		}

		commentsMap[c.ID] = &c
		if c.ParentCommentID == nil {
			comments = append(comments, c)
		}
	}

	for i := range comments {
		comments[i].Replies = getReplies(comments[i].ID, commentsMap)
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
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, id)
	if v, _ := result.RowsAffected(); v == 0 {
		return repoerrs.ErrCommentNotFound
	}
	return nil
}

func getReplies(id uuid.UUID, commentsMap map[uuid.UUID]*types.Comment) []types.Comment {
	var replies []types.Comment
	for _, cm := range commentsMap {
		if cm.ParentCommentID != nil && *cm.ParentCommentID == id {
			cm.Replies = getReplies(cm.ID, commentsMap)
			replies = append(replies, *cm)
		}
	}
	return replies
}
