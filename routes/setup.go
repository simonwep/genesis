package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {

	// Create router
	router := gin.New()

	// Auth endpoints
	router.POST("/login", Login)
	router.POST("/register", Register)

	// Data endpoints
	router.PUT("/data/:key", SetData)
	router.GET("/data", Data)

	return router
}
