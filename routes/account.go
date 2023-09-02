package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/simonwep/genesis/core"
	"net/http"
)

type updateBody struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

func UpdateAccount(c *gin.Context) {
	validate := validator.New()
	user := authenticateUser(c)

	if user == nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	var body updateBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.Status(http.StatusBadRequest)
		return
	} else if _, err := core.AuthenticateUser(user.Name, body.CurrentPassword); err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	newUser := core.User{
		Name:     user.Name,
		Admin:    user.Admin,
		Password: body.NewPassword,
	}

	if err := validate.Struct(&newUser); err != nil {
		c.Status(http.StatusBadRequest)
	} else if err := core.UpsertUser(newUser, true); err != nil {
		c.Status(http.StatusBadRequest)
	} else {
		c.Status(http.StatusOK)
	}
}
