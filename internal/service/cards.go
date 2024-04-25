package service

import (
	"context"
	"errors"

	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/domain"
	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/repository"
	"github.com/sirupsen/logrus"
)

type CardsService struct {
	cardsRepo repository.Cards
	usersRepo repository.Users
}

var (
	ErrCardNotFound = errors.New("card not found")
)

func NewCardsService(cardsRepo repository.Cards, usersRepo repository.Users) *CardsService {
	return &CardsService{cardsRepo, usersRepo}
}

func (s *CardsService) GetAll(ctx context.Context, userId int) ([]domain.Card, error) {
	_, err := s.usersRepo.GetById(ctx, userId)
	if err != nil {
		logrus.Errorf("Service cards error GetAll: %s", err)
		if errors.Is(repository.ErrUserNotFound, err) {
			return nil, ErrUserNotFound
		}
		return nil, ErrInternal
	}
	cards, err := s.cardsRepo.GetAll(ctx, userId)
	if err != nil {
		logrus.Errorf("Service GetAllCards calling repository error: %s", err)
		return nil, ErrInternal
	}

	//Доделать конфиденциальный вывод карт
	return cards, nil
}

func (s *CardsService) Get(ctx context.Context, userId int, cardId int) (domain.Card, error) {
	card, err := s.cardsRepo.Get(ctx, cardId)
	if err != nil {
		logrus.Errorf("Service GetCard calling repository error: %s", err)
		if errors.Is(repository.ErrCardNotFound, err) {
			return card, ErrCardNotFound
		}
		return card, ErrInternal
	}
	if card.UserId != userId {
		return card, ErrCardNotFound
	}
	return card, nil
}

func (s *CardsService) Create(ctx context.Context, userId int, card domain.Card) (int, error) {
	_, err := s.usersRepo.GetById(ctx, userId)
	if err != nil {
		logrus.Errorf("Service CreateCard getUser from repository when creating card error: %s", err)
		if errors.Is(repository.ErrUserNotFound, err) {
			return 0, ErrUserNotFound
		}
		return 0, ErrInternal
	}
	id, err := s.cardsRepo.Create(ctx, card)
	if err != nil {
		logrus.Errorf("Service CreateCard are broken when calling repository error: %s", err)
		return 0, ErrInternal
	}
	return id, nil
}

func (s *CardsService) Delete(ctx context.Context, userId int, cardId int) error {
	card, err := s.cardsRepo.Get(ctx, userId)
	if err != nil {
		logrus.Errorf("Service DeleteCard error when getting card: %s", err)
		if errors.Is(repository.ErrCardNotFound, err) {
			return ErrCardNotFound
		}
		return ErrInternal
	}
	if card.UserId != userId {
		logrus.Error("User id of card is not got userId")
		return ErrCardNotFound
	}
	err = s.cardsRepo.Delete(ctx, cardId)
	if err != nil {
		logrus.Errorf("Service DeleteCard error when deleting card from repo: %s", err)
		return ErrInternal
	}
	return nil
}