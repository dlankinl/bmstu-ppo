package base

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateAuthToken(username, jwtKey string) (tokenString string, err error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24),
		})

	tokenString, err = token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("формирование JWT-ключа: %w", err)
	}

	return tokenString, nil
}

func VerifyAuthToken(tokenString, jwtKey string) (err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return fmt.Errorf("парсинг токена: %w", err)
	}

	if !token.Valid {
		return fmt.Errorf("токен невалидный")
	}

	return nil
}
