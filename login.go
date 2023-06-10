package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type LoginBody struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type JWTClaim struct {
	User string `json:"user"`
	jwt.RegisteredClaims
}

func Login(c *gin.Context) {
	var body LoginBody

	if err := c.BindJSON(&body); err != nil {
		c.Status(http.StatusBadRequest)
	} else if !ValidateUserName(body.User) {
		c.Status(http.StatusUnauthorized)
	} else {
		loginUser(c, &body)
	}
}

func loginUser(c *gin.Context, login *LoginBody) {
	user, err := GetUser(login.User, login.Password)

	if user == nil || err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaim{
		User: login.User,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(Env().JWTExpires)),
		},
	})

	if tokenString, err := token.SignedString(Env().JWTSecret); err != nil {
		c.Status(http.StatusInternalServerError)
	} else {
		c.Header("Token", tokenString)
		c.Status(http.StatusOK)
	}
}
