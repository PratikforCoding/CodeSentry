package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PratikforCoding/CodeSentry/internal/handlers"
	"github.com/PratikforCoding/CodeSentry/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockAnalysisRepository implements repository.AnalysisRepositoryInterface
type MockAnalysisRepository struct {
	SaveAnalysisFunc    func(req models.AnalyzeRequest, res models.AnalysisResponse) error
	GetAllAnalysesFunc  func(language string) ([]models.Analysis, error)
	GetAnalysisByIDFunc func(id string) (models.Analysis, error)
	UpdateAnalysisFunc  func(id string, updateReq models.UpdateAnalysisRequest) error
	DeleteAnalysisFunc  func(id string) error
}

func (m *MockAnalysisRepository) SaveAnalysis(req models.AnalyzeRequest, res models.AnalysisResponse) error {
	if m.SaveAnalysisFunc != nil {
		return m.SaveAnalysisFunc(req, res)
	}
	return nil
}

func (m *MockAnalysisRepository) GetAllAnalyses(language string) ([]models.Analysis, error) {
	return m.GetAllAnalysesFunc(language)
}

func (m *MockAnalysisRepository) GetAnalysisByID(id string) (models.Analysis, error) {
	return m.GetAnalysisByIDFunc(id)
}

func (m *MockAnalysisRepository) UpdateAnalysis(id string, updateReq models.UpdateAnalysisRequest) error {
	return m.UpdateAnalysisFunc(id, updateReq)
}

func (m *MockAnalysisRepository) DeleteAnalysis(id string) error {
	return m.DeleteAnalysisFunc(id)
}

func setupRouterWithMockRepo(mockRepo *MockAnalysisRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	handler := handlers.NewAnalysesHandlerWithRepo(mockRepo)

	router.GET("/analyses", handler.GetAnalyses)
	router.GET("/analyses/:id", handler.GetAnalysis)
	router.PUT("/analyses/:id", handler.UpdateAnalysis)
	router.DELETE("/analyses/:id", handler.DeleteAnalysis)

	return router
}

func TestGetAnalyses_Success(t *testing.T) {
	mockRepo := &MockAnalysisRepository{
		GetAllAnalysesFunc: func(language string) ([]models.Analysis, error) {
			id, err := primitive.ObjectIDFromHex("6858401a0b6bfe01d3921ec7")
			if err != nil {
				t.Fatal(err)
			}
			id2, err := primitive.ObjectIDFromHex("68583ff00b6bfe01d3921ec6")
			if err != nil {
				t.Fatal(err)
			}
			return []models.Analysis{
				{ID: id, Language: language},
				{ID: id2, Language: language},
			}, nil
		},
	}

	router := setupRouterWithMockRepo(mockRepo)

	req, _ := http.NewRequest("GET", "/analyses?language=go", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	// The count is a float64 because JSON numbers unmarshal as float64
	assert.Equal(t, float64(2), resp["count"])
}

func TestGetAnalyses_Error(t *testing.T) {
	mockRepo := &MockAnalysisRepository{
		GetAllAnalysesFunc: func(language string) ([]models.Analysis, error) {
			return nil, errors.New("db error")
		},
	}

	router := setupRouterWithMockRepo(mockRepo)

	req, _ := http.NewRequest("GET", "/analyses", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetAnalysis_Success(t *testing.T) {
	mockRepo := &MockAnalysisRepository{
		GetAnalysisByIDFunc: func(id string) (models.Analysis, error) {
			objectID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				panic(err)
			}
			return models.Analysis{ID: objectID, Language: "go"}, nil
		},
	}

	router := setupRouterWithMockRepo(mockRepo)

	validID := "60b8d295f1d2c12a34567890"
	req, _ := http.NewRequest("GET", "/analyses/"+validID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var analysis models.Analysis
	err := json.Unmarshal(w.Body.Bytes(), &analysis)
	assert.NoError(t, err)

	// Compare the hex string representation of ObjectID
	assert.Equal(t, validID, analysis.ID.Hex())
}

func TestGetAnalysis_NotFound(t *testing.T) {
	mockRepo := &MockAnalysisRepository{
		GetAnalysisByIDFunc: func(id string) (models.Analysis, error) {
			return models.Analysis{}, errors.New("analysis not found")
		},
	}

	router := setupRouterWithMockRepo(mockRepo)

	req, _ := http.NewRequest("GET", "/analyses/unknown", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateAnalysis_Success(t *testing.T) {
	mockRepo := &MockAnalysisRepository{
		UpdateAnalysisFunc: func(id string, updateReq models.UpdateAnalysisRequest) error {
			return nil
		},
	}

	router := setupRouterWithMockRepo(mockRepo)

	updateReq := models.UpdateAnalysisRequest{
		Title: "New Title",
		Tags:  []string{"tag1", "tag2"},
	}
	body, _ := json.Marshal(updateReq)

	req, _ := http.NewRequest("PUT", "/analyses/123", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateAnalysis_BadRequest(t *testing.T) {
	mockRepo := &MockAnalysisRepository{}

	router := setupRouterWithMockRepo(mockRepo)

	req, _ := http.NewRequest("PUT", "/analyses/123", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteAnalysis_Success(t *testing.T) {
	mockRepo := &MockAnalysisRepository{
		DeleteAnalysisFunc: func(id string) error {
			return nil
		},
	}

	router := setupRouterWithMockRepo(mockRepo)

	req, _ := http.NewRequest("DELETE", "/analyses/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteAnalysis_NotFound(t *testing.T) {
	mockRepo := &MockAnalysisRepository{
		DeleteAnalysisFunc: func(id string) error {
			return errors.New("analysis not found")
		},
	}

	router := setupRouterWithMockRepo(mockRepo)

	req, _ := http.NewRequest("DELETE", "/analyses/unknown", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
