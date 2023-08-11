package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/simonwep/genesis/core"
	"net/http"
)

type updateBody struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

func UpdateAccount(c *gin.Context) {
	user := authenticateUser(c)

	if user == nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	var body updateBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Status(http.StatusBadRequest)
	} else if _, err := core.AuthenticateUser(user.User, body.CurrentPassword); err != nil {
		c.Status(http.StatusUnauthorized)
	} else if err := core.UpsertUser(core.User{
		User:     user.User,
		Admin:    user.Admin,
		Password: body.NewPassword,
	}, true); err != nil {
		c.Status(http.StatusBadRequest)
	} else {
		c.Status(http.StatusOK)
	}
}
