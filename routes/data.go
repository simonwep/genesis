package routes

import (
	"errors"
	"github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"
	"github.com/simonwep/genesis/core"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func Data(c *gin.Context) {
	user := authenticateUser(c)

	if user == nil {
		c.Status(http.StatusUnauthorized)
	} else if data, err := core.GetAllDataFromUser(user.User); err != nil {
		c.Status(http.StatusInternalServerError)
		core.Logger.Error("failed to retrieve data", zap.Error(err))
	} else {
		c.Data(http.StatusOK, "application/json; charset=utf-8", data)
	}
}

func DataByKey(c *gin.Context) {
	key := c.Param("key")
	user := authenticateUser(c)

	if user == nil {
		c.Status(http.StatusUnauthorized)
	} else if !core.Config.AppKeyPattern.MatchString(key) {
		c.Status(http.StatusNotFound)
	} else if data, err := core.GetDataFromUser(user.User, key); err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			c.Status(http.StatusNoContent)
		} else {
			c.Status(http.StatusInternalServerError)
			core.Logger.Error("failed to retrieve unit of data", zap.Error(err))
		}
	} else {
		c.Data(http.StatusOK, "application/json; charset=utf-8", data)
	}
}

func SetData(c *gin.Context) {
	key := c.Param("key")
	user := authenticateUser(c)

	if user == nil {
		c.Status(http.StatusUnauthorized)
	} else if !core.Config.AppKeyPattern.MatchString(key) {
		c.Status(http.StatusBadRequest)
	} else if count := core.GetDataCountForUser(user.User, key); count > core.Config.AppKeysPerUser {
		c.Status(http.StatusForbidden)
	} else if size, err := getContentLength(c); err != nil || size > core.Config.AppDataMaxSize {
		c.Status(http.StatusRequestEntityTooLarge)
	} else if body, err := c.GetRawData(); err != nil {
		c.Status(http.StatusBadRequest)
	} else if err := core.SetDataForUser(user.User, key, body); err != nil {
		c.Status(http.StatusInternalServerError)
		core.Logger.Error("failed to set data", zap.Error(err))
	} else {
		c.Status(http.StatusOK)
	}
}

func DeleteData(c *gin.Context) {
	key := c.Param("key")
	user := authenticateUser(c)

	if user == nil {
		c.Status(http.StatusUnauthorized)
	} else if err := core.DeleteDataFromUser(user.User, key); err != nil {
		c.Status(http.StatusInternalServerError)
		core.Logger.Error("failed to delete data", zap.Error(err))
	} else {
		c.Status(http.StatusOK)
	}
}

func getContentLength(c *gin.Context) (int64, error) {
	return strconv.ParseInt(c.GetHeader("Content-Length"), 10, 64)
}
