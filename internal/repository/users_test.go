package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/domain"
	"github.com/jmoiron/sqlx"
)

func TestUsersPostgres_Create(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	mock.ExpectQuery("INSERT INTO users").
		WithArgs("username", "name", "email", "password").WillReturnRows(mock.NewRows([]string{"id"}).AddRow(1))

	authRepository := NewUsersRepository(sqlxDB)

	id, err := authRepository.Create(context.Background(), domain.User{
		Username: "username",
		Name:     "name",
		Email:    "email",
		Password: "password",
	})

	if err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if id != 1 {
		t.Errorf("there were unfilfilled result, want 1, got %d", id)
	}
}
