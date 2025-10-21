package store

import (
	"context"
	"database/sql"
	"time"
)

type Comment struct {
	Id        int64      `json:"id"`
	PostId    int64      `json:"postId"`
	UserId    int64      `json:"userId"`
	Content   string     `json:"content"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type CommentStore struct {
	db *sql.DB
}

func (s *CommentStore) Create(ctx context.Context, comment *Comment) error {
	query := `
		INSERT INTO comments (content, userId, postId)
		VALUES (?, ?, ?)
	`

	res, err := s.db.ExecContext(ctx, query, comment.Content, comment.UserId, comment.PostId)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	comment.Id = id

	var createdAt, updatedAt sql.NullTime

	err = s.db.QueryRowContext(ctx,
		"SELECT createdAt, updatedAt FROM comments WHERE id = ?",
		comment.Id,
	).Scan(&createdAt, &updatedAt)
	if err != nil {
		return err
	}

	if createdAt.Valid {
		comment.CreatedAt = &createdAt.Time
	}
	if updatedAt.Valid {
		comment.UpdatedAt = &updatedAt.Time
	}

	return nil
}

func (s *CommentStore) GetByPostId(ctx context.Context, postId int64) ([]Comment, error) {
	query := `
		SELECT id, postId, userId, content, createdAt, updatedAt
		FROM comments
		WHERE postId = ?
		ORDER BY createdAt ASC
	`

	rows, err := s.db.QueryContext(ctx, query, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment

	for rows.Next() {
		var c Comment
		var createdAt, updatedAt sql.NullTime

		if err := rows.Scan(
			&c.Id,
			&c.PostId,
			&c.UserId,
			&c.Content,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		if createdAt.Valid {
			c.CreatedAt = &createdAt.Time
		}
		if updatedAt.Valid {
			c.UpdatedAt = &updatedAt.Time
		}

		comments = append(comments, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
