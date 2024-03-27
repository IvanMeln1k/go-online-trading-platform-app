package password

import (
	"crypto/sha256"
	"fmt"
)

type PasswordManagerI interface {
	HashPassword(password string) string
	CheckPassword(password string, hashedPassword string) bool
}

type PasswordManager struct {
	salt string
}

func NewPasswordManager(salt string) *PasswordManager {
	return &PasswordManager{
		salt: salt,
	}
}

func (pm *PasswordManager) HashPassword(password string) string {
	bytes := sha256.Sum256([]byte(password + pm.salt))
	return fmt.Sprintf("%x", bytes)
}

func (pm *PasswordManager) CheckPassword(password string, hashedPassword string) bool {
	return pm.HashPassword(password) == hashedPassword
}
