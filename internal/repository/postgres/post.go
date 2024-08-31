package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/escoutdoor/social/internal/repository/repoerrs"
	"github.com/escoutdoor/social/internal/types"
	"github.com/google/uuid"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

func (s *PostRepository) Create(ctx context.Context, userID uuid.UUID, input types.CreatePostReq) (*types.Post, error) {
	stmt, err := s.db.PrepareContext(ctx, `
		INSERT INTO POSTS(CONTENT, USER_ID, PHOTO_URL) VALUES($1, $2, $3)
		RETURNING ID, CONTENT, USER_ID, PHOTO_URL, CREATED_AT, UPDATED_AT
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	args := []interface{}{input.Content, userID, input.PhotoURL}
	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		return scanPost(rows)
	}
	return nil, repoerrs.ErrPostNotFound
}

func (s *PostRepository) Update(ctx context.Context, postID uuid.UUID, input types.Post) (*types.Post, error) {
	stmt, err := s.db.PrepareContext(ctx, `
		UPDATE POSTS SET
			CONTENT = $1,
			PHOTO_URL = $2
		WHERE ID = $3
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	args := []interface{}{input.Content, input.PhotoURL, postID}
	if _, err = stmt.ExecContext(ctx, args...); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, postID)
}

func (s *PostRepository) GetByID(ctx context.Context, id uuid.UUID) (*types.Post, error) {
	stmt, err := s.db.PrepareContext(ctx, `
		SELECT 
			p.ID,
			p.CONTENT,
			p.USER_ID,
			p.PHOTO_URL,
			COUNT(l.ID) as LIKES,
			p.UPDATED_AT,
			p.CREATED_AT
		FROM POSTS p 
		LEFT JOIN POST_LIKES l ON p.ID = l.POST_ID
		WHERE p.ID = $1
		GROUP BY p.ID
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var post types.Post
	err = stmt.QueryRowContext(ctx, id).Scan(
		&post.ID,
		&post.Content,
		&post.UserID,
		&post.PhotoURL,
		&post.Likes,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repoerrs.ErrPostNotFound
		}
		return nil, err
	}
	return &post, err
}

func (s *PostRepository) GetAll(ctx context.Context) ([]types.Post, error) {
	stmt, err := s.db.PrepareContext(ctx, `
		SELECT 
			p.ID,
			p.CONTENT,
			p.USER_ID,
			p.PHOTO_URL,
			COUNT(l.ID) AS LIKES,
			p.CREATED_AT,
			p.UPDATED_AT
		FROM POSTS p
		LEFT JOIN POST_LIKES l ON p.ID = l.POST_ID
		GROUP BY p.ID
		ORDER BY p.CREATED_AT
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []types.Post
	for rows.Next() {
		var p types.Post
		err = rows.Scan(
			&p.ID,
			&p.Content,
			&p.UserID,
			&p.PhotoURL,
			&p.Likes,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func (s *PostRepository) Delete(ctx context.Context, id uuid.UUID) error {
	stmt, err := s.db.PrepareContext(ctx, `
		DELETE FROM POSTS WHERE ID = $1
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}
	if ra, _ := res.RowsAffected(); ra == 0 {
		return repoerrs.ErrPostNotFound
	}
	return nil
}

func scanPost(rows *sql.Rows) (*types.Post, error) {
	var post types.Post
	err := rows.Scan(
		&post.ID,
		&post.Content,
		&post.UserID,
		&post.PhotoURL,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	return &post, err
}
