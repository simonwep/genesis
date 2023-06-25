package main

import (
	"github.com/simonwep/genisis/core"
	"github.com/simonwep/genisis/routes"
)

func main() {
	router := routes.SetupRoutes()
	router.SetTrustedProxies(nil)
	router.Run("0.0.0.0:" + core.Config.AppPort)
}
