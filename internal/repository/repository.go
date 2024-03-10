package repository

import (
	"context"
	"errors"

	"github.com/IvanMeln1k/go-online-trading-platform-app/domain"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable = "users"
)

var (
	ErrUserNotFound            = errors.New("error user not found")
	ErrInternal                = errors.New("internal error")
	ErrSessionExpiredOrInvalid = errors.New("error session expired or invalid")
)

type Users interface {
	Create(ctx context.Context, user domain.User) (int, error)
	GetById(ctx context.Context, id int) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	GetByUserName(ctx context.Context, username string) (domain.User, error)
	Update(ctx context.Context, id int, data domain.UserUpdate) (domain.User, error)
}

type Sessions interface {
	Create(ctx context.Context, session domain.Session) error
	Get(ctx context.Context, refreshToken string) (domain.Session, error)
	Delete(ctx context.Context, userId int, refreshToken string) error
	GetCnt(ctx context.Context, userId int) (int, error)
	GetAll(ctx context.Context, userId int) ([]domain.Session, error)
	DeleteAll(ctx context.Context, userId int) error
}

type Repository struct {
	Users
	Sessions
}

func NewRepository(db *sqlx.DB, rdb *redis.Client) *Repository {
	return &Repository{
		Users:    NewUsersRepository(db),
		Sessions: NewSessionsRepository(rdb),
	}
}
