package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/simonwep/genisis/core"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type LoginBody struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
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

	if _, err := core.GetUser(body.User); err != nil {
		c.Status(http.StatusInternalServerError)
		core.Logger.Error("failed to check for user", zap.Error(err))
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
		c.JSON(http.StatusOK, LoginResponse{
			Token:     tokenString,
			ExpiresAt: time.Now().UnixMilli() + core.Config.JWTExpires.Milliseconds(),
		})
	}
}

func validateUserName(name string) bool {
	users := core.Config.AppInitialUsers

	for _, user := range users {
		if user.Name == name {
			return true
		}
	}

	return false
}
