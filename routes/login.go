package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/simonwep/genisis/core"
	"go.uber.org/zap"
	"net/http"
)

type LoginBody struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var body LoginBody

	if err := c.ShouldBindJSON(&body); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	// quick path if username is invalid
	if !validateUserName(body.User) {
		c.Status(http.StatusUnauthorized)
		return
	}

	user, err := core.AuthenticateUser(body.User, body.Password)
	if user == nil || err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	tokenString, err := core.CreateAuthToken(user)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		core.Logger.Error("failed to create auth token", zap.Error(err))
	} else {
		c.Header("Authorization", "Bearer "+tokenString)
		c.Status(http.StatusOK)
	}
}
