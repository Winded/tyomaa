package services

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/winded/tyomaa/backend/db"
	"github.com/winded/tyomaa/backend/util/token"
)

func AuthenticateUser(name, password string) *db.User {
	var user db.User
	if err := db.Instance.Where("name = ?", name).First(&user).Error; err != nil {
		return nil
	}

	if user.CheckPassword(password) != nil {
		return nil
	}

	return &user
}

func CreateToken(user *db.User) (string, error) {
	expTime := time.Now().Add(time.Hour * 24 * 30)

	claims := token.Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime.Unix(),
		},
	}

	return token.Create(claims)
}
