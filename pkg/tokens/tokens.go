package tokens

import (
	"errors"
	"math/rand"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrTokenExpired = errors.New("token is expired")
	ErrTokenInvalid = errors.New("token is invalid")
)

type TokenManagerI interface {
	CreateRefreshToken() (string, error)
	CreateAccessToken(userId int) (string, error)
	ParseAccesToken(tokenString string) (int, error)
}

type TokenManager struct {
	jwt_key   string
	accessTTL time.Duration
}

func NewTokenManager(jwt_key string, accessTTL time.Duration) *TokenManager {
	return &TokenManager{
		jwt_key:   jwt_key,
		accessTTL: accessTTL,
	}
}

func (tm *TokenManager) CreateRefreshToken() (string, error) {
	r := make([]byte, 32)
	source := rand.NewSource(time.Now().Unix())
	res := rand.New(source)
	if _, err := res.Read(r); err != nil {
		return "", nil
	}
	return string(r), nil
}

type StandardClaimsWithUserId struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func (tm *TokenManager) CreateAccessToken(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &StandardClaimsWithUserId{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tm.accessTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userId,
	})
	return token.SignedString([]byte(tm.jwt_key))
}

func (tm *TokenManager) ParseAccesToken(tokenString string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &StandardClaimsWithUserId{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrTokenInvalid
			}
			claims, ok := t.Claims.(*StandardClaimsWithUserId)
			if !ok {
				return nil, ErrTokenInvalid
			}
			if claims.ExpiresAt <= time.Now().Unix() {
				return nil, ErrTokenExpired
			}
			return []byte(tm.jwt_key), nil
		})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*StandardClaimsWithUserId)
	if !ok {
		return 0, ErrTokenInvalid
	}

	return claims.UserId, nil
}
