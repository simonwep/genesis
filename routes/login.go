package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/simonwep/genisis/core"
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
	} else {
		loginUser(c, &body)
	}
}

func loginUser(c *gin.Context, login *LoginBody) {
	user, err := core.AuthenticateUser(login.User, login.Password)

	if user == nil || err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	tokenString, err := core.CreateAuthToken(user)

	if err != nil {
		c.Status(http.StatusInternalServerError)
	} else {
		c.Header("Authorization", "Bearer "+tokenString)
		c.Status(http.StatusOK)
	}
}
