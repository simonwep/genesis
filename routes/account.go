package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/simonwep/genesis/core"
	"net/http"
)

type updateBody struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword" validate:"required,gte=8,lte=64"`
}

// UpdateAccount godoc
// @Summary      Update account password
// @Description  Update the password for the currently authenticated user
// @Tags         account
// @Accept       json
// @Produce      json
// @Param        request body UpdatePasswordRequest true "Password update request"
// @Success      200 "Password updated successfully"
// @Failure      400 {object} ErrorResponse "Invalid JSON or validation failed"
// @Failure      401 {object} ErrorResponse "Unauthorized or current password incorrect"
// @Security     CookieAuth
// @Router       /account/update [post]
func UpdateAccount(c *gin.Context) {
	validate := validator.New()
	user := authenticateUser(c)

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var body updateBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	} else if _, err := core.AuthenticateUser(user.Name, body.CurrentPassword); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "current password incorrect"})
		return
	}

	if err := validate.Struct(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed, must contain currentPassword and newPassword"})
	} else if err := core.UpdateUser(user.Name, core.PartialUser{
		Admin:    nil,
		Password: &body.NewPassword,
	}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to update user"})
	} else {
		c.Status(http.StatusOK)
	}
}
