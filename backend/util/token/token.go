package token

import (
	"errors"
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
)

var Secret []byte

func init() {
	hmacSec, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		hmacSec = "test"
	}
	Secret = []byte(hmacSec)
}

type Claims struct {
	UserID uint `json:"userId"`
	jwt.StandardClaims
}

func Create(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(Secret)
}

func Verify(token string) (Claims, error) {
	t, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return Secret, nil
	})
	if err != nil {
		return Claims{}, err
	}

	claims, ok := t.Claims.(*Claims)
	if !ok || !t.Valid {
		return *claims, errors.New("Invalid token")
	}

	return *claims, nil
}
