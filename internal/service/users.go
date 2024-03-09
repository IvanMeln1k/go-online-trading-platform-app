package service

import "github.com/IvanMeln1k/go-online-trading-platform-app/internal/repository"

type UsersService struct {
	repo repository.Users
}

func NewUsersService(repo repository.Users) *UsersService {
	return &UsersService{
		repo: repo,
	}
}
