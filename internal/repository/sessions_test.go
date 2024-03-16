package repository

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/domain"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestSessions_Create(t *testing.T) {
	rdb, mock := redismock.NewClientMock()

	sessionsRepository := NewSessionsRepository(rdb)

	type args struct {
		session domain.Session
	}

	type mockBehavior func(args args)

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      error
	}{
		{
			name: "ok",
			args: args{
				session: domain.Session{
					UserId:       1,
					RefreshToken: "some-refresh-token",
					ExpiresAt:    time.Now().Add(24 * time.Hour),
				},
			},
			mockBehavior: func(args args) {
				userSessionsKey := sessionsRepository.getUserSessionsKey(args.session.UserId)
				sessionKey := sessionsRepository.getSessionKey(args.session.RefreshToken)

				mock.ExpectTxPipeline()

				mock.ExpectZAdd(userSessionsKey, redis.Z{
					Score:  0,
					Member: args.session.RefreshToken,
				}).SetVal(1)

				mock.ExpectHSet(sessionKey, "userId", args.session.UserId).SetVal(1)
				mock.ExpectExpireAt(sessionKey, args.session.ExpiresAt).SetVal(true)

				mock.ExpectTxPipelineExec()
			},
			wantErr: nil,
		},
		{
			name: "err add item to usersessions",
			args: args{
				session: domain.Session{
					UserId:       1,
					RefreshToken: "some-refresh-token",
					ExpiresAt:    time.Now().Add(24 * time.Hour),
				},
			},
			mockBehavior: func(args args) {
				userSessionKey := sessionsRepository.getUserSessionsKey(args.session.UserId)

				mock.ExpectTxPipeline()

				mock.ExpectZAdd(userSessionKey, redis.Z{
					Score:  0,
					Member: args.session.RefreshToken,
				}).SetErr(errors.New("some redis error"))
			},
			wantErr: ErrInternal,
		},
		{
			name: "err add session",
			args: args{
				session: domain.Session{
					UserId:       1,
					RefreshToken: "some-refresh-token",
					ExpiresAt:    time.Now().Add(24 * time.Hour),
				},
			},
			mockBehavior: func(args args) {
				userSessionKey := sessionsRepository.getUserSessionsKey(args.session.UserId)
				sessionKey := sessionsRepository.getSessionKey(args.session.RefreshToken)

				mock.ExpectTxPipeline()

				mock.ExpectZAdd(userSessionKey, redis.Z{
					Score:  0,
					Member: args.session.RefreshToken,
				}).SetVal(1)

				mock.ExpectHSet(sessionKey, "userId", args.session.UserId).
					SetErr(errors.New("some redis error"))
			},
			wantErr: ErrInternal,
		},
		{
			name: "err add ttl session",
			args: args{
				session: domain.Session{
					UserId:       1,
					RefreshToken: "some-refresh-token",
					ExpiresAt:    time.Now().Add(24 * time.Hour),
				},
			},
			mockBehavior: func(args args) {
				userSessionKey := sessionsRepository.getUserSessionsKey(args.session.UserId)
				sessionKey := sessionsRepository.getSessionKey(args.session.RefreshToken)

				mock.ExpectTxPipeline()

				mock.ExpectZAdd(userSessionKey, redis.Z{
					Score:  0,
					Member: args.session.RefreshToken,
				}).SetVal(1)

				mock.ExpectHSet(sessionKey, "userId", args.session.UserId).SetVal(1)

				mock.ExpectExpireAt(sessionKey, args.session.ExpiresAt).
					SetErr(errors.New("some redis error"))
			},
			wantErr: ErrInternal,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.args)

			err := sessionsRepository.Create(context.Background(), test.args.session)

			if test.wantErr != nil {
				assert.ErrorIs(t, test.wantErr, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSessions_Get(t *testing.T) {
	rdb, mock := redismock.NewClientMock()

	sessionsRepository := NewSessionsRepository(rdb)

	type args struct {
		refreshToken string
	}

	type mockBehavior func(args args, session domain.Session)

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		want         domain.Session
		wantErr      error
	}{
		{
			name: "ok",
			args: args{
				refreshToken: "some-refresh-token",
			},
			mockBehavior: func(args args, session domain.Session) {
				mock.ExpectHGetAll(sessionsRepository.getSessionKey(args.refreshToken)).
					SetVal(map[string]string{
						"userId": strconv.Itoa(session.UserId),
					})
				mock.ExpectTTL(sessionsRepository.getSessionKey(args.refreshToken)).
					SetVal(time.Until(session.ExpiresAt))
			},
			want: domain.Session{
				UserId:       1,
				RefreshToken: "some-refresh-token",
				ExpiresAt:    time.Now().UTC().Add(24 * time.Hour),
			},
			wantErr: nil,
		},
		{
			name: "err get ttl",
			args: args{
				refreshToken: "some-refresh-token-durachek",
			},
			mockBehavior: func(args args, session domain.Session) {
				mock.ExpectHGetAll(sessionsRepository.getSessionKey(args.refreshToken)).
					SetVal(map[string]string{
						"userId": strconv.Itoa(session.UserId),
					})
				mock.ExpectTTL(sessionsRepository.getSessionKey(args.refreshToken)).
					SetErr(errors.New("some redis errror"))
			},
			want:    domain.Session{},
			wantErr: ErrInternal,
		},
		{
			name: "session doesn't exist",
			args: args{
				refreshToken: "durachek",
			},
			mockBehavior: func(args args, session domain.Session) {
				mock.ExpectHGetAll(sessionsRepository.getSessionKey(args.refreshToken)).
					SetVal(nil)
			},
			want:    domain.Session{},
			wantErr: ErrSessionExpiredOrInvalid,
		},
		{
			name: "error get session",
			args: args{
				refreshToken: "lox",
			},
			mockBehavior: func(args args, session domain.Session) {
				mock.ExpectHGetAll(sessionsRepository.getSessionKey(args.refreshToken)).
					SetErr(errors.New("some redis error"))
			},
			want:    domain.Session{},
			wantErr: ErrInternal,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.args, test.want)

			got, err := sessionsRepository.Get(context.Background(), test.args.refreshToken)

			assert.ErrorIs(t, test.wantErr, err)
			if test.wantErr == nil {
				assert.Equal(t, test.want.RefreshToken, got.RefreshToken)
				assert.Equal(t, test.want.UserId, got.UserId)
				assert.Equal(t, true, got.ExpiresAt.Sub(test.want.ExpiresAt) < time.Second)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSessions_Delete(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	SessionsRepository := NewSessionsRepository(rdb)

	type args struct {
		userId       int
		refreshToken string
	}

	type mockBehavior func(args args)

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      error
	}{
		{
			name: "ok",
			args: args{
				userId:       1,
				refreshToken: "ok-refresh-token",
			},
			mockBehavior: func(args args) {
				userSessionsKey := SessionsRepository.getUserSessionsKey(args.userId)
				sessionKey := SessionsRepository.getSessionKey(args.refreshToken)
				mock.ExpectTxPipeline()
				mock.ExpectZRem(userSessionsKey, args.refreshToken).SetVal(1)
				mock.ExpectDel(sessionKey).SetVal(1)
				mock.ExpectTxPipelineExec()
			},
			wantErr: nil,
		},
		{
			name: "err delete session from user sessions list",
			args: args{
				userId:       1,
				refreshToken: "lox",
			},
			mockBehavior: func(args args) {
				userSessionsKey := SessionsRepository.getUserSessionsKey(args.userId)
				mock.ExpectTxPipeline()
				mock.ExpectZRem(userSessionsKey, args.refreshToken).
					SetErr(errors.New("some redis error"))
			},
			wantErr: ErrInternal,
		},
		{
			name: "err delete session",
			args: args{
				userId:       1,
				refreshToken: "hahaha",
			},
			mockBehavior: func(args args) {
				userSessionsKey := SessionsRepository.getUserSessionsKey(args.userId)
				sessionKey := SessionsRepository.getSessionKey(args.refreshToken)
				mock.ExpectTxPipeline()
				mock.ExpectZRem(userSessionsKey, args.refreshToken).SetVal(1)
				mock.ExpectDel(sessionKey).SetErr(errors.New("some redis error"))
			},
			wantErr: ErrInternal,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.args)

			err := SessionsRepository.Delete(context.Background(),
				test.args.userId, test.args.refreshToken)

			assert.ErrorIs(t, test.wantErr, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSessions_GetCnt(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	SessionsRepository := NewSessionsRepository(rdb)

	type args struct {
		userId int
	}

	type mockBehavior func(args args, cnt int)

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		want         int
		wantErr      error
	}{
		{
			name: "ok",
			args: args{
				userId: 10,
			},
			mockBehavior: func(args args, cnt int) {
				userSessionsKey := SessionsRepository.getUserSessionsKey(args.userId)

				mock.ExpectZCard(userSessionsKey).SetVal(int64(cnt))
			},
			want:    3,
			wantErr: nil,
		},
		{
			name: "error get cnt",
			args: args{
				userId: 10,
			},
			mockBehavior: func(args args, cnt int) {
				userSessionsKey := SessionsRepository.getUserSessionsKey(args.userId)

				mock.ExpectZCard(userSessionsKey).SetErr(errors.New("some redis error"))
			},
			want:    0,
			wantErr: ErrInternal,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.args, test.want)

			got, err := SessionsRepository.GetCnt(context.Background(), test.args.userId)

			assert.ErrorIs(t, test.wantErr, err)
			if test.wantErr == nil {
				assert.Equal(t, test.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSessions_GetAll(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	sessionsRepository := NewSessionsRepository(rdb)

	type args struct {
		userId int
	}

	type mockBehavior func(args args, sessions []domain.Session)

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		want         []domain.Session
		wantErr      error
	}{
		{
			name: "ok",
			args: args{
				userId: 1,
			},
			mockBehavior: func(args args, sessions []domain.Session) {
				userSessionKey := sessionsRepository.getUserSessionsKey(args.userId)
				mock.ExpectZCard(userSessionKey).SetVal(int64(len(sessions)))
				refreshTokens := make([]string, len(sessions))
				for i := 0; i < len(sessions); i++ {
					refreshTokens[i] = sessions[i].RefreshToken
				}
				mock.ExpectZRange(userSessionKey, 0, int64(len(sessions))).SetVal(refreshTokens)
				for i := 0; i < len(sessions); i++ {
					sessionKey := sessionsRepository.getSessionKey(sessions[i].RefreshToken)
					mock.ExpectHGetAll(sessionKey).SetVal(map[string]string{
						"userId": strconv.Itoa(sessions[i].UserId),
					})
					mock.ExpectTTL(sessionKey).SetVal(time.Until(sessions[i].ExpiresAt))
				}
			},
			want: []domain.Session{
				{
					UserId:       1,
					RefreshToken: "refresh-token-1",
					ExpiresAt:    time.Now().UTC().Add(1 * time.Hour),
				},
				{
					UserId:       1,
					RefreshToken: "refresh-token-2",
					ExpiresAt:    time.Now().UTC().Add(2 * time.Hour),
				},
				{
					UserId:       1,
					RefreshToken: "refresh-token-3",
					ExpiresAt:    time.Now().UTC().Add(3 * time.Hour),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.args, test.want)

			got, err := sessionsRepository.GetAll(context.Background(), test.args.userId)

			assert.ErrorIs(t, test.wantErr, err)
			if test.wantErr == nil {
				assert.Equal(t, len(test.want), len(got))
				if len(test.want) == len(got) {
					for i, v := range got {
						assert.Equal(t, test.want[i].UserId, v.UserId)
						assert.Equal(t, test.want[i].RefreshToken, v.RefreshToken)
						assert.True(t, v.ExpiresAt.Sub(v.ExpiresAt) < time.Second)
					}
				}
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestSessions_DeleteAll(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	SessionsRepository := NewSessionsRepository(rdb)

	type args struct {
		userId int
	}

	type mockBehavior func(args args)

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      error
	}{
		{
			name: "ok",
			args: args{
				userId: 5,
			},
			mockBehavior: func(args args) {
				refreshTokens := []string{
					"srt1", "srt2", "srt3",
				}
				userSessionsKey := SessionsRepository.getUserSessionsKey(args.userId)
				mock.ExpectZCard(userSessionsKey).SetVal(int64(len(refreshTokens)))

				sessionsKeys := make([]string, len(refreshTokens))
				for i := 0; i < len(sessionsKeys); i++ {
					sessionsKeys[i] = SessionsRepository.getSessionKey(refreshTokens[i])
				}

				mock.ExpectZRange(userSessionsKey, 0, int64(len(refreshTokens))).SetVal(refreshTokens)
				mock.ExpectTxPipeline()
				mock.ExpectDel(userSessionsKey).SetVal(1)
				mock.ExpectDel(sessionsKeys...).SetVal(int64(len(sessionsKeys)))
				mock.ExpectTxPipelineExec()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.args)

			err := SessionsRepository.DeleteAll(context.Background(), test.args.userId)

			assert.ErrorIs(t, test.wantErr, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
