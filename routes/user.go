package routes

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/simonwep/genesis/core"
	"go.uber.org/zap"
	"net/http"
)

func CreateUser(c *gin.Context) {
	validate := validator.New()
	var body core.User

	if !isAsAdminAuthenticated(c) {
		c.Status(http.StatusForbidden)
	} else if err := c.ShouldBindJSON(&body); err != nil {
		c.Status(http.StatusBadRequest)
	} else if !core.Config.AppUserPattern.MatchString(body.Name) {
		c.Status(http.StatusBadRequest)
	} else if err := validate.Struct(&body); err != nil {
		c.Status(http.StatusBadRequest)
	} else if err := core.CreateUser(body); err != nil {
		if errors.Is(err, core.ErrUserAlreadyExists) {
			c.Status(http.StatusConflict)
		} else {
			c.Status(http.StatusInternalServerError)
			core.Logger.Error("failed to create user", zap.Error(err))
		}
	} else {
		c.Status(http.StatusCreated)
	}
}

func UpdateUser(c *gin.Context) {
	user := authenticateUser(c)
	validate := validator.New()
	name := c.Param("name")
	var body core.PartialUser

	if user == nil || !user.Admin {
		c.Status(http.StatusForbidden)
	} else if name == user.Name {
		c.Status(http.StatusForbidden)
	} else if err := c.ShouldBindJSON(&body); err != nil {
		c.Status(http.StatusBadRequest)
	} else if err := validate.Struct(&body); err != nil {
		c.Status(http.StatusBadRequest)
	} else if _, err := core.GetUser(name); err != nil {
		c.Status(http.StatusInternalServerError)
		core.Logger.Error("failed to retrieve user", zap.Error(err))
	} else if err := core.UpdateUser(name, body); err != nil {
		c.Status(http.StatusBadRequest)
	} else {
		c.Status(http.StatusOK)
	}
}

func DeleteUser(c *gin.Context) {
	name := c.Param("name")

	if !isAsAdminAuthenticated(c) {
		c.Status(http.StatusForbidden)
	} else {
		if err := core.DeleteUser(name); err != nil {
			c.Status(http.StatusInternalServerError)
			core.Logger.Error("Failed to delete user", zap.String("name", name), zap.Error(err))
		} else {
			c.Status(http.StatusOK)
		}
	}
}

func GetUser(c *gin.Context) {
	user := authenticateUser(c)

	if user == nil || !user.Admin {
		c.Status(http.StatusForbidden)
	} else if list, err := core.GetUsers(user.Name); err != nil {
		c.Status(http.StatusInternalServerError)
		core.Logger.Error("failed to retrieve users", zap.Error(err))
	} else {
		c.JSON(http.StatusOK, list)
	}
}

func isAsAdminAuthenticated(c *gin.Context) bool {
	user := authenticateUser(c)
	return user != nil && user.Admin
}
