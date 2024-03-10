package repository

import (
	"context"

	"github.com/IvanMeln1k/go-online-trading-platform-app/domain"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable = "users"
)

type Users interface {
	Create(ctx context.Context, user domain.User) (int, error)
	GetById(ctx context.Context, id int) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	GetByUserName(ctx context.Context, username string) (domain.User, error)
	Update(ctx context.Context, id int, data domain.UserUpdate) (domain.User, error)
}

type Repository struct {
	Users
}

func NewRepository(db *sqlx.DB, rdb *redis.Client) *Repository {
	return &Repository{
		Users: NewUsersRepository(db),
	}
}
