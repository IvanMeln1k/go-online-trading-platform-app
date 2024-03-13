package domain

import (
	"errors"
)

var (
	ErrUserUpdateHasNoValues = errors.New("user update has no values")
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password" db:"hash_password"`
}

type UserUpdate struct {
	Username *string `json:"username"`
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func (u *UserUpdate) Validate() error {
	if u.Username == nil && u.Name == nil && u.Email == nil && u.Password == nil {
		return ErrUserUpdateHasNoValues
	}
	return nil
}
