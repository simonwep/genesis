package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/simonwep/genesis/core"
	"github.com/simonwep/genesis/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Genesis API
// @version         1.0
// @description     A generic JSON API for small, private frontend apps with user management and data storage
// @description     Authentication uses HTTP-only cookies with JWT tokens (cookie name: 'gt')
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    https://github.com/simonwep/genesis

// @license.name  MIT
// @license.url   https://github.com/simonwep/genesis/blob/main/LICENSE

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name gt
// @description JWT token stored in HTTP-only cookie

func SetupRoutes() *gin.Engine {

	// Set mode
	gin.SetMode(core.Config.AppGinMode)

	// Create router
	root := gin.New()

	// Middleware
	root.Use(gin.Recovery())

	// Wrap routes under common path
	router := root.Group(core.Config.BaseUrl)

	// Auth and account endpoints
	router.POST("/login", Login)
	router.POST("/account/update", UpdateAccount)
	router.POST("/logout", Logout)

	// User endpoints
	router.GET("/user", GetUser)
	router.POST("/user", CreateUser)
	router.POST("/user/:name", UpdateUser)
	router.DELETE("/user/:name", DeleteUser)

	// Data endpoints
	router.POST("/data/:key", middleware.LimitBodySize(core.Config.AppDataMaxSize), middleware.MinifyJson(), SetData)
	router.DELETE("/data/:key", DeleteData)
	router.GET("/data/:key", DataByKey)
	router.GET("/data", Data)

	// Heal check endpoints
	router.GET("/health", Health)

	// Swagger documentation
	if core.Config.SwaggerEnabled {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return root
}
