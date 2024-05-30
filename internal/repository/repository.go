package repository

import (
	"context"
	"errors"

	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

const (
	usersTable = "users"
)

var (
	ErrInternal = errors.New("internal error")
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type Users interface {
	// Creates a user in DB
	Create(ctx context.Context, user domain.User) (int, error)
	// Returns the user from DB by id
	GetById(ctx context.Context, id int) (domain.User, error)
	// Returns the user from DB by email
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	// Returns the user from DB by username
	GetByUserName(ctx context.Context, username string) (domain.User, error)
	// Updates the user by id
	Update(ctx context.Context, id int, data domain.UserUpdate) (domain.User, error)
}

type Cards interface {
	GetAll(ctx context.Context, userId int) ([]domain.Card, error)
	Get(ctx context.Context, cardId int) (domain.Card, error)
	Create(ctx context.Context, card domain.Card) (int, error)
	Delete(ctx context.Context, cardId int) error
}

type Sessions interface {
	// Creates a session record in DB
	Create(ctx context.Context, session domain.Session) error
	// Returns the session from DB by refreshToken (refreshToken = sessionId)
	// Returns error if the session is invalid or expired (doesn't exist)
	Get(ctx context.Context, refreshToken string) (domain.Session, error)
	// Deletes user's session record from DB by refreshToken (refreshToken = sessionId)
	Delete(ctx context.Context, userId int, refreshToken string) error
	// Returns count of user's sessions in DB by userId
	GetCnt(ctx context.Context, userId int) (int, error)
	// Returns all user's sessions from DB by userId
	GetAll(ctx context.Context, userId int) ([]domain.Session, error)
	// Deletes all user's sessions from DB by userId
	DeleteAll(ctx context.Context, userId int) error
}
type Products interface {
	GetMyAll(ctx context.Context, userId int) ([]domain.Product, error)
	Get(ctx context.Context, productId int) (domain.Product, error)
	Create(ctx context.Context, product domain.Product) (int, error)
	Delete(ctx context.Context, productId int) error
	GetAll(ctx context.Context, filter domain.Filter) ([]domain.Product, error)
}

type Repository struct {
	Users
	Sessions
	Cards
	Products
}

func NewRepository(db *sqlx.DB, rdb *redis.Client) *Repository {
	return &Repository{
		Users:    NewUsersRepository(db),
		Sessions: NewSessionsRepository(rdb),
		Cards:    NewCardsRepository(db),
		Products: NewProductsRepository(db),
	}
}
