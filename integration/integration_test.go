package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/PratikforCoding/CodeSentry/internal/handlers"
	"github.com/PratikforCoding/CodeSentry/internal/models"
	"github.com/PratikforCoding/CodeSentry/internal/repository"
	"github.com/PratikforCoding/CodeSentry/internal/services"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	testCtx             context.Context
	testDBClient        *mongo.Client
	testDB              *mongo.Database
	testRepo            repository.AnalysisRepositoryInterface
	testAnalyzer        *services.SecurityAnalyzer
	testAnalysesHandler *handlers.AnalysesHandler
	testAnalyzerHandler *handlers.AnalyzerHandler
)

func setupMongoClient() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:example@localhost:27017/codesentry_test?authSource=admin"))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}

func TestMain(m *testing.M) {
	var err error
	testCtx = context.Background()

	testDBClient, err = setupMongoClient()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	testDB = testDBClient.Database("codesentry_test")

	// Clean test database before running tests
	if err := testDB.Drop(testCtx); err != nil {
		log.Fatalf("Failed to drop test database: %v", err)
	}

	// Initialize repository and handlers
	testRepo = repository.NewAnalysisRepository(testDB)
	testAnalyzer = services.NewSecurityAnalyzer()
	testAnalysesHandler = handlers.NewAnalysesHandlerWithRepo(testRepo)
	testAnalyzerHandler = handlers.NewAnalyzerHandler(testDB)

	code := m.Run()

	// Disconnect MongoDB client after tests
	if err := testDBClient.Disconnect(testCtx); err != nil {
		log.Printf("Error disconnecting MongoDB client: %v", err)
	}

	os.Exit(code)
}

func TestAnalyzeCodeHandler(t *testing.T) {
	// Prepare the AnalyzeRequest payload
	reqBody := models.AnalyzeRequest{
		Code: `
			eval("dangerous")
			SELECT * FROM users WHERE username = 'admin'
			password = "1234"
		`,
		Language: "go",
		Options: models.AnalyzeOptions{
			CheckSecurity: true,
		},
	}

	// Marshal payload to JSON
	bodyBytes, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	// Create a new HTTP POST request
	req := httptest.NewRequest(http.MethodPost, "/analyze", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create a Gin context with the ResponseRecorder
	gin.SetMode(gin.TestMode) // Set Gin to test mode
	c, _ := gin.CreateTestContext(rr)
	c.Request = req

	// Call your AnalyzeCode handler with the Gin context
	testAnalyzerHandler.AnalyzeCode(c)

	// Assert the HTTP status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse the response body
	var resp models.AnalysisResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)

	// Assert that security issues are detected
	assert.NotEmpty(t, resp.SecurityIssues, "Expected security issues in response")

	// Optionally check for specific issue types
	var foundEval, foundSQL, foundPassword bool
	for _, issue := range resp.SecurityIssues {
		switch issue.Type {
		case "eval_functions":
			foundEval = true
		case "sql_injection":
			foundSQL = true
		case "hardcoded_password":
			foundPassword = true
		}
	}

	assert.True(t, foundEval, "Should detect eval usage")
	assert.True(t, foundSQL, "Should detect SQL injection")
	assert.True(t, foundPassword, "Should detect hardcoded password")
}

func TestGetAnalysesHandler(t *testing.T) {
	// Insert a dummy analysis to test retrieval
	err := testRepo.SaveAnalysis(models.AnalyzeRequest{
		Code:     "print('hello')",
		Language: "python",
	}, models.AnalysisResponse{
		Language: "python",
	})
	assert.NoError(t, err, "Failed to save analysis")

	req := httptest.NewRequest(http.MethodGet, "/analyses?language=python", nil)
	rr := httptest.NewRecorder()

	// Gin context setup
	c, _ := gin.CreateTestContext(rr)
	c.Request = req

	testAnalysesHandler.GetAnalyses(c)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)

	analyses, ok := resp["analyses"].([]interface{})
	assert.True(t, ok)
	assert.NotEmpty(t, analyses)
}

func seedTestData(t *testing.T) string {
	err := testRepo.SaveAnalysis(models.AnalyzeRequest{
		Code:     "print('test data')",
		Language: "python",
	}, models.AnalysisResponse{})
	assert.NoError(t, err)

	analyses, err := testRepo.GetAllAnalyses("python")
	assert.NoError(t, err)
	assert.NotEmpty(t, analyses)

	return analyses[0].ID.Hex()
}

func TestUpdateDeleteAnalysisHandler(t *testing.T) {
	// Save an analysis to update and delete
	err := testRepo.SaveAnalysis(models.AnalyzeRequest{
		Code:     "print('update')",
		Language: "python",
	}, models.AnalysisResponse{})
	assert.NoError(t, err)

	analyses, err := testRepo.GetAllAnalyses("python")
	assert.NoError(t, err)
	assert.NotEmpty(t, analyses)
	if len(analyses) == 0 {
		t.Fatal("No analyses found, cannot proceed")
	}

	analysisID := analyses[0].ID.Hex()

	// Update request
	updateReq := models.UpdateAnalysisRequest{
		Title: "Updated Title",
	}
	bodyBytes, err := json.Marshal(updateReq)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, "/analyses/"+analysisID, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: analysisID}}

	testAnalysesHandler.UpdateAnalysis(c)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Delete request
	reqDel := httptest.NewRequest(http.MethodDelete, "/analyses/"+analysisID, nil)
	rrDel := httptest.NewRecorder()
	cDel, _ := gin.CreateTestContext(rrDel)
	cDel.Request = reqDel
	cDel.Params = gin.Params{gin.Param{Key: "id", Value: analysisID}}

	testAnalysesHandler.DeleteAnalysis(cDel)
	assert.Equal(t, http.StatusOK, rrDel.Code)
}
