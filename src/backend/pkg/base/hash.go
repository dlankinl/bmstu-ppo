package base

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type IHashCrypto interface {
	GenerateHashPass(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type HashCrypto struct {
}

func NewHashCrypto() IHashCrypto {
	return HashCrypto{}
}

func (c HashCrypto) GenerateHashPass(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("генерация хэша пароля: %w", err)
	}

	return string(hash), nil
}

func (c HashCrypto) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
