package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/simonwep/genisis/core"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"net/http"
)

type LoginBody struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var body LoginBody
	var userCreated bool

	if err := c.ShouldBindJSON(&body); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	// quick path if username is invalid
	if !validateUserName(body.User) {
		c.Status(http.StatusUnauthorized)
		return
	}

	if user, err := core.GetUser(body.User); err != nil {
		c.Status(http.StatusInternalServerError)
		core.Logger.Error("failed to check for user", zap.Error(err))
		return
	} else if user == nil {
		if err := core.CreateUser(body.User, body.Password); err != nil {
			c.Status(http.StatusUnauthorized)
			core.Logger.Error("failed to register user", zap.Error(err))
			return
		} else {
			userCreated = true
		}
	} else {
		userCreated = false
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

		if userCreated {
			c.Status(http.StatusCreated)
		} else {
			c.Status(http.StatusOK)
		}
	}
}

func validateUserName(name string) bool {
	users := core.Config.AppAllowedUsers
	return len(users) == 0 || slices.Contains(users, name)
}
