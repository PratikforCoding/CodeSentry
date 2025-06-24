package handlers

import (
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"

	"github.com/PratikforCoding/CodeSentry/internal/models"
	"github.com/PratikforCoding/CodeSentry/internal/repository"
	"github.com/gin-gonic/gin"
)

type AnalysesHandler struct {
	analysisRepo repository.AnalysisRepositoryInterface
}

func NewAnalysesHandlerWithRepo(repo repository.AnalysisRepositoryInterface) *AnalysesHandler {
	return &AnalysesHandler{
		analysisRepo: repo,
	}
}

func NewAnalysesHandler(db *mongo.Database) *AnalysesHandler {
	return &AnalysesHandler{
		analysisRepo: repository.NewAnalysisRepository(db),
	}
}

func (ah *AnalysesHandler) GetAnalyses(c *gin.Context) {

	language := c.Query("language")

	analyses, err := ah.analysisRepo.GetAllAnalyses(language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to fetch analyses",
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"analyses": analyses,
		"count":    len(analyses),
	})
}

func (ah *AnalysesHandler) GetAnalysis(c *gin.Context) {
	id := c.Param("id")

	analysis, err := ah.analysisRepo.GetAnalysisByID(id)
	if err != nil {
		if err.Error() == "analysis not found" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Analysis not found",
				Code:    404,
				Message: "Analysis with the given ID does not exist",
			})
			return
		}
		if err.Error() == "invalid ID format" {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "Invalid ID format",
				Code:    400,
				Message: "The provided ID is not in valid format",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to fetch analysis",
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, analysis)
}

func (ah *AnalysesHandler) UpdateAnalysis(c *gin.Context) {
	id := c.Param("id")

	var updateReq models.UpdateAnalysisRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Code:    400,
			Message: err.Error(),
		})
		return
	}

	err := ah.analysisRepo.UpdateAnalysis(id, updateReq)
	if err != nil {
		if err.Error() == "analysis not found" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Analysis not found",
				Code:    404,
				Message: "Analysis with the given ID does not exist",
			})
			return
		}
		if err.Error() == "invalid ID format" {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "Invalid ID format",
				Code:    400,
				Message: "The provided ID is not in valid format",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to update analysis",
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Analysis updated successfully",
	})
}

func (ah *AnalysesHandler) DeleteAnalysis(c *gin.Context) {
	id := c.Param("id")

	err := ah.analysisRepo.DeleteAnalysis(id)
	if err != nil {
		if err.Error() == "analysis not found" {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error:   "Analysis not found",
				Code:    404,
				Message: "Analysis with the given ID does not exist",
			})
			return
		}
		if err.Error() == "invalid ID format" {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "Invalid ID format",
				Code:    400,
				Message: "The provided ID is not in valid format",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to delete analysis",
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Analysis deleted successfully",
	})
}
