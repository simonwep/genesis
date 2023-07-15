package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/simonwep/genisis/core"
	"github.com/simonwep/genisis/middleware"
)

func SetupRoutes() *gin.Engine {

	// Set mode
	gin.SetMode(core.Config.GinMode)

	// Create router
	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())

	// Auth endpoint
	router.POST("/login", Login)

	// Data endpoints
	router.PUT("/data/:key", middleware.LimitBodySize(core.Config.AppValueMaxSize), middleware.MinifyJson(), SetData)
	router.DELETE("/data/:key", DeleteData)
	router.GET("/data/:key", DataByKey)
	router.GET("/data", Data)

	return router
}
