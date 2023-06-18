package routes

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"
	"github.com/simonwep/genisis/core"
	"net/http"
	"strings"
)

func Data(c *gin.Context) {
	user := getUser(c)

	if user == nil {
		c.Status(http.StatusUnauthorized)
	} else if data, err := core.GetAllDataFromUser(user.User); err != nil {
		c.Status(http.StatusInternalServerError)
	} else {
		c.Data(http.StatusOK, "application/json; charset=utf-8", data)
	}
}

func DataByKey(c *gin.Context) {
	key := c.Param("key")
	user := getUser(c)

	if user == nil {
		c.Status(http.StatusUnauthorized)
	} else if !core.Config().AppAllowedKeyPattern.MatchString(key) {
		c.Status(http.StatusBadRequest)
	} else if data, err := core.GetDataFromUser(user.User, key); err != nil {
		if err == badger.ErrKeyNotFound {
			c.Status(http.StatusNotFound)
		} else {
			c.Status(http.StatusInternalServerError)
		}
	} else {
		c.Data(http.StatusOK, "application/json; charset=utf-8", data)
	}
}

func SetData(c *gin.Context) {
	var body map[string]interface{}
	key := c.Param("key")
	user := getUser(c)

	if user == nil {
		c.Status(http.StatusUnauthorized)
	} else if !core.Config().AppAllowedKeyPattern.MatchString(key) {
		c.Status(http.StatusBadRequest)
	} else if err := c.BindJSON(&body); err != nil {
		c.Status(http.StatusBadRequest)
	} else if err := core.SetDataForUser(user.User, key, body); err != nil {
		c.Status(http.StatusInternalServerError)
	} else {
		c.Status(http.StatusOK)
	}
}

func getUser(c *gin.Context) *core.User {
	token := strings.Split(c.GetHeader("Authorization"), "Bearer ")[1]

	if parsed, err := core.ParseAuthToken(token); err != nil {
		return nil
	} else if user, err := core.GetUser(parsed.User); err != nil {
		return nil
	} else {
		return user
	}
}
