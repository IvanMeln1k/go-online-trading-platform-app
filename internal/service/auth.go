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
	"github.com/sirupsen/logrus"
)

var (
	ErrEmailAlreadyInUse      = errors.New("email already in use")
	ErrUsernameAlreadyInUse   = errors.New("username already in use")
	ErrUserNotFound           = errors.New("user not found")
	ErrInvalidEmailOrPassowrd = errors.New("invalid email or password")
)

type AuthService struct {
	usersRepo       repository.Users
	sessionsRepo    repository.Sessions
	tokenManager    tokens.TokenManagerI
	passwordManager password.PasswordManagerI
	emailSender     email.EmailSender
	refreshTTL      time.Duration
}

func NewAuthService(usersRepo repository.Users, sessionsRepo repository.Sessions,
	tokenManager tokens.TokenManagerI, passwordManager password.PasswordManagerI,
	emailSender email.EmailSender, refreshTTL time.Duration) *AuthService {
	return &AuthService{
		usersRepo:       usersRepo,
		sessionsRepo:    sessionsRepo,
		tokenManager:    tokenManager,
		passwordManager: passwordManager,
		emailSender:     emailSender,
		refreshTTL:      refreshTTL,
	}
}

func parseTime(timestring string) time.Duration {
	timeDuration, err := time.ParseDuration(timestring)
	if err != nil {
		logrus.Fatalf("error parse time duration in auth service: %s", err)
	}
	return timeDuration
}

func (s *AuthService) SignUp(ctx context.Context, user domain.User) (int, error) {
	_, err := s.usersRepo.GetByEmail(ctx, user.Email)
	if err != nil && err != repository.ErrUserNotFound {
		return 0, ErrInternal
	}
	if err == nil {
		return 0, ErrEmailAlreadyInUse
	}
	_, err = s.usersRepo.GetByUserName(ctx, user.Username)
	if err != nil && err != repository.ErrUserNotFound {
		return 0, ErrInternal
	}
	if err == nil {
		return 0, ErrUsernameAlreadyInUse
	}
	user.Password = s.passwordManager.HashPassword(user.Password)
	id, err := s.usersRepo.Create(ctx, user)
	if err != nil {
		return 0, ErrInternal
	}
	return id, nil
}

func (s *AuthService) SignIn(ctx context.Context, email string,
	password string) (domain.Tokens, error) {
	user, err := s.usersRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(ErrUsernameAlreadyInUse, err) {
			return domain.Tokens{}, ErrInvalidEmailOrPassowrd
		}
		return domain.Tokens{}, ErrInternal
	}
	validPassword := s.passwordManager.CheckPassword(password, user.Password)
	if !validPassword {
		return domain.Tokens{}, ErrInvalidEmailOrPassowrd
	}

	accessToken, err := s.tokenManager.CreateAccessToken(user.Id)
	if err != nil {
		return domain.Tokens{}, ErrInternal
	}

	refreshToken, err := s.tokenManager.CreateRefreshToken()
	if err != nil {
		return domain.Tokens{}, ErrInternal
	}
	session := s.createSession(user.Id, refreshToken)
	err = s.sessionsRepo.Create(ctx, session)
	if err != nil {
		return domain.Tokens{}, ErrInternal
	}

	return domain.Tokens{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil
}

func (s *AuthService) createSession(userId int, refreshToken string) domain.Session {
	return domain.Session{
		RefreshToken: refreshToken,
		UserId:       userId,
		ExpiresAt:    time.Now().UTC(),
	}
}
