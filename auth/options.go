package auth

import (
	"encoding/json"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type Jwt struct {
}

var secretKey = []byte("your_secret_key")

func Encrypt(payload interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["data"] = payload

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("Ошибка при подписании токена: %v", err)
	}

	return tokenString, nil
}

func Decrypt(tokenString string, v interface{}) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return fmt.Errorf("Ошибка при разборе токена: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		data := claims["data"].(map[string]interface{})
		bytes, _ := json.Marshal(data)
		json.Unmarshal(bytes, &v)
		return nil
	} else {
		return fmt.Errorf("Недействительный токен")
	}
}
