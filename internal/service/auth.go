package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/domain"
	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/repository"
	"github.com/IvanMeln1k/go-online-trading-platform-app/pkg/email"
	"github.com/IvanMeln1k/go-online-trading-platform-app/pkg/password"
	"github.com/IvanMeln1k/go-online-trading-platform-app/pkg/tokens"
	"github.com/sirupsen/logrus"
)

var (
	ErrEmailAlreadyInUse       = errors.New("email already in use")
	ErrUsernameAlreadyInUse    = errors.New("username already in use")
	ErrUserNotFound            = errors.New("user not found")
	ErrInvalidEmailOrPassowrd  = errors.New("invalid email or password")
	ErrSendEmailVerification   = errors.New("error send email verification")
	ErrSessionInvalidOrExpired = errors.New("session is invalid or expired")
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

	emailToken, err := s.tokenManager.CreateEmailToken(user.Email)
	if err != nil {
		return 0, ErrSendEmailVerification
	}
	err = s.emailSender.Send("templates/verification.html", user.Email,
		"GO Online-Trading-Platform verification email", map[string]string{
			"Link": fmt.Sprintf("localhost:8000/auth/verify?email=%s", emailToken),
		})
	if err != nil {
		return 0, ErrSendEmailVerification
	}

	return id, nil
}

func (s *AuthService) VerifyEmail(ctx context.Context, email string) error {
	user, err := s.usersRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(repository.ErrUserNotFound, err) {
			return ErrUserNotFound
		}
		return ErrInternal
	}
	if user.EmailVerified {
		return ErrInternal
	}
	var emailVefiried bool
	emailVefiried = true
	_, err = s.usersRepo.Update(ctx, user.Id, domain.UserUpdate{
		EmailVefiried: &emailVefiried,
	})
	if err != nil {
		return ErrInternal
	}
	return nil
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

	cntSessions, err := s.sessionsRepo.GetCnt(ctx, user.Id)
	if err != nil {
		logrus.Errorf("error get cnt session in auth service: %s", err)
	}
	if err == nil {
		if cntSessions > 5 {
			logrus.Printf("cntsessions > 5")
			sessions, err := s.sessionsRepo.GetAll(ctx, user.Id)
			if err != nil {
				logrus.Errorf("error get all sessions in auth service: %s", err)
			}
			if err == nil {
				domain.SortSessionsByTime(sessions)
				cntSessions = len(sessions)
				for cntSessions > 5 {
					err = s.sessionsRepo.Delete(ctx, user.Id, sessions[cntSessions-1].RefreshToken)
					if err != nil {
						logrus.Errorf("error delete session when added new session: %s", err)
						break
					}
					cntSessions--
					sessions = sessions[:cntSessions]
				}
			}
		}
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
		ExpiresAt:    time.Now().UTC().Add(s.refreshTTL),
	}
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	session, err := s.sessionsRepo.Get(ctx, refreshToken)
	if err != nil {
		logrus.Errorf("error delete session: %s", err)
		if errors.Is(repository.ErrSessionExpiredOrInvalid, err) {
			return ErrSessionInvalidOrExpired
		}
		return ErrInternal
	}

	err = s.sessionsRepo.Delete(ctx, session.UserId, refreshToken)
	if err != nil {
		logrus.Errorf("error delete session: %s", err)
		return ErrInternal
	}

	return nil
}

func (s *AuthService) LogoutAll(ctx context.Context, refreshToken string) error {
	session, err := s.sessionsRepo.Get(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, repository.ErrSessionExpiredOrInvalid) {
			return ErrSessionInvalidOrExpired
		}
		return ErrInternal
	}
	err = s.sessionsRepo.DeleteAll(ctx, session.UserId)
	if err != nil {
		return ErrInternal
	}
	return nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (domain.Tokens, error) {
	session, err := s.sessionsRepo.Get(ctx, refreshToken)
	if err != nil {
		logrus.Errorf("error get session when refreshing: %s", err)
		if errors.Is(repository.ErrSessionExpiredOrInvalid, err) {
			return domain.Tokens{}, ErrSessionInvalidOrExpired
		}
		return domain.Tokens{}, ErrInternal
	}

	accessToken, err := s.tokenManager.CreateAccessToken(session.UserId)
	if err != nil {
		logrus.Errorf("error create new access token: %s", err)
		return domain.Tokens{}, ErrInternal
	}
	newRefreshToken, err := s.tokenManager.CreateRefreshToken()
	if err != nil {
		logrus.Errorf("error create new refreshToken: %s", err)
		return domain.Tokens{}, ErrInternal
	}

	err = s.sessionsRepo.Delete(ctx, session.UserId, session.RefreshToken)
	if err != nil {
		logrus.Errorf("error delete session: %s", err)
		return domain.Tokens{}, ErrInternal
	}

	session.RefreshToken = newRefreshToken
	session.ExpiresAt = time.Now().UTC().Add(s.refreshTTL)

	err = s.sessionsRepo.Create(ctx, session)
	if err != nil {
		logrus.Errorf("error create new session when refreshing: %s", err)
		return domain.Tokens{}, ErrInternal
	}

	return domain.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
