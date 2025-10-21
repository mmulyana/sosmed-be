package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("record not found")
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (*Post, error)
		GetAll(context.Context) ([]Post, error)
	}

	Users interface {
		Create(context.Context, *User) error
	}

	Comments interface {
		Create(context.Context, *Comment) error
		GetByPostId(context.Context, int64) ([]Comment, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	comments := &CommentStore{db}

	return Storage{
		Posts:    &PostStore{db: db, commentStore: comments},
		Users:    &UserStore{db},
		Comments: comments,
	}
}
