package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

var (
	ErrUserNotFound = errors.New("error user not found")
)

type UsersRepository struct {
	db *sqlx.DB
}

func NewUsersRepository(db *sqlx.DB) *UsersRepository {
	return &UsersRepository{
		db: db,
	}
}

func (r *UsersRepository) Create(ctx context.Context, user domain.User) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (username, name, email, hash_password, role)"+
		" VALUES ($1, $2, $3, $4, $5) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.Username, user.Name, user.Email, user.Password, user.Role)
	if err := row.Scan(&id); err != nil {
		logrus.Errorf("error create user into db: %s", err)
		return 0, ErrInternal
	}

	return id, nil
}

func (r *UsersRepository) get(ctx context.Context, key string,
	value interface{}) (domain.User, error) {
	var user domain.User

	query := fmt.Sprintf(`SELECT * FROM %s WHERE %s=$1`, usersTable, key)
	err := r.db.GetContext(ctx, &user, query, value)
	if err != nil {
		logrus.Errorf("error get user from db: %s", err)
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
		names = append(names, fmt.Sprintf("%s=$%d", name, argId))
		values = append(values, value)
		argId++
	}

	if data.Username != nil {
		addProp("username", *data.Username)
	}
	if data.Name != nil {
		addProp("name", *data.Name)
	}
	if data.Email != nil {
		addProp("email", *data.Email)
	}
	if data.Password != nil {
		addProp("password", *data.Password)
	}
	if data.Role != nil {
		addProp("role", *data.Role)
	}
	if data.EmailVefiried != nil {
		addProp("email_verified", *data.EmailVefiried)
	}

	setQuery := strings.Join(names, ", ")
	values = append(values, id)
	query := fmt.Sprintf(`UPDATE %s u SET %s WHERE id=$%d RETURNING u.*`,
		usersTable, setQuery, argId)
	err := r.db.Get(&user, query, values...)
	if err != nil {
		logrus.Errorf("error update user: %s", err)
		if errors.Is(err, sql.ErrNoRows) {
			return user, ErrUserNotFound
		}
		return user, ErrInternal
	}

	return user, nil
}
