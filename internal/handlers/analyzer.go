package handlers

import (
	"github.com/PratikforCoding/CodeSentry/internal/models"
	"github.com/PratikforCoding/CodeSentry/internal/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type AnalyzerHandler struct {
	AnalyzerService services.AnalyzerServiceInterface
}

func NewAnalyzerHandler(db *mongo.Database) *AnalyzerHandler {
	return &AnalyzerHandler{
		AnalyzerService: services.NewAnalyzerService(db),
	}
}

func (ah *AnalyzerHandler) AnalyzeCode(c *gin.Context) {
	var req models.AnalyzeRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Code:    400,
			Message: "Invalid request",
		})
		return
	}

	if !req.Options.CheckSecurity && !req.Options.CheckStyle && !req.Options.CheckComplexity && !req.Options.CheckMetrics {
		req.Options.CheckSecurity = true
		req.Options.CheckStyle = true
		req.Options.CheckComplexity = true
		req.Options.CheckMetrics = true
	}

	response := ah.AnalyzerService.AnalyzeCode(req)
	c.JSON(http.StatusOK, response)
}

func (ah *AnalyzerHandler) AnalyzeComplexity(c *gin.Context) {
	var req models.AnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Code:    400,
			Message: err.Error(),
		})
		return
	}
	req.Options.CheckComplexity = true
	//req.Options.CheckSecurity = false
	//req.Options.CheckStyle = false
	//req.Options.CheckMetrics = false

	response := ah.AnalyzerService.AnalyzeCode(req)
	c.JSON(http.StatusOK, gin.H{
		"language":         response.Language,
		"complexity_score": response.ComplexityScore,
		"overall_score":    response.OverallScore,
	})
}

func (ah *AnalyzerHandler) AnalyzeSecurity(c *gin.Context) {
	var req models.AnalyzeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	req.Options.CheckSecurity = true
	req.Options.CheckComplexity = false
	req.Options.CheckStyle = false
	req.Options.CheckMetrics = false

	response := ah.AnalyzerService.AnalyzeCode(req)
	c.JSON(http.StatusOK, gin.H{
		"language":        response.Language,
		"security_issues": response.SecurityIssues,
		"overall_score":   response.OverallScore,
	})
}

func (ah *AnalyzerHandler) AnalyzeStyle(c *gin.Context) {
	var req models.AnalyzeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request",
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	req.Options.CheckStyle = true
	req.Options.CheckComplexity = false
	req.Options.CheckSecurity = false
	req.Options.CheckMetrics = false

	response := ah.AnalyzerService.AnalyzeCode(req)
	c.JSON(http.StatusOK, gin.H{
		"language":          response.Language,
		"style_suggestions": response.StyleSuggestions,
		"overall_score":     response.OverallScore,
	})
}

func (ah *AnalyzerHandler) GetSupportedLanguages(c *gin.Context) {
	languages := []string{"go", "javascript", "python", "java", "sql"}
	c.JSON(http.StatusOK, gin.H{
		"supported_languages": languages,
	})
}
