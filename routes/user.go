package routes

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/simonwep/genesis/core"
	"go.uber.org/zap"
	"net/http"
)

// CreateUser godoc
// @Summary      Create a new user
// @Description  Create a new user (admin only)
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user body CreateUserRequest true "User details"
// @Success      201 {object} SuccessResponse "User created successfully"
// @Failure      400 {object} ErrorResponse "Invalid JSON or validation failed"
// @Failure      403 {object} ErrorResponse "Forbidden - admin only"
// @Failure      409 {object} ErrorResponse "User already exists"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Security     CookieAuth
// @Router       /user [post]
func CreateUser(c *gin.Context) {
	validate := validator.New()
	var body core.User

	if !isAsAdminAuthenticated(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "only admins can create users"})
	} else if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
	} else if !core.Config.AppUserPattern.MatchString(body.Name) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user name, must match " + core.Config.AppUserPattern.String()})
	} else if err := validate.Struct(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation of json failed, must contain name, password and admin"})
	} else if err := core.CreateUser(body); err != nil {
		if errors.Is(err, core.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			core.Logger.Error("failed to create user", zap.Error(err))
		}
	} else {
		c.JSON(http.StatusCreated, gin.H{"message": "user created"})
	}
}

// UpdateUser godoc
// @Summary      Update a user
// @Description  Update user details by name (admin only, cannot update self)
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        name path string true "Username"
// @Param        user body UpdateUserRequest true "User update details"
// @Success      200 "User updated successfully"
// @Failure      400 {object} ErrorResponse "Invalid JSON or validation failed"
// @Failure      403 {object} ErrorResponse "Forbidden - admin only or cannot update self"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Security     CookieAuth
// @Router       /user/{name} [post]
func UpdateUser(c *gin.Context) {
	user := authenticateUser(c)
	validate := validator.New()
	name := c.Param("name")
	var body core.PartialUser

	if user == nil || !user.Admin {
		c.JSON(http.StatusForbidden, gin.H{"error": "user not found or you are not an admin"})
	} else if name == user.Name {
		c.JSON(http.StatusForbidden, gin.H{"error": "you cannot update yourself"})
	} else if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
	} else if err := validate.Struct(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation of json failed, may contain admin or password"})
	} else if _, err := core.GetUser(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user"})
		core.Logger.Error("failed to retrieve user", zap.Error(err))
	} else if err := core.UpdateUser(name, body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "update failed"})
	} else {
		c.Status(http.StatusOK)
	}
}

// DeleteUser godoc
// @Summary      Delete a user
// @Description  Delete user by name (admin only)
// @Tags         user
// @Produce      json
// @Param        name path string true "Username"
// @Success      200 "User deleted successfully"
// @Failure      403 {object} ErrorResponse "Forbidden - admin only"
// @Failure      500 {object} ErrorResponse "Failed to delete user"
// @Security     CookieAuth
// @Router       /user/{name} [delete]
func DeleteUser(c *gin.Context) {
	name := c.Param("name")

	if !isAsAdminAuthenticated(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	} else {
		if err := core.DeleteUser(name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
			core.Logger.Error("Failed to delete user", zap.String("name", name), zap.Error(err))
		} else {
			c.Status(http.StatusOK)
		}
	}
}

// GetUser godoc
// @Summary      Get all users
// @Description  List all users (admin only)
// @Tags         user
// @Produce      json
// @Success      200 {array} core.PublicUser "List of users"
// @Failure      403 {object} ErrorResponse "Forbidden - admin only"
// @Failure      500 {object} ErrorResponse "Failed to retrieve users"
// @Security     CookieAuth
// @Router       /user [get]
func GetUser(c *gin.Context) {
	user := authenticateUser(c)

	if user == nil || !user.Admin {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	} else if list, err := core.GetUsers(user.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve users"})
		core.Logger.Error("failed to retrieve users", zap.Error(err))
	} else {
		c.JSON(http.StatusOK, list)
	}
}

func isAsAdminAuthenticated(c *gin.Context) bool {
	user := authenticateUser(c)
	return user != nil && user.Admin
}
