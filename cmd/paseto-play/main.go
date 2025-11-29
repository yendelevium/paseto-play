package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	pasetotokens "github.com/yendelevium/paseto-play/internal/auth"
	"github.com/yendelevium/paseto-play/internal/config"
	"github.com/yendelevium/paseto-play/internal/routes"
)

// This runs BEFORE main
func init() {
	config.LoadEnv()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	makerPublic, err := pasetotokens.MakePasetoKeyPair()
	if err != nil {
		log.Fatalf("Failed to make key pair for paseto: %v", err)
	}
	makerLocal := pasetotokens.MakePasetoLocalKey()

	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()

	// You CAN nest the routes inside /api doing something like this
	// apiGroup := router.Group("/api")
	// routes.AddRoutes(apiGroup, maker)
	// But directly under / is also fine
	routes.AddRoutes(&r.RouterGroup, makerPublic, makerLocal)

	// Start server on port 8080 (default)
	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
