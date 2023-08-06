package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/simonwep/genesis/core"
	"go.uber.org/zap"
	"net/http"
)

func CreateUser(c *gin.Context) {
	var body core.User

	if !isAsAdminAuthenticated(c) {
		c.Status(http.StatusUnauthorized)
	} else if err := c.ShouldBindJSON(&body); err != nil {
		c.Status(http.StatusBadRequest)
	} else if !core.Config.AppUserPattern.MatchString(body.User) {
		c.Status(http.StatusBadRequest)
	} else if err := core.UpsertUser(body, false); err != nil {
		c.Status(http.StatusInternalServerError)
		core.Logger.Error("failed to create user", zap.Error(err))
	} else {
		c.Status(http.StatusCreated)
	}
}

func UpdateUser(c *gin.Context) {
	var body core.User

	if !isAsAdminAuthenticated(c) {
		c.Status(http.StatusUnauthorized)
	} else if err := c.ShouldBindJSON(&body); err != nil {
		c.Status(http.StatusBadRequest)
	} else if err := core.UpsertUser(body, true); err != nil {
		c.Status(http.StatusInternalServerError)
		core.Logger.Error("failed to create user", zap.Error(err))
	} else {
		c.Status(http.StatusOK)
	}
}

func DeleteUser(c *gin.Context) {
	name := c.Param("name")

	if !isAsAdminAuthenticated(c) {
		c.Status(http.StatusUnauthorized)
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
		c.Status(http.StatusUnauthorized)
	} else if list, err := core.GetUsers(); err != nil {
		c.Status(http.StatusInternalServerError)
		core.Logger.Error("failed to retrieve users", zap.Error(err))
	} else {
		c.JSON(http.StatusOK, list)
	}

}

func isAsAdminAuthenticated(c *gin.Context) bool {
	user := AuthenticateUser(c)
	return user != nil && user.Admin
}
