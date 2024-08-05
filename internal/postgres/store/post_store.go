package store

import (
	"context"
	"database/sql"

	"github.com/escoutdoor/social/internal/types"
	"github.com/google/uuid"
)

type PostStore struct {
	db *sql.DB
}

func NewPostStore(db *sql.DB) *PostStore {
	return &PostStore{
		db: db,
	}
}

func (s *PostStore) Create(ctx context.Context, userID uuid.UUID, input types.CreatePostReq) (uuid.UUID, error) {
	var id uuid.UUID
	stmt, err := s.db.PrepareContext(ctx, `
		INSERT INTO POSTS(CONTENT, USER_ID, PHOTO_URL) VALUES($1, $2, $3)
		RETURNING ID
	`)
	if err != nil {
		return id, err
	}

	if err = stmt.QueryRowContext(ctx, input.Content, userID, input.PhotoURL).Scan(&id); err != nil {
		return id, err
	}
	return id, nil
}

func (s *PostStore) Update(ctx context.Context, postID uuid.UUID, input types.Post) (*types.Post, error) {
	stmt, err := s.db.PrepareContext(ctx, `
		UPDATE POSTS SET
			CONTENT = $1,
			PHOTO_URL = $2
		WHERE ID = $3
	`)
	if err != nil {
		return nil, err
	}

	if _, err = stmt.ExecContext(
		ctx,
		input.Content,
		input.PhotoURL,
		postID,
	); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, postID)
}

func (s *PostStore) GetByID(ctx context.Context, id uuid.UUID) (*types.Post, error) {
	stmt, err := s.db.PrepareContext(ctx, `
		SELECT * FROM POSTS WHERE ID = $1
	`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		return scanPost(rows)
	}
	return nil, ErrPostNotFound
}

func (s *PostStore) GetAll(ctx context.Context) ([]types.Post, error) {
	stmt, err := s.db.PrepareContext(ctx, `SELECT * FROM POSTS`)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	var posts []types.Post
	for rows.Next() {
		p, err := scanPost(rows)
		if err != nil {
			return nil, err
		}
		posts = append(posts, *p)
	}
	return posts, nil
}

func (s *PostStore) Delete(ctx context.Context, id uuid.UUID) error {
	stmt, err := s.db.PrepareContext(ctx, `
		DELETE FROM POSTS WHERE ID = $1
	`)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}
	if ra, _ := res.RowsAffected(); ra == 0 {
		return ErrPostNotFound
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
