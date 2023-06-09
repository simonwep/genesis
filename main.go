package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func GetEnv(key string) string {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {

	// Set mode
	gin.SetMode(GetEnv("GIN_MODE"))

	// Create router
	router := gin.New()

	// Bind endpoints
	router.POST("/register", Register)

	// Configure and start server
	router.SetTrustedProxies(nil)
	router.Run("localhost:8080")
}
