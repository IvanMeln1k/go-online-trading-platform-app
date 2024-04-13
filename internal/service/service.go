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
	VerifyEmail(ctx context.Context, email string) error
	Logout(ctx context.Context, refreshToken string) error
	LogoutAll(ctx context.Context, refreshToken string) error
	Refresh(ctx context.Context, refreshToken string) (domain.Tokens, error)
	ResendEmail(ctx context.Context, id int) error
}

type Cards interface {
	GetAll(ctx context.Context, userId int) ([]domain.Card, error)
	Get(ctx context.Context, userId int, cardId int) (domain.Card, error)
	Create(ctx context.Context, userId int, card domain.Card) (int, error)
	Delete(ctx context.Context, userId int, cardId int) error
}

type Service struct {
	Auth
	Cards
}

type Deps struct {
	Repo                  *repository.Repository
	TokenManager          tokens.TokenManagerI
	PasswordManager       password.PasswordManagerI
	EmailSender           email.EmailSender
	RefreshTTL            time.Duration
	VerificationEmailAddr string
}

func NewService(deps Deps) *Service {
	return &Service{
		Auth: NewAuthService(deps.Repo.Users, deps.Repo.Sessions, deps.TokenManager,
			deps.PasswordManager, deps.EmailSender, deps.RefreshTTL, deps.VerificationEmailAddr),
	}
}
