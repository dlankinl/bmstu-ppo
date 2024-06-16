package base

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JwtPayload struct {
	ID   string
	Role string
}

func GenerateAuthToken(id, jwtKey, role string) (tokenString string, err error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub":  id,
			"exp":  time.Now().Add(time.Hour * 24).Unix(),
			"role": role,
		})

	tokenString, err = token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", fmt.Errorf("формирование JWT-ключа: %w", err)
	}

	return tokenString, nil
}

func VerifyAuthToken(tokenString, jwtKey string) (payload *JwtPayload, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("парсинг токена: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("токен невалидный")
	}

	payload = new(JwtPayload)
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		payload.ID = fmt.Sprint(claims["sub"])
		payload.Role = fmt.Sprint(claims["role"])
	}

	return payload, nil
}
