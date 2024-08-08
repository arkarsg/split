package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword creats a bcrypt hash string with cost=14
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Failed to hash password: %w", err)
	}
	return string(bytes), err
}

// CheckPasswordHash checks if the `input` password is equals to bcrypt `hash` string
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
