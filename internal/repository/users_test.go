package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestUsersPostgres_Create(t *testing.T) {
	sqlmock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxmock := sqlx.NewDb(sqlmock, "postgres")

	usersRepo := NewUsersRepository(sqlxmock)

	type args struct {
		ctx  context.Context
		user domain.User
	}

	type mockBehavior func(args args, id int)

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		want         int
		wantErr      bool
		err          error
	}{
		{
			name: "ok",
			args: args{
				ctx: context.Background(),
				user: domain.User{
					Username: "IvanMelnik",
					Name:     "Ivan",
					Email:    "IvanMelnikovF@gmail.com",
					Password: "pass",
				},
			},
			mockBehavior: func(args args, id int) {
				rows := mock.NewRows([]string{"id"}).AddRow(49)
				mock.ExpectQuery(`INSERT INTO users \(username, name, email, hash_password\)`+
					` VALUES \(\$1, \$2, \$3, \$4\) RETURNING id`).
					WithArgs(args.user.Username, args.user.Name, args.user.Email, args.user.Password).
					WillReturnRows(rows)
			},
			want:    49,
			wantErr: false,
			err:     nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.args, test.want)

			got, err := usersRepo.Create(test.args.ctx, test.args.user)

			if test.wantErr {
				assert.ErrorIs(t, err, test.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUsersPostgres_GetById(t *testing.T) {
	sqlmock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	sqlxmock := sqlx.NewDb(sqlmock, "postgres")

	usersRepository := NewUsersRepository(sqlxmock)

	type args struct {
		id int
	}

	type mockBehavior func(args args, user domain.User)

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		want         domain.User
		wantErr      error
	}{
		{
			name: "ok",
			args: args{
				id: 1,
			},
			mockBehavior: func(args args, user domain.User) {
				rows := mock.NewRows([]string{"id", "username", "name", "email", "hash_password"}).
					AddRow(user.Id, user.Username, user.Name, user.Email, user.Password)
				mock.ExpectQuery("SELECT (.+) FROM users WHERE id=(.+)").WithArgs(args.id).
					WillReturnRows(rows)
			},
			want: domain.User{
				Id:       1,
				Username: "username",
				Name:     "name",
				Email:    "email",
				Password: "password",
			},
			wantErr: nil,
		},
		{
			name: "not found",
			args: args{
				id: 1,
			},
			mockBehavior: func(args args, user domain.User) {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE id=(.+)").WithArgs(args.id).
					WillReturnError(sql.ErrNoRows)
			},
			want:    domain.User{},
			wantErr: ErrUserNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.args, test.want)

			got, err := usersRepository.GetById(context.Background(), test.args.id)

			assert.ErrorIs(t, test.wantErr, err)
			assert.Equal(t, test.want, got)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUsersPostgres_GetByEmail(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub db connection", err)
	}
	sqlxDB := sqlx.NewDb(sqlDB, "postgres")

	UsersRepository := NewUsersRepository(sqlxDB)

	type args struct {
		email string
	}

	type mockBehavior func(args args, user domain.User)

	tests := []struct {
		name         string
		args         args
		want         domain.User
		mockBehavior mockBehavior
		wantErr      error
	}{
		{
			name: "ok",
			args: args{
				email: "email",
			},
			want: domain.User{
				Id:       1,
				Username: "username",
				Name:     "name",
				Email:    "email",
				Password: "password",
			},
			mockBehavior: func(args args, user domain.User) {
				rows := mock.NewRows([]string{"id", "username", "name", "email", "hash_password"}).
					AddRow(user.Id, user.Username, user.Name, user.Email, user.Password)
				mock.ExpectQuery("SELECT (.+) FROM users WHERE email=(.+)").WithArgs(args.email).
					WillReturnRows(rows)
			},
			wantErr: nil,
		},
		{
			name: "not found",
			args: args{
				email: "email",
			},
			want: domain.User{},
			mockBehavior: func(args args, user domain.User) {
				rows := mock.NewRows([]string{"id", "username", "name", "email", "hash_password"})
				mock.ExpectQuery("SELECT (.+) FROM users WHERE email=(.+)").WithArgs(args.email).
					WillReturnRows(rows)
			},
			wantErr: ErrUserNotFound,
		},
		{
			name: "internal error",
			args: args{
				email: "email",
			},
			want: domain.User{},
			mockBehavior: func(args args, user domain.User) {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE email=(.+)").WithArgs(args.email).
					WillReturnError(errors.New("some sql error"))
			},
			wantErr: ErrInternal,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.args, test.want)

			got, err := UsersRepository.GetByEmail(context.Background(), test.args.email)

			assert.ErrorIs(t, test.wantErr, err)
			if test.wantErr == nil {
				assert.Equal(t, test.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUsersPostres_GetByUserName(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not excepted when opening a stub db connection", err)
	}
	sqlxDB := sqlx.NewDb(sqlDB, "postgres")

	usersRepository := NewUsersRepository(sqlxDB)

	type args struct {
		username string
	}

	type mockBehavior func(args args, user domain.User)

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		want         domain.User
		wantErr      error
	}{
		{
			name: "ok",
			args: args{
				username: "username",
			},
			mockBehavior: func(args args, user domain.User) {
				rows := mock.NewRows([]string{"id", "username", "name", "email", "hash_password"}).
					AddRow(user.Id, user.Username, user.Name, user.Email, user.Password)
				mock.ExpectQuery("SELECT (.+) FROM users WHERE username=(.+)").WithArgs(args.username).
					WillReturnRows(rows)
			},
			want: domain.User{
				Id:       1,
				Username: "username",
				Name:     "name",
				Email:    "email",
				Password: "password",
			},
			wantErr: nil,
		},
		{
			name: "not found",
			args: args{
				username: "username",
			},
			mockBehavior: func(args args, user domain.User) {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE username=(.+)").
					WithArgs(args.username).WillReturnError(sql.ErrNoRows)
			},
			want:    domain.User{},
			wantErr: ErrUserNotFound,
		},
		{
			name: "internal error",
			args: args{
				username: "username",
			},
			mockBehavior: func(args args, user domain.User) {
				mock.ExpectQuery("SELECT (.+) FROM users WHERE username=(.+)").
					WithArgs(args.username).WillReturnError(errors.New("some error"))
			},
			want:    domain.User{},
			wantErr: ErrInternal,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.args, test.want)

			got, err := usersRepository.GetByUserName(context.Background(), test.args.username)

			assert.ErrorIs(t, test.wantErr, err)
			if test.wantErr == nil {
				assert.Equal(t, test.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUsersPostgres_Update(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub db connection", err)
	}
	sqlxDB := sqlx.NewDb(sqlDB, "postgres")

	UsersRepository := NewUsersRepository(sqlxDB)

	type args struct {
		id   int
		user domain.UserUpdate
	}

	type mockBehavior func(args args, user domain.User)

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		want         domain.User
		wantErr      error
	}{
		{
			name: "ok/full",
			args: args{
				id: 1,
				user: domain.UserUpdate{
					Username:      stringPointer("username"),
					Name:          stringPointer("name"),
					Email:         stringPointer("email"),
					Password:      stringPointer("password"),
					EmailVefiried: boolPointer(true),
				},
			},
			mockBehavior: func(args args, user domain.User) {
				rows := mock.NewRows([]string{"id", "username", "name", "email", "hash_password",
					"email_verified"}).
					AddRow(user.Id, user.Username, user.Name, user.Email, user.Password,
						user.EmailVerified)
				mock.ExpectQuery("UPDATE users u SET (.+) WHERE id=(.+) RETURNING (.+)").
					WithArgs(*args.user.Username, *args.user.Name, *args.user.Email,
						*args.user.Password, *args.user.EmailVefiried, args.id).
					WillReturnRows(rows)
			},
			want: domain.User{
				Id:            1,
				Username:      "username",
				Name:          "name",
				Email:         "email",
				Password:      "password",
				EmailVerified: true,
			},
			wantErr: nil,
		},
		{
			name: "ok/email",
			args: args{
				user: domain.UserUpdate{
					Email: stringPointer("email"),
				},
			},
			mockBehavior: func(args args, user domain.User) {
				rows := mock.NewRows([]string{"id", "username", "name", "email", "hash_password",
					"email_verified"}).
					AddRow(user.Id, user.Username, user.Name, user.Email, user.Password,
						user.EmailVerified)
				mock.ExpectQuery("UPDATE users u SET email = (.+) WHERE id=(.+) RETURNING (.+)").
					WithArgs(*args.user.Email, args.id).WillReturnRows(rows)
			},
			want: domain.User{
				Id:            1,
				Username:      "username",
				Name:          "name",
				Email:         "email",
				Password:      "password",
				EmailVerified: true,
			},
			wantErr: nil,
		},
		{
			name: "not found",
			args: args{
				user: domain.UserUpdate{
					Email: stringPointer("email"),
				},
			},
			mockBehavior: func(args args, user domain.User) {
				mock.ExpectQuery("UPDATE users u SET email = (.+) WHERE id=(.+) RETURNING (.+)").
					WithArgs(*args.user.Email, args.id).WillReturnError(sql.ErrNoRows)
			},
			want:    domain.User{},
			wantErr: ErrUserNotFound,
		},
		{
			name: "internal error",
			args: args{
				user: domain.UserUpdate{
					Email: stringPointer("email"),
				},
				id: 1,
			},
			mockBehavior: func(args args, user domain.User) {
				mock.ExpectQuery("UPDATE users u SET email = (.+) WHERE id=(.+) RETURNING (.+)").
					WithArgs(*args.user.Email, args.id).WillReturnError(errors.New("some error"))
			},
			want:    domain.User{},
			wantErr: ErrInternal,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.args, test.want)

			got, err := UsersRepository.Update(context.Background(), test.args.id, test.args.user)

			if test.wantErr != nil {
				assert.ErrorIs(t, test.wantErr, err)
			} else {
				assert.Equal(t, test.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func stringPointer(s string) *string {
	return &s
}

func boolPointer(b bool) *bool {
	return &b
}
