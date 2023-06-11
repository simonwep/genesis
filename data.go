package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Data(c *gin.Context) {
	user := getUser(c)

	if user == nil {
		c.Status(http.StatusUnauthorized)
	} else if data, err := GetAllDataFromUser(user.User); err != nil {
		c.Status(http.StatusInternalServerError)
	} else {
		c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(data))
	}
}

func SetData(c *gin.Context) {
	var body map[string]interface{}
	key := c.Param("key")
	user := getUser(c)

	if user == nil {
		c.Status(http.StatusUnauthorized)
	} else if err := c.BindJSON(&body); err != nil {
		c.Status(http.StatusBadRequest)
	} else if err := SetDataForUser(user.User, key, body); err != nil {
		c.Status(http.StatusInternalServerError)
	} else {
		c.Status(http.StatusOK)
	}
}

func getUser(c *gin.Context) *User {
	token := strings.Split(c.GetHeader("Authorization"), "Bearer ")[1]

	if parsed, err := ParseAuthToken(token); err != nil {
		return nil
	} else if user, err := GetUser(parsed.User); err != nil {
		return nil
	} else {
		return user
	}
}
