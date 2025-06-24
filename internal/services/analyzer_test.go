package services_test

import (
	"testing"

	"github.com/PratikforCoding/CodeSentry/internal/models"
	"github.com/PratikforCoding/CodeSentry/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mock definitions ---

// Mock for LanguageDetector
type MockLanguageDetector struct {
	mock.Mock
}

func (m *MockLanguageDetector) DetectLanguage(code string) models.Language {
	args := m.Called(code)
	return args.Get(0).(models.Language) // Extract first return value and assert type
}

// Mock for ComplexityAnalyzer
type MockComplexityAnalyzer struct {
	mock.Mock
}

func (m *MockComplexityAnalyzer) AnalyzeComplexity(code string) int {
	args := m.Called(code)
	return args.Int(0)
}

func (m *MockComplexityAnalyzer) CountFunctions(code string, lang models.Language) int {
	args := m.Called(code, lang)
	return args.Int(0)
}

func (m *MockComplexityAnalyzer) CalculateNestingDepth(code string, lang models.Language) int {
	args := m.Called(code, lang)
	return args.Int(0)
}

// Mock for SecurityAnalyzer
type MockSecurityAnalyzer struct {
	mock.Mock
}

func (m *MockSecurityAnalyzer) AnalyzeSecurity(code string) []models.SecurityIssue {
	args := m.Called(code)
	return args.Get(0).([]models.SecurityIssue)
}

// Mock for StyleAnalyzer
type MockStyleAnalyzer struct {
	mock.Mock
}

func (m *MockStyleAnalyzer) AnalyzeStyle(code string) []models.StyleSuggestion {
	args := m.Called(code)
	return args.Get(0).([]models.StyleSuggestion)
}

// Mock for repository.AnalysisRepository
type MockAnalysisRepository struct {
	mock.Mock
}

func (m *MockAnalysisRepository) SaveAnalysis(req models.AnalyzeRequest, res models.AnalysisResponse) error {
	args := m.Called(req, res)
	return args.Error(0)
}

func (m *MockAnalysisRepository) GetAllAnalyses(language string) ([]models.Analysis, error) {
	args := m.Called(language)
	return args.Get(0).([]models.Analysis), args.Error(1)
}

func (m *MockAnalysisRepository) GetAnalysisByID(id string) (models.Analysis, error) {
	args := m.Called(id)
	return args.Get(0).(models.Analysis), args.Error(1)
}

func (m *MockAnalysisRepository) UpdateAnalysis(id string, updateReq models.UpdateAnalysisRequest) error {
	args := m.Called(id, updateReq)
	return args.Error(0)
}

func (m *MockAnalysisRepository) DeleteAnalysis(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// --- Test function ---

func TestAnalyzeCode(t *testing.T) {
	// Create mocks
	langDetector := new(MockLanguageDetector)
	complexityAnalyzer := new(MockComplexityAnalyzer)
	securityAnalyzer := new(MockSecurityAnalyzer)
	styleAnalyzer := new(MockStyleAnalyzer)
	mockRepo := new(MockAnalysisRepository)

	// Inject mocks via the new constructor
	service := services.NewAnalyzerServiceWithDeps(
		langDetector,
		complexityAnalyzer,
		securityAnalyzer,
		styleAnalyzer,
		mockRepo,
	)

	// Define input request with options enabled
	req := models.AnalyzeRequest{
		Code:     "package main\nfunc main() {}",
		Language: "",
		Options: models.AnalyzeOptions{
			CheckComplexity: true,
			CheckSecurity:   true,
			CheckStyle:      true,
			CheckMetrics:    false,
		},
	}

	// Setup mock expectations
	langDetector.On("DetectLanguage", req.Code).Return(models.Go)
	complexityAnalyzer.On("AnalyzeComplexity", req.Code).Return(5)
	securityAnalyzer.On("AnalyzeSecurity", req.Code).Return([]models.SecurityIssue{
		{Severity: "HIGH", Description: "Test security issue"},
	})
	styleAnalyzer.On("AnalyzeStyle", req.Code).Return([]models.StyleSuggestion{
		{
			Severity: "WARNING",
			Message:  "Use camelCase",
			Line:     1,                   // optionally provide line number
			Column:   0,                   // optionally provide column number
			Rule:     "naming-convention", // optionally provide rule name
		},
	})
	mockRepo.On("SaveAnalysis", req, mock.AnythingOfType("models.AnalysisResponse")).Return(nil)

	// Call the method under test
	resp := service.AnalyzeCode(req)

	// Assertions
	assert.Equal(t, "go", resp.Language)
	assert.Equal(t, 5, resp.ComplexityScore)
	assert.Len(t, resp.SecurityIssues, 1)
	assert.Len(t, resp.StyleSuggestions, 1)

	expectedScore := 100.0 - 10 - 2
	assert.Equal(t, expectedScore, resp.OverallScore)

	// Verify SaveAnalysis called once
	mockRepo.AssertCalled(t, "SaveAnalysis", req, resp)

	// Assert all expectations met
	langDetector.AssertExpectations(t)
	complexityAnalyzer.AssertExpectations(t)
	securityAnalyzer.AssertExpectations(t)
	styleAnalyzer.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}
