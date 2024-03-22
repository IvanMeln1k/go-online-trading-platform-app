package service

import (
	"errors"

	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/repository"
)

type UsersService struct {
	usersRepo repository.Users
	cardsRepo repository.Cards
}

func NewUsersService(usersRepo repository.Users, cardsRepo repository.Cards) *UsersService {
	return &UsersService{
		usersRepo: usersRepo,
		cardsRepo: cardsRepo,
	}
}

var (
	ErrCardNotFound = errors.New("error card not found")
)
