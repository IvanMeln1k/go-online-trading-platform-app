package service

import (
	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/repository"
	"github.com/IvanMeln1k/go-online-trading-platform-app/pkg/email"
)

type Users interface {
}

type Auth interface {
}

type Service struct {
	Users
	Auth
}

func NewService(repo *repository.Repository, emailSender email.EmailSender) *Service {
	return &Service{
		Users: NewUsersService(repo.Users),
		Auth:  NewAuthService(),
	}
}
