package core

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTClaim struct {
	User string `json:"user"`
	jwt.RegisteredClaims
}

func CreateAccessToken(user *User) (string, error) {
	return createToken(user, time.Now().Add(Config.JWTAccessExpiration))
}

func CreateRefreshToken(user *User) (string, error) {
	return createToken(user, time.Now().Add(Config.JWTRefreshExpiration))
}

func ParseAuthToken(token string) (*JWTClaim, error) {
	var claims JWTClaim

	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return Config.JWTSecret, nil
	})

	return &claims, err
}

func createToken(user *User, expires time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaim{
		User: user.User,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
		},
	})

	return token.SignedString(Config.JWTSecret)
}
