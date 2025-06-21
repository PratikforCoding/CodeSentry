package api

import (
	"github.com/PratikforCoding/CodeSentry/internal/handlers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/health", handlers.HealthCheck)

	analyzerHandler := handlers.NewAnalyzerHandler()
	analysesHandler := handlers.NewAnalysesHandler()

	v1 := router.Group("/api/v1")
	{
		// Existing analysis routes
		v1.POST("/analyze", analyzerHandler.AnalyzeCode)
		v1.POST("/analyze/complexity", analyzerHandler.AnalyzeComplexity)
		v1.POST("/analyze/security", analyzerHandler.AnalyzeSecurity)
		v1.POST("/analyze/style", analyzerHandler.AnalyzeStyle)
		v1.GET("/languages", analyzerHandler.GetSupportedLanguages)

		// New CRUD routes for analyses
		v1.GET("/analyses", analysesHandler.GetAnalyses)
		v1.GET("/analyses/:id", analysesHandler.GetAnalysis)
		v1.PUT("/analyses/:id", analysesHandler.UpdateAnalysis)
		v1.DELETE("/analyses/:id", analysesHandler.DeleteAnalysis)
	}
	return router
}
