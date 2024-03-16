package service

import "github.com/IvanMeln1k/go-online-trading-platform-app/internal/repository"

type AuthService struct {
	usersRepo    repository.Users
	sessionsRepo repository.Sessions
}

func NewAuthService(usersRepo repository.Users, sessionsRepo repository.Sessions) *AuthService {
	return &AuthService{
		usersRepo:    usersRepo,
		sessionsRepo: sessionsRepo,
	}
}
