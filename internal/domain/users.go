package domain

import (
	"errors"
)

var (
	ErrUserUpdateHasNoValues = errors.New("user update has no values")
)

const (
	RoleUser    = "user"
	RoleManager = "manager"
	RoleSeller  = "seller"
)

type User struct {
	Id            int    `json:"id"`
	Username      string `json:"username" validate:"required"`
	Name          string `json:"name" validate:"required"`
	Email         string `json:"email" validate:"required,email"`
	Password      string `json:"password" db:"hash_password" validate:"required"`
	EmailVerified bool   `json:"emailVerified" db:"email_verified"`
	Role          string `json:"role" db:"role"`
}

type UserUpdate struct {
	Username      *string
	Name          *string
	Email         *string
	Password      *string
	EmailVefiried *bool
	Role          *string
}

func (u *UserUpdate) Validate() error {
	if u.Username == nil && u.Name == nil && u.Email == nil && u.Password == nil && u.Role == nil {
		return ErrUserUpdateHasNoValues
	}
	return nil
}
