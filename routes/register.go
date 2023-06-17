package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/simonwep/genisis/core"
	"golang.org/x/exp/slices"
	"log"
	"net/http"
)

type APIError struct {
	Err string `json:"error"`
}

func Register(c *gin.Context) {
	var body LoginBody

	if err := c.BindJSON(&body); err != nil {
		c.Status(http.StatusBadRequest)
		log.Printf("Received invalid body: %v", err)
	} else if !validateUserName(body.User) {
		c.Status(http.StatusUnauthorized)
	} else if err := core.CreateUser(body.User, body.Password); err != nil {
		c.Status(http.StatusUnauthorized)
		log.Printf("Failed to create user: %v", err)
	} else {
		c.Status(http.StatusCreated)
	}
}

func validateUserName(name string) bool {
	users := core.Config().AppAllowedUsers
	return len(users) == 0 || slices.Contains(users, name)
}
