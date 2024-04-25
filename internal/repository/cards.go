package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type CardsRepository struct {
	db *sqlx.DB
}

var (
	ErrCardNotFound = errors.New("card not found")
)

func (r *CardsRepository) Create(ctx context.Context, card domain.Card) (int, error) {
	var id int
	row := r.db.QueryRow(`INSERT INTO cards (number, data, cvv, users_id) VALUES ($1, $2, $3, $4) RETURNING id`, card.Number, card.Data, card.Cvv, card.UserId)
	if err := row.Scan(&id); err != nil {
		logrus.Errorf("Creation error of card: %s", err)
		return 0, ErrInternal
	}
	return id, nil
}

func (r *CardsRepository) Get(ctx context.Context, cardId int) (domain.Card, error) {
	var Card domain.Card
	row := r.db.QueryRow(`SELECT * FROM cards WHERE id = $1`, cardId)
	if err := row.Scan(&Card); err != nil {
		logrus.Errorf("Error get card from postgresql: %s", err)
		if errors.Is(sql.ErrNoRows, err) {
			return Card, ErrCardNotFound
		}
		return Card, ErrInternal
	}
	return Card, nil

}

func (r *CardsRepository) GetAll(ctx context.Context, userId int) ([]domain.Card, error) {
	var cards []domain.Card
	err := r.db.Select(cards, "SELECT * FROM cards WHERE user_id = $1", userId)
	if err != nil {
		logrus.Errorf("Error get cards from postgresql: %s", err)
		return cards, ErrInternal
	}
	return cards, nil
}

func (r *CardsRepository) Delete(ctx context.Context, cardId int) error {
	_, err := r.db.Exec("DELETE FROM cards WHERE id = $1", cardId)
	if err != nil {
		logrus.Errorf("Error deleting card from posgresql: %s", err)
		return ErrInternal
	}
	return nil
}