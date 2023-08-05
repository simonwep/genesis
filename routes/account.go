package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/simonwep/genesis/core"
	"go.uber.org/zap"
	"net/http"
)

type updateBody struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

func UpdateAccount(c *gin.Context) {
	user := AuthenticateUser(c)

	if user == nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	// TODO: use upsert
	var body updateBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Status(http.StatusBadRequest)
	} else if _, err := core.AuthenticateUser(user.User, body.CurrentPassword); err != nil {
		c.Status(http.StatusBadRequest)
	} else if err := core.SetPasswordForUser(user.User, body.NewPassword); err != nil {
		c.Status(http.StatusInternalServerError)
		core.Logger.Error("failed to update user password", zap.Error(err))
	} else {
		c.Status(http.StatusOK)
	}
}
