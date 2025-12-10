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

// Data godoc
// @Summary      Get all data
// @Description  Retrieve all data for the authenticated user as a JSON object
// @Tags         data
// @Produce      json
// @Success      200 {object} map[string]interface{} "User data as JSON object"
// @Failure      401 {object} ErrorResponse "Unauthorized"
// @Failure      500 {object} ErrorResponse "Failed to retrieve data"
// @Security     CookieAuth
// @Router       /data [get]
func Data(c *gin.Context) {
	user := authenticateUser(c)

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	} else if data, err := core.GetAllDataFromUser(user.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve data"})
		core.Logger.Error("failed to retrieve data", zap.Error(err))
	} else {
		c.Data(http.StatusOK, "application/json; charset=utf-8", data)
	}
}

// DataByKey godoc
// @Summary      Get data by key
// @Description  Retrieve the data stored for a specific key
// @Tags         data
// @Produce      json
// @Param        key path string true "Data key"
// @Success      200 {object} map[string]interface{} "Data for the specified key"
// @Failure      204 "No content found for key"
// @Failure      401 {object} ErrorResponse "Unauthorized"
// @Failure      404 {object} ErrorResponse "Key not found or invalid key pattern"
// @Failure      500 {object} ErrorResponse "Failed to retrieve data"
// @Security     CookieAuth
// @Router       /data/{key} [get]
func DataByKey(c *gin.Context) {
	key := c.Param("key")
	user := authenticateUser(c)

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	} else if !core.Config.AppKeyPattern.MatchString(key) {
		c.JSON(http.StatusNotFound, gin.H{"error": "key must match " + core.Config.AppKeyPattern.String()})
	} else if data, err := core.GetDataFromUser(user.Name, key); err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			c.JSON(http.StatusNoContent, gin.H{"error": "key not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve unit of data"})
			core.Logger.Error("failed to retrieve unit of data", zap.Error(err))
		}
	} else {
		c.Data(http.StatusOK, "application/json; charset=utf-8", data)
	}
}

// SetData godoc
// @Summary      Set data by key
// @Description  Store or update data for a specific key. JSON data is minified and validated.
// @Tags         data
// @Accept       json
// @Produce      json
// @Param        key path string true "Data key"
// @Param        data body map[string]interface{} true "JSON data to store"
// @Success      200 "Data stored successfully"
// @Failure      400 {object} ErrorResponse "Invalid key pattern or invalid body"
// @Failure      401 {object} ErrorResponse "Unauthorized"
// @Failure      403 {object} ErrorResponse "Too many keys (limit exceeded)"
// @Failure      413 {object} ErrorResponse "Request entity too large"
// @Failure      500 {object} ErrorResponse "Failed to set data"
// @Security     CookieAuth
// @Router       /data/{key} [post]
func SetData(c *gin.Context) {
	key := c.Param("key")
	user := authenticateUser(c)

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	} else if !core.Config.AppKeyPattern.MatchString(key) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key must match " + core.Config.AppKeyPattern.String()})
	} else if count := core.GetDataCountForUser(user.Name, key); count > core.Config.AppKeysPerUser {
		c.JSON(http.StatusForbidden, gin.H{"error": "too many keys, limit is " + strconv.FormatInt(core.Config.AppKeysPerUser, 10)})
	} else if size, err := getContentLength(c); err != nil || size > core.Config.AppDataMaxSize {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "request entity too large, limit is " + strconv.FormatInt(core.Config.AppDataMaxSize, 10) + " kilobytes"})
	} else if body, err := c.GetRawData(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
	} else if err := core.SetDataForUser(user.Name, key, body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set data"})
		core.Logger.Error("failed to set data", zap.Error(err))
	} else {
		c.Status(http.StatusOK)
	}
}

// DeleteData godoc
// @Summary      Delete data by key
// @Description  Remove data for a specific key (always returns 200, even if key doesn't exist)
// @Tags         data
// @Produce      json
// @Param        key path string true "Data key"
// @Success      200 "Data deleted successfully"
// @Failure      401 {object} ErrorResponse "Unauthorized"
// @Failure      500 {object} ErrorResponse "Failed to delete data"
// @Security     CookieAuth
// @Router       /data/{key} [delete]
func DeleteData(c *gin.Context) {
	key := c.Param("key")
	user := authenticateUser(c)

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	} else if err := core.DeleteDataFromUser(user.Name, key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete data"})
		core.Logger.Error("failed to delete data", zap.Error(err))
	} else {
		c.Status(http.StatusOK)
	}
}

func getContentLength(c *gin.Context) (int64, error) {
	return strconv.ParseInt(c.GetHeader("Content-Length"), 10, 64)
}
