package handlers_test

import (
	"bytes"
	"encoding/json"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PratikforCoding/CodeSentry/internal/handlers"
	"github.com/PratikforCoding/CodeSentry/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockAnalyzerService mocks services.AnalyzerService interface
type MockAnalyzerService struct {
	AnalyzeCodeFunc func(req models.AnalyzeRequest) models.AnalysisResponse
}

func (m *MockAnalyzerService) AnalyzeCode(req models.AnalyzeRequest) models.AnalysisResponse {
	return m.AnalyzeCodeFunc(req)
}

func setupRouterWithMockService(mockService *MockAnalyzerService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := &handlers.AnalyzerHandler{
		AnalyzerService: mockService,
	}

	router.POST("/analyze", handler.AnalyzeCode)
	router.POST("/analyze/complexity", handler.AnalyzeComplexity)
	router.POST("/analyze/security", handler.AnalyzeSecurity)
	router.POST("/analyze/style", handler.AnalyzeStyle)
	router.GET("/languages", handler.GetSupportedLanguages)

	return router
}

func TestAnalyzerHandler_AnalyzeCode(t *testing.T) {
	mockService := &MockAnalyzerService{
		AnalyzeCodeFunc: func(req models.AnalyzeRequest) models.AnalysisResponse {
			return models.AnalysisResponse{Language: "go", OverallScore: 90.5}
		},
	}

	router := setupRouterWithMockService(mockService)

	reqBody := models.AnalyzeRequest{
		Code: "package main\nfunc main() {}",
		Options: struct {
			CheckSecurity   bool `json:"check_security"`
			CheckStyle      bool `json:"check_style"`
			CheckComplexity bool `json:"check_complexity"`
			CheckMetrics    bool `json:"check_metrics"`
		}{
			CheckSecurity:   false,
			CheckStyle:      false,
			CheckComplexity: false,
			CheckMetrics:    false,
		},
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/analyze", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.AnalysisResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "go", resp.Language)
	assert.Equal(t, 90.5, resp.OverallScore)
}

func TestAnalyzerHandler_AnalyzeCode_InvalidJSON(t *testing.T) {
	mockService := &MockAnalyzerService{}

	router := setupRouterWithMockService(mockService)

	req, _ := http.NewRequest("POST", "/analyze", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAnalyzerHandler_AnalyzeComplexity(t *testing.T) {
	mockService := &MockAnalyzerService{
		AnalyzeCodeFunc: func(req models.AnalyzeRequest) models.AnalysisResponse {
			return models.AnalysisResponse{
				Language:        "go",
				ComplexityScore: 5,
				OverallScore:    75.0,
			}
		},
	}

	router := setupRouterWithMockService(mockService)

	reqBody := models.AnalyzeRequest{
		Code: "func main() {}",
	}
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/analyze/complexity", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "go", resp["language"])
	assert.Equal(t, float64(5), resp["complexity_score"])
	assert.Equal(t, float64(75.0), resp["overall_score"])
}

func TestAnalyzerHandler_GetSupportedLanguages(t *testing.T) {
	mockService := &MockAnalyzerService{}
	router := setupRouterWithMockService(mockService)

	req, _ := http.NewRequest("GET", "/languages", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string][]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp["supported_languages"], "go")
	assert.Contains(t, resp["supported_languages"], "python")
}
