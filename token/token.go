package token

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Token struct {
	JWTSecret []byte
}

func (t *Token) GenerateToken(user_id int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user_id,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Hour + 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(t.JWTSecret)
}

func (t *Token) ParseToken(tokenString string) (int, error) {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, err := token.Method.(*jwt.SigningMethodHMAC); !err {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return t.JWTSecret, nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return 0, fmt.Errorf("Invalid token")

	}
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid `user_id` claim")
	}
	return int(userIDFloat), nil
}
