package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/simonwep/genisis/core"
	"github.com/simonwep/genisis/middleware"
)

func SetupRoutes() *gin.Engine {

	// Set mode
	gin.SetMode(core.Config.AppGinMode)

	// Create router
	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())

	// Auth and account endpoints
	router.POST("/login", Login)
	router.POST("/login/refresh", Refresh)
	router.POST("/account/update", Update)
	router.POST("/logout", Logout)

	// Data endpoints
	router.POST("/data/:key", middleware.LimitBodySize(core.Config.AppValueMaxSize), middleware.MinifyJson(), SetData)
	router.DELETE("/data/:key", DeleteData)
	router.GET("/data/:key", DataByKey)
	router.GET("/data", Data)

	return router
}
