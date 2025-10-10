package store

import (
	"context"
	"database/sql"
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
		VALUES ($1, $2, $3) RETURNING id, createdAt, updatedAt
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserId,
	).Scan(
		&post.Id,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil

}
