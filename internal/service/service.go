package service

import (
	"context"
	"errors"
	"time"

	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/domain"
	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/repository"
	"github.com/IvanMeln1k/go-online-trading-platform-app/pkg/email"
	"github.com/IvanMeln1k/go-online-trading-platform-app/pkg/password"
	"github.com/IvanMeln1k/go-online-trading-platform-app/pkg/tokens"
)

var (
	ErrInternal = errors.New("internal error")
)

type Auth interface {
	SignUp(ctx context.Context, user domain.User) (int, error)
	SignIn(ctx context.Context, email string, password string) (domain.Tokens, error)
	// Logout(ctx context.Context, refreshToken string) error
	// LogoutAll(ctx context.Context, refreshToken string) error
}

type Service struct {
	Auth
}

type Deps struct {
	Repo            *repository.Repository
	TokenManager    tokens.TokenManagerI
	PasswordManager password.PasswordManagerI
	EmailSender     email.EmailSender
	RefreshTTL      time.Duration
}

func NewService(deps Deps) *Service {
	return &Service{
		Auth: NewAuthService(deps.Repo.Users, deps.Repo.Sessions, deps.TokenManager,
			deps.PasswordManager, deps.EmailSender, deps.RefreshTTL),
	}
}
