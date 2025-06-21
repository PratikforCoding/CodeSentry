package services

import (
	"github.com/PratikforCoding/CodeSentry/internal/models"
	"github.com/PratikforCoding/CodeSentry/internal/repository"
	"github.com/PratikforCoding/CodeSentry/internal/utils"
	"strings"
)

type AnalyzerService struct {
	languageDetector   *LanguageDetector
	complexityAnalyzer *ComplexityAnalyzer
	securityAnalyzer   *SecurityAnalyzer
	styleAnalyzer      *StyleAnalyzer
	repo               *repository.AnalysisRepository
}

func NewAnalyzerService() *AnalyzerService {
	return &AnalyzerService{
		languageDetector:   NewLanguageDetector(),
		complexityAnalyzer: NewComplexityAnalyzer(),
		securityAnalyzer:   NewSecurityAnalyzer(),
		styleAnalyzer:      NewStyleAnalyzer(),
		repo:               repository.NewAnalysisRepository(),
	}
}

func (as *AnalyzerService) AnalyzeCode(req models.AnalyzeRequest) models.AnalysisResponse {

	// Detect language
	language := req.Language
	if language == "" {
		detectLang := as.languageDetector.DetectLanguage(req.Code)
		language = string(detectLang)
	}

	response := models.AnalysisResponse{
		Language: language,
	}

	// Analyze complexity
	if req.Options.CheckComplexity {
		response.ComplexityScore = as.complexityAnalyzer.AnalyzeComplexity(req.Code)
	}

	// Analyze security
	if req.Options.CheckSecurity {
		response.SecurityIssues = as.securityAnalyzer.AnalyzeSecurity(req.Code)
	}

	// Analyze style
	if req.Options.CheckStyle {
		response.StyleSuggestions = as.styleAnalyzer.AnalyzeStyle(req.Code)
	}

	// Analyze metrics {
	if req.Options.CheckMetrics {
		response.Metrics = as.calculateMetrics(req.Code)
	}

	// Calculate overall acore
	response.OverallScore = as.calculateOverallScore(response)

	_ = as.repo.SaveAnalysis(req, response)

	return response
}

func (as *AnalyzerService) calculateMetrics(code string) models.CodeMetrics {

	parser := utils.NewParser(code)
	_ = parser.Tokenize()

	language := as.languageDetector.DetectLanguage(code)

	totalLines, blankLines, commentLines := utils.CountLines(code)
	functions := as.complexityAnalyzer.CountFunctions(code, language)
	nestingDepth := as.complexityAnalyzer.CalculateNestingDepth(code)
	complexity := as.complexityAnalyzer.AnalyzeComplexity(code)

	linesOfCode := totalLines - blankLines - commentLines
	commentRatio := float64(commentLines) / float64(totalLines)

	lines := strings.Split(code, "\n")
	totalLength := 0
	for _, line := range lines {
		totalLength += len(line)
	}
	avgLineLength := float64(totalLength) / float64(len(lines))

	identifierCount := parser.CountTokenType(models.IDENTIFIER)
	keywordCount := parser.CountTokenType(models.KEYWORD)
	operatorCount := parser.CountTokenType(models.OPERATOR)

	return models.CodeMetrics{
		LinesOfCode:          linesOfCode,
		LinesOfComments:      commentLines,
		BlankLines:           blankLines,
		IdentifierCount:      identifierCount,
		KeywordCount:         keywordCount,
		OperatorCount:        operatorCount,
		Functions:            functions,
		Classes:              0, // TODO: Implement class counting
		CommentRatio:         commentRatio,
		AverageLineLength:    avgLineLength,
		MaxNestingDepth:      nestingDepth,
		CyclomaticComplexity: complexity,
	}
}

func (as *AnalyzerService) calculateOverallScore(response models.AnalysisResponse) float64 {
	score := 100.0
	if response.ComplexityScore > 10 {
		score -= float64(response.ComplexityScore-10) * 2
	}

	for _, issue := range response.SecurityIssues {
		switch issue.Severity {
		case "CRITICAL":
			score -= 20
		case "HIGH":
			score -= 10
		case "MEDIUM":
			score -= 5
		case "LOW":
			score -= 2
		}
	}

	for _, suggestion := range response.StyleSuggestions {
		switch suggestion.Severity {
		case "WARNING":
			score -= 2
		case "INFO":
			score -= 1
		}
	}

	if score < 0 {
		score = 0
	}
	return score
}
