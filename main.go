package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	// Set mode
	gin.SetMode(Env().GinMode)

	// Create router
	router := gin.New()

	// Auth endpoints
	router.POST("/login", Login)
	router.POST("/register", Register)

	// Data endpoints
	router.PUT("/data/:key", SetData)
	router.GET("/data", Data)

	// Configure and start server
	router.SetTrustedProxies(nil)
	router.Run("localhost:8080")
}
