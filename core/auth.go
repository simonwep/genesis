package core

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type JWTClaim struct {
	User string `json:"user"`
	jwt.RegisteredClaims
}

func CreateAccessToken(user *User) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaim{
		User: user.User,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(Config.JWTAccessExpiration)),
		},
	}).SignedString(Config.JWTSecret)
}

func CreateRefreshToken(user *User) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaim{
		User: user.User,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(Config.JWTRefreshExpiration)),
			ID:        uuid.NewString(),
		},
	}).SignedString(Config.JWTSecret)
}

func ParseAuthToken(token string) (*JWTClaim, error) {
	var claims JWTClaim

	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return Config.JWTSecret, nil
	})

	if len(claims.ID) != 0 {
		blacklisted, err := IsTokenBlacklisted(claims.ID)

		if blacklisted || err != nil {
			return nil, err
		}
	}

	return &claims, err
}
