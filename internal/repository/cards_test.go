package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCards_Create(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an errror '%s' was not expected when opening a stub db connection", err)
	}
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	cardsRepo := NewCardsRepository(sqlxDB)

	type args struct {
		card domain.Card
	}

	type mockBehavior func(args args, id int)

	tests := []struct {
		name         string
		args         args
		want         int
		mockBehavior mockBehavior
		wantErr      bool
		err          error
	}{
		{
			name: "ok",
			args: args{
				card: domain.Card{
					Number: "1234 5678 9012 3456",
					Data:   "13/11",
					Cvv:    "123",
					UserId: 1,
				},
			},
			want: 1,
			mockBehavior: func(args args, id int) {
				rows := mock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery(`INSERT INTO cards`).WithArgs(args.card.Number,
					args.card.Data, args.card.Cvv, args.card.UserId).WillReturnRows(rows)
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "sql error",
			args: args{
				card: domain.Card{
					Number: "",
					Data:   "",
					Cvv:    "",
					UserId: 1,
				},
			},
			want: 0,
			mockBehavior: func(args args, id int) {
				mock.ExpectQuery("INSERT INTO cards").WithArgs(args.card.Number,
					args.card.Data, args.card.Cvv, args.card.UserId).
					WillReturnError(errors.New("sql error"))
			},
			wantErr: true,
			err:     ErrInternal,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.args, test.want)

			got, err := cardsRepo.Create(context.Background(), test.args.card)

			if test.wantErr {
				assert.ErrorIs(t, test.err, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCards_Get(t *testing.T) {
	sqlMock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub db connection", err)
	}
	defer sqlMock.Close()

	sqlxDB := sqlx.NewDb(sqlMock, "sqlmock")
	cardsRepo := NewCardsRepository(sqlxDB)

	type args struct {
		cardId int
	}

	type mockBehavior func(args args, card domain.Card)

	tests := []struct {
		name         string
		args         args
		want         domain.Card
		mockBehavior mockBehavior
		wantErr      bool
		err          error
	}{
		{
			name: "ok",
			args: args{
				cardId: 1,
			},
			want: domain.Card{
				Id:     1,
				Number: "1234 1234 1234 1234",
				Data:   "13/11",
				Cvv:    "123",
				UserId: 1,
			},
			mockBehavior: func(args args, card domain.Card) {
				rows := mock.NewRows([]string{"id", "number", "data", "cvv", "user_id"}).
					AddRow(card.Id, card.Number, card.Data, card.Cvv, card.UserId)
				mock.ExpectQuery("SELECT (.+) FROM cards WHERE id = (.+)").
					WithArgs(args.cardId).
					WillReturnRows(rows)
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "card not found",
			args: args{
				cardId: 12,
			},
			want: domain.Card{},
			mockBehavior: func(args args, card domain.Card) {
				mock.ExpectQuery("SELECT (.+) FROM cards WHERE id = (.+)").
					WithArgs(args.cardId).WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
			err:     ErrCardNotFound,
		},
		{
			name: "some sql error",
			args: args{
				cardId: 123,
			},
			want: domain.Card{},
			mockBehavior: func(args args, card domain.Card) {
				mock.ExpectQuery("SELECT (.+) FROM cards WHERE id = (.+)").
					WithArgs(args.cardId).WillReturnError(errors.New("some sql error"))
			},
			wantErr: true,
			err:     ErrInternal,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.args, test.want)

			got, err := cardsRepo.Get(context.Background(), test.args.cardId)

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

func TestCards_GetAll(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub db connection", err)
	}

	sqlxDb := sqlx.NewDb(mockDb, "sqlmock")
	cardsRepo := NewCardsRepository(sqlxDb)

	type args struct {
		userId int
	}

	type mockBehavior func(args args, cards []domain.Card)

	tests := []struct {
		name         string
		args         args
		want         []domain.Card
		mockBehavior mockBehavior
		wantErr      bool
		err          error
	}{
		{
			name: "ok",
			args: args{
				userId: 1,
			},
			want: []domain.Card{
				{
					Id:     1,
					Number: "1234 1234 1234 1234",
					Data:   "13/11",
					Cvv:    "123",
					UserId: 1,
				},
				{
					Id:     22,
					Number: "4321 4321 4321 4321",
					Data:   "11/13",
					Cvv:    "321",
					UserId: 1,
				},
			},
			mockBehavior: func(args args, cards []domain.Card) {
				rows := mock.NewRows([]string{"id", "number", "data", "cvv", "user_id"})
				for _, card := range cards {
					rows.AddRow(card.Id, card.Number, card.Data, card.Cvv, card.UserId)
				}
				mock.ExpectQuery("SELECT (.+) FROM cards WHERE user_id = (.+)").
					WithArgs(args.userId).WillReturnRows(rows)
			},
		},
		{
			name: "no cards",
			args: args{
				userId: 1,
			},
			want: []domain.Card{},
			mockBehavior: func(args args, cards []domain.Card) {
				rows := mock.NewRows([]string{"id", "number", "data", "cvv", "user_id"})
				mock.ExpectQuery("SELECT (.+) FROM cards WHERE user_id = (.+)").
					WithArgs(args.userId).WillReturnRows(rows)
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "some sql error",
			args: args{
				userId: 1,
			},
			want: nil,
			mockBehavior: func(args args, cards []domain.Card) {
				mock.ExpectQuery("SELECT (.+) FROM cards WHERE user_id = (.+)").
					WithArgs(args.userId).WillReturnError(errors.New("some sql error"))
			},
			wantErr: true,
			err:     ErrInternal,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.args, test.want)

			got, err := cardsRepo.GetAll(context.Background(), test.args.userId)

			if test.wantErr {
				assert.ErrorIs(t, err, test.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, got, test.want)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCards_Delete(t *testing.T) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub db connection", err)
	}

	sqlxDb := sqlx.NewDb(mockDb, "sqlmock")
	cardsRepo := NewCardsRepository(sqlxDb)

	type args struct {
		cardId int
	}

	type mockBehavior func(args args)

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      bool
		err          error
	}{
		{
			name: "ok",
			args: args{
				cardId: 1,
			},
			mockBehavior: func(args args) {
				mock.ExpectExec("DELETE FROM cards WHERE id = (.+)").WithArgs(args.cardId)
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "some sql error",
			args: args{
				cardId: 2,
			},
			mockBehavior: func(args args) {
				mock.ExpectExec("DELETE FROM cards WHERE id = (.+)").
					WillReturnError(errors.New("some sql error"))
			},
			wantErr: true,
			err:     ErrInternal,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.args)

			err := cardsRepo.Delete(context.Background(), test.args.cardId)

			if test.wantErr {
				assert.ErrorIs(t, err, test.err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
