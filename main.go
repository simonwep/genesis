package main

import (
	"github.com/gin-gonic/gin"
	"github.com/simonwep/genisis/core"
	"github.com/simonwep/genisis/routes"
)

func main() {

	// Set mode
	gin.SetMode(core.Config().GinMode)

	router := routes.SetupRoutes()

	// Configure and start server
	router.SetTrustedProxies(nil)
	router.Run("localhost:8080")
}
