package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

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

func (r *UsersRepository) get(ctx context.Context, key string,
	value interface{}) (domain.User, error) {
	var user domain.User

	query := fmt.Sprintf(`SELECT * FROM %s WHERE %s = $1`, usersTable, key)
	err := r.db.Get(&user, query, value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, ErrUserNotFound
		}
		return user, ErrInternal
	}

	return user, nil
}

func (r *UsersRepository) GetById(ctx context.Context, id int) (domain.User, error) {
	return r.get(ctx, "id", id)
}

func (r *UsersRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	return r.get(ctx, "email", email)
}

func (r *UsersRepository) GetByUserName(ctx context.Context, username string) (domain.User, error) {
	return r.get(ctx, "username", username)
}

func (r *UsersRepository) Update(ctx context.Context, id int, data domain.UserUpdate) (domain.User, error) {
	var user domain.User

	var names []string
	var values []interface{}
	argId := 1

	addProp := func(name string, value interface{}) {
		names = append(names, fmt.Sprintf("%s = $%d", name, argId))
		values = append(values, value)
		argId++
	}

	if data.Email != nil {
		addProp("email", *data.Email)
	}
	if data.Name != nil {
		addProp("name", *data.Email)
	}
	if data.Password != nil {
		addProp("password", *data.Password)
	}
	if data.Username != nil {
		addProp("username", *data.Username)
	}

	setQuery := strings.Join(names, ", ")
	values = append(values, id)
	query := fmt.Sprintf(`UPDATE %s u SET %s WHERE id = $%d RETURNING u.*`,
		usersTable, setQuery, argId+1)
	err := r.db.Get(&user, query, values...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, ErrUserNotFound
		}
		return user, ErrInternal
	}

	return user, nil
}
