package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginBody struct {
	User     string `json:"user"`
	Password string `json:"password"`
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
	user, err := AuthenticateUser(login.User, login.Password)

	if user == nil || err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	tokenString, err := CreateAuthToken(user)

	if err != nil {
		c.Status(http.StatusInternalServerError)
	} else {
		c.Header("Authorization", "Bearer "+tokenString)
		c.Status(http.StatusOK)
	}
}
