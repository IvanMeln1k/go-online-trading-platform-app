package service

import "github.com/IvanMeln1k/go-online-trading-platform-app/internal/repository"

type Service struct {
}

func NewService(repo *repository.Repository) *Service {
	return &Service{}
}
