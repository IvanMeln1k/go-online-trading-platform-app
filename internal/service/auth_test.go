package service

import (
	"context"
	"fmt"
	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/domain"
	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/repository"
	mock_repository "github.com/IvanMeln1k/go-online-trading-platform-app/internal/repository/mocks"
	mock_email "github.com/IvanMeln1k/go-online-trading-platform-app/pkg/email/mocks"
	mock_password "github.com/IvanMeln1k/go-online-trading-platform-app/pkg/password/mocks"
	mock_tokens "github.com/IvanMeln1k/go-online-trading-platform-app/pkg/tokens/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAuth_SignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersRepo := mock_repository.NewMockUsers(ctrl)
	sessionsRepo := mock_repository.NewMockSessions(ctrl)
	tokenManager := mock_tokens.NewMockTokenManagerI(ctrl)
	passwordManager := mock_password.NewMockPasswordManagerI(ctrl)
	emailSender := mock_email.NewMockEmailSender(ctrl)
	refreshTTL := 30 * 24 * time.Hour
	verificationAddr := "some-verification-addr"

	authService := NewAuthService(usersRepo, sessionsRepo, tokenManager, passwordManager, emailSender,
		refreshTTL, verificationAddr)

	type args struct {
		user domain.User
	}

	type mockBehavior func(args args, id int)

	tests := []struct {
		name         string
		args         args
		want         int
		mockBehavior mockBehavior
		wantErr      bool
		err          error
	}{
		{
			name: "ok",
			args: args{
				user: domain.User{
					Username: "Ivan",
					Name:     "Ivan",
					Email:    "IvanMelnikovF@gmail.com",
					Password: "pass",
				},
			},
			want: 1,
			mockBehavior: func(args args, id int) {
				usersRepo.EXPECT().GetByEmail(context.Background(), args.user.Email).
					Return(domain.User{}, repository.ErrUserNotFound)
				usersRepo.EXPECT().GetByUserName(context.Background(), args.user.Username).
					Return(domain.User{}, repository.ErrUserNotFound)

				hashPass := "hash"
				passwordManager.EXPECT().HashPassword(args.user.Password).Return(hashPass)

				args.user.Password = hashPass
				usersRepo.EXPECT().Create(context.Background(), args.user).Return(id, nil)

				emailToken := "emailToken"
				tokenManager.EXPECT().CreateEmailToken(args.user.Email).Return(emailToken, nil)

				emailSender.EXPECT().Send("templates/verification.html", args.user.Email,
					"GO Online-Trading-Platform verification email", map[string]string{
						"Link": fmt.Sprintf("%s?email=%s", verificationAddr, emailToken),
					}).Return(nil)
			},
		},
		{
			name: "email already in use",
			args: args{
				user: domain.User{
					Username: "Ivan",
					Name:     "Ivan",
					Email:    "IvanMelnikoF@gmail.com",
					Password: "pass",
				},
			},
			want: 0,
			mockBehavior: func(args args, id int) {
				usersRepo.EXPECT().GetByEmail(context.Background(), args.user.Email).
					Return(domain.User{}, nil)
			},
			wantErr: true,
			err:     ErrEmailAlreadyInUse,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.args, test.want)

			got, err := authService.SignUp(context.Background(), test.args.user)

			if test.wantErr {
				assert.ErrorIs(t, err, test.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			}
		})
	}
}
