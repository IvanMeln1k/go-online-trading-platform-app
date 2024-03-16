package service

import (
	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/repository"
	"github.com/IvanMeln1k/go-online-trading-platform-app/pkg/email"
)

type Auth interface {
}

type Service struct {
	Auth
}

func NewService(repo *repository.Repository, emailSender email.EmailSender) *Service {
	return &Service{
		Auth: NewAuthService(repo.Users, repo.Sessions),
	}
}
