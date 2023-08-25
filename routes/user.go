package routes

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/simonwep/genesis/core"
	"go.uber.org/zap"
	"net/http"
)

func CreateUser(c *gin.Context) {
	var body core.User

	if !isAsAdminAuthenticated(c) {
		c.Status(http.StatusForbidden)
	} else if err := c.ShouldBindJSON(&body); err != nil {
		c.Status(http.StatusBadRequest)
	} else if !core.Config.AppUserPattern.MatchString(body.Name) {
		c.Status(http.StatusBadRequest)
	} else if err := core.UpsertUser(body, false); err != nil {
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
	name := c.Param("name")
	var body core.User

	if !isAsAdminAuthenticated(c) {
		c.Status(http.StatusForbidden)
	} else if err := c.ShouldBindJSON(&body); err != nil {
		c.Status(http.StatusBadRequest)
	} else if usr, err := core.GetUser(body.Name); err != nil {
		c.Status(http.StatusInternalServerError)
		core.Logger.Error("failed to retrieve user", zap.Error(err))
	} else if name != body.Name && usr != nil {
		c.Status(http.StatusConflict)
	} else if err := core.UpsertUser(body, true); err != nil {
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
	if !isAsAdminAuthenticated(c) {
		c.Status(http.StatusForbidden)
	} else if list, err := core.GetUsers(); err != nil {
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
