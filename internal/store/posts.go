package store

import (
	"context"
	"database/sql"
	"errors"
)

type Post struct {
	Id           int64     `json:"id"`
	Content      string    `json:""`
	Title        string    `json:"title"`
	UserId       int64     `json:"user_id"`
	CreatedAt    string    `json:"createdAt"`
	UpdatedAt    string    `json:"updatedAt"`
	Comments     []Comment `json:"comments,omitempty"`
	CommentCount int64     `json:"commentCount,omitempty"`
}

type PostStore struct {
	db           *sql.DB
	commentStore *CommentStore
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

	if s.commentStore != nil {
		comments, err := s.commentStore.GetByPostId(ctx, post.Id)
		if err != nil {
			return nil, err
		}
		post.Comments = comments
	}

	return &post, nil
}

func (s *PostStore) GetAll(ctx context.Context) ([]Post, error) {
	query := `
		SELECT 
			p.id,
			p.userId,
			p.title,
			p.content,
			p.createdAt,
			p.updatedAt,
			COUNT(c.id) AS commentCount
		FROM posts p
		LEFT JOIN comments c ON c.postId = p.id
		GROUP BY p.id
		ORDER BY p.createdAt DESC
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(
			&post.Id,
			&post.UserId,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.CommentCount,
		); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
