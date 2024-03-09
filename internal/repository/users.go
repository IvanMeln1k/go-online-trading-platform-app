package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/IvanMeln1k/go-online-trading-platform-app/domain"
	"github.com/jmoiron/sqlx"
)

type UsersRepository struct {
	db *sqlx.DB
}

func NewUsersRepository(db *sqlx.DB) *UsersRepository {
	return &UsersRepository{
		db: db,
	}
}

var (
	ErrUserNotFound = errors.New("error user not found")
	ErrInternal     = errors.New("internal error")
)

func (r *UsersRepository) Create(ctx context.Context, user domain.User) (int, error) {
	var id int

	query := fmt.Sprintf(`INSERT INTO %s (email, username, name, hash_password)
	VALUES ($1, $2, $3, $4) RETURNING id`, usersTable)
	res := r.db.QueryRow(query, user.Email, user.Username, user.Name, user.Password)
	if err := res.Scan(&id); err != nil {
		return 0, ErrInternal
	}

	return id, nil
}

func (r *UsersRepository) GetById(ctx context.Context, id int) (domain.User, error) {
	var user domain.User

	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1`, usersTable)
	err := r.db.Get(&user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, ErrUserNotFound
		}
		return user, ErrInternal
	}

	return user, nil
}
