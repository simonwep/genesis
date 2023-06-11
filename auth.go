package main

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTClaim struct {
	User string `json:"user"`
	jwt.RegisteredClaims
}

func CreateAuthToken(user *User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaim{
		User: user.User,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(Env().JWTExpires)),
		},
	})

	return token.SignedString(Env().JWTSecret)
}

func ParseAuthToken(token string) (*JWTClaim, error) {
	var claims JWTClaim

	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return Env().JWTSecret, nil
	})

	return &claims, err
}
