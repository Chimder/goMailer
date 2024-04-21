package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type Jwt struct {
}

var secretKey = []byte("your_secret_key")

func Encrypt(payload interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	fmt.Println("CLAIM", claims)
	claims["data"] = payload
	fmt.Println("ONECLAIN", claims)
	// claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("Ошибка при подписании токена: %v", err)
	}

	return tokenString, nil
}
func Decrypt(tokenString string) (interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("Ошибка при разборе токена: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["data"], nil
	} else {
		return nil, fmt.Errorf("Недействительный токен")
	}
}
