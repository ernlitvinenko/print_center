package storage

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword генерирует bcrypt-хеш из пароля.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

// checkPassword сравнивает пароль с bcrypt-хешем.
func checkPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
