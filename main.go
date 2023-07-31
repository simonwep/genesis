package main

import (
	"github.com/simonwep/genesis/core"
	"github.com/simonwep/genesis/routes"
)

func main() {
	router := routes.SetupRoutes()
	router.SetTrustedProxies(nil)
	router.Run("0.0.0.0:" + core.Config.AppPort)
}
