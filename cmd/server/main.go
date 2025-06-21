package main

import (
	"github.com/PratikforCoding/CodeSentry/api"
	"github.com/PratikforCoding/CodeSentry/internal/database"
	"github.com/PratikforCoding/CodeSentry/pkg/config"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	cfg := config.Load()
	database.InitMongoDB(cfg)
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := api.SetupRoutes()
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	log.Print("Starting server on port: " + port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to run server", err)
	}
}
