package repository

import (
	"context"

	"github.com/IvanMeln1k/go-online-trading-platform-app/domain"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable = "users"
)

type Users interface {
	Create(ctx context.Context, user domain.User) (int, error)
	GetById(ctx context.Context, id int) (domain.User, error)
}

type Repository struct {
	Users
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Users: NewUsersRepository(db),
	}
}
