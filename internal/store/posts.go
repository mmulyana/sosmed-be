package store

import (
	"context"
	"database/sql"
	"errors"
)

type Post struct {
	Id        int64  `json:"id"`
	Content   string `json:""`
	Title     string `json:"title"`
	UserId    int64  `json:"user_id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (content, title, userId)
		VALUES (?, ?, ?)
	`

	res, err := s.db.ExecContext(ctx, query, post.Content, post.Title, post.UserId)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	post.Id = id

	err = s.db.QueryRowContext(ctx,
		"SELECT createdAt, updatedAt FROM posts WHERE id = ?",
		post.Id,
	).Scan(&post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `
		SELECT id, userId, title, content, createdAt, updatedAt
		FROM posts
		WHERE id = ?
	`

	var post Post
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.Id,
		&post.UserId,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}
