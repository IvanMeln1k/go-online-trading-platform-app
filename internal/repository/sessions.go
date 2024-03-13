package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/IvanMeln1k/go-online-trading-platform-app/internal/domain"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

type SessionsRepository struct {
	rdb *redis.Client
}

func NewSessionsRepository(rdb *redis.Client) *SessionsRepository {
	return &SessionsRepository{
		rdb: rdb,
	}
}

// Returns the session key in redis by refreshToken
func (r *SessionsRepository) getSessionKey(refreshToken string) string {
	return fmt.Sprintf("sessions:%s", refreshToken)
}

// Returns key of the user's sessions list
func (r *SessionsRepository) getUserSessionsKey(userId int) string {
	return fmt.Sprintf("userSessions:%d", userId)
}

func (r *SessionsRepository) Create(ctx context.Context, session domain.Session) error {
	pipe := r.rdb.Pipeline()

	sessionKey := r.getSessionKey(session.RefreshToken)
	userSessionsKey := r.getUserSessionsKey(session.UserId)

	_, err := pipe.ZAdd(userSessionsKey, redis.Z{
		Score:  0,
		Member: session.RefreshToken,
	}).Result()
	if err != nil {
		logrus.Errorf("error add sessions into userSessions: %s", err)
		pipe.Discard()
		return ErrInternal
	}

	_, err = pipe.HSet(sessionKey, "id", session.UserId).Result()
	if err != nil {
		logrus.Errorf("error create session into redis: %s", err)
		pipe.Discard()
		return ErrInternal
	}

	_, err = pipe.ExpireAt(sessionKey, session.ExpiresAt).Result()
	if err != nil {
		logrus.Errorf("error set expiresat session: %s", err)
		pipe.Discard()
		return ErrInternal
	}

	_, err = pipe.Exec()
	if err != nil {
		logrus.Errorf("error exec transactoin: %s", err)
		pipe.Discard()
		return ErrInternal
	}

	return nil
}

func (r *SessionsRepository) bindSession(session *domain.Session, sessionMap map[string]string) error {
	val, ok := sessionMap["userId"]
	if !ok {
		return errors.New("sessionMap hasn't userId")
	}
	userId, err := strconv.Atoi(val)
	if err != nil {
		return errors.New("sessionsMap's userId isn't int")
	}

	session.UserId = userId

	return nil
}

func (r *SessionsRepository) Get(ctx context.Context, refreshToken string) (domain.Session, error) {
	var session domain.Session

	sessionKey := r.getSessionKey(refreshToken)

	res, err := r.rdb.HGetAll(sessionKey).Result()
	if err != nil {
		logrus.Errorf("")
		return session, ErrInternal
	}

	err = r.bindSession(&session, res)
	if err != nil {
		logrus.Errorf("error bind session: %s", err)
		return session, ErrSessionExpiredOrInvalid
	}

	ttl, err := r.rdb.TTL(refreshToken).Result()
	if err != nil {
		logrus.Errorf("error get ttl session: %s", err)
		return session, ErrSessionExpiredOrInvalid
	}

	session.ExpiresAt = time.Now().UTC().Add(ttl)
	session.RefreshToken = refreshToken

	return session, nil
}

func (r *SessionsRepository) Delete(ctx context.Context, userId int, refreshToken string) error {
	pipe := r.rdb.Pipeline()

	userSessionsKey := r.getUserSessionsKey(userId)
	sessionKey := r.getSessionKey(refreshToken)

	_, err := pipe.ZRem(userSessionsKey, sessionKey).Result()
	if err != nil {
		pipe.Discard()
		logrus.Errorf("error del session from userSessions: %s", err)
		return ErrInternal
	}

	_, err = pipe.Del(userSessionsKey).Result()
	if err != nil {
		pipe.Discard()
		logrus.Errorf("error del session: %s", err)
		return ErrInternal
	}

	_, err = pipe.Exec()
	if err != nil {
		logrus.Errorf("error exec transaction redis: %s", err)
		return ErrInternal
	}

	return nil
}

func (r *SessionsRepository) GetCnt(ctx context.Context, userId int) (int, error) {
	var cnt int

	userSessionsKey := r.getUserSessionsKey(userId)

	_, err := r.rdb.ZCard(userSessionsKey).Result()
	if err != nil {
		logrus.Errorf("error get zcard user sessions: %s", err)
		return cnt, ErrInternal
	}

	return cnt, nil
}

func (r *SessionsRepository) GetAll(ctx context.Context, userId int) ([]domain.Session, error) {
	var sessions []domain.Session

	cnt, err := r.GetCnt(ctx, userId)
	if err != nil {
		return nil, ErrInternal
	}

	userSessionsKey := r.getUserSessionsKey(userId)

	refreshTokens, err := r.rdb.ZRange(userSessionsKey, 0, int64(cnt)).Result()
	if err != nil {
		logrus.Errorf("error get refresh tokens by id from redis: %s", err)
		return nil, ErrInternal
	}
	sessionKeys := make([]string, len(refreshTokens))
	for i := 0; i < len(sessionKeys); i++ {
		sessionKeys[i] = r.getSessionKey(refreshTokens[i])
	}

	for i := 0; i < len(sessionKeys); i++ {
		session, err := r.Get(ctx, refreshTokens[i])
		if err != nil {
			if errors.Is(err, ErrSessionExpiredOrInvalid) {
				_, err := r.rdb.Del(sessionKeys[i]).Result()
				if err != nil {
					logrus.Errorf("error delete invalid session: %s", err)
					return nil, ErrInternal
				}
				continue
			}
			logrus.Errorf("error get session by token: %s", err)
			return nil, ErrInternal
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (r *SessionsRepository) DeleteAll(ctx context.Context, userId int) error {
	pipe := r.rdb.Pipeline()

	userSessionsKey := r.getUserSessionsKey(userId)

	sessions, err := r.GetAll(ctx, userId)
	if err != nil {
		logrus.Errorf("error get user sessions: %s", err)
	}

	_, err = pipe.Del(userSessionsKey).Result()
	if err != nil {
		logrus.Errorf("error del user sessions: %s", err)
		pipe.Discard()
		return ErrInternal
	}

	sessionsKeys := make([]string, len(sessions))
	for i := 0; i < len(sessions); i++ {
		sessionsKeys[i] = r.getSessionKey(sessions[i].RefreshToken)
	}
	_, err = pipe.Del(sessionsKeys...).Result()
	if err != nil {
		logrus.Errorf("error del sessions: %s", err)
		pipe.Discard()
		return ErrInternal
	}

	return nil
}
