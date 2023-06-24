package main

import (
	"github.com/simonwep/genisis/routes"
)

func main() {
	router := routes.SetupRoutes()
	router.SetTrustedProxies(nil)
	router.Run("localhost:8080")
}
