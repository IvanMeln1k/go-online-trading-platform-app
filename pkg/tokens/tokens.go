package tokens

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrTokenExpired = errors.New("token is expired")
	ErrTokenInvalid = errors.New("token is invalid")
)

//go:generate mockgen -source=tokens.go -destination=mocks/mock.go

type TokenManagerI interface {
	CreateRefreshToken() (string, error)
	CreateAccessToken(userId int) (string, error)
	ParseAccessToken(tokenString string) (int, error)
	CreateEmailToken(email string) (string, error)
	ParseEmailToken(tokenString string) (string, error)
}

type TokenManager struct {
	jwt_key   string
	accessTTL time.Duration
	emailTTL  time.Duration
}

func NewTokenManager(jwt_key string, accessTTL time.Duration, emailTTL time.Duration) *TokenManager {
	return &TokenManager{
		jwt_key:   jwt_key,
		accessTTL: accessTTL,
		emailTTL:  emailTTL,
	}
}

func (tm *TokenManager) CreateRefreshToken() (string, error) {
	b := make([]byte, 32)

	src := rand.NewSource(time.Now().Unix())
	r := rand.New(src)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

type StandardClaimsWithUserId struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type StandardClaimsWithEmail struct {
	jwt.StandardClaims
	Email string `json:"email"`
}

func (tm *TokenManager) createJWTToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tm.jwt_key))
}

func (tm *TokenManager) getStandartClaims(ttl time.Duration) jwt.StandardClaims {
	return jwt.StandardClaims{
		ExpiresAt: time.Now().Add(ttl).Unix(),
		IssuedAt:  time.Now().Unix(),
	}
}

func (tm *TokenManager) CreateEmailToken(email string) (string, error) {
	return tm.createJWTToken(&StandardClaimsWithEmail{
		tm.getStandartClaims(tm.emailTTL),
		email,
	})
}

func (tm *TokenManager) CreateAccessToken(userId int) (string, error) {
	return tm.createJWTToken(&StandardClaimsWithUserId{
		tm.getStandartClaims(tm.accessTTL),
		userId,
	})
}

func (tm *TokenManager) ParseEmailToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &StandardClaimsWithEmail{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrTokenInvalid
			}
			claims, ok := t.Claims.(*StandardClaimsWithEmail)
			if !ok {
				return nil, ErrTokenInvalid
			}
			if claims.ExpiresAt <= time.Now().Unix() {
				return nil, ErrTokenExpired
			}
			return []byte(tm.jwt_key), nil
		})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*StandardClaimsWithEmail)
	if !ok {
		return "", err
	}

	return claims.Email, nil
}

func (tm *TokenManager) ParseAccessToken(tokenString string) (int, error) {
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
