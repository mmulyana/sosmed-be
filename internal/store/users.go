package store

import (
	"context"
	"database/sql"
)

type User struct {
	Id        int64  `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (username, password)
		VALUES ($1, $2)
		RETURNING id, createdAt, updatedAt
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Password,
	).Scan(
		&user.Id,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}
