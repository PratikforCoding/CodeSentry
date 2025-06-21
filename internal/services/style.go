package services

import (
	"github.com/PratikforCoding/CodeSentry/internal/models"
	"github.com/PratikforCoding/CodeSentry/internal/utils"
	"strings"
)

type StyleAnalyzer struct{}

func NewStyleAnalyzer() *StyleAnalyzer {
	return &StyleAnalyzer{}
}

func (sa *StyleAnalyzer) AnalyzeStyle(code string) []models.StyleSuggestion {
	var suggestions []models.StyleSuggestion
	lines := strings.Split(code, "\n")

	for lineNum, line := range lines {
		if len(line) > 120 {
			suggestions = append(suggestions, models.StyleSuggestion{
				Line:     lineNum + 1,
				Column:   120,
				Message:  "Line exceeds maximum length of 120 characters",
				Rule:     "max-line-length",
				Severity: "WARNING",
			})
		}

		if utils.StylePatterns["trailing_whitespace"].MatchString(line) {
			suggestions = append(suggestions, models.StyleSuggestion{
				Line:     lineNum + 1,
				Column:   len(lines),
				Message:  "Trailing whitespace detected",
				Rule:     "no-trailing-whitespace",
				Severity: "INFO",
			})
		}

		if utils.StylePatterns["mixed_indentation"].MatchString(line) {
			suggestions = append(suggestions, models.StyleSuggestion{
				Line:     lineNum + 1,
				Column:   0,
				Message:  "Mixed spaces and tabs for indentation",
				Rule:     "consistent-indentation",
				Severity: "WARNING",
			})
		}

		if utils.StylePatterns["snake_case_violation"].MatchString(line) ||
			utils.StylePatterns["camel_case_violation"].MatchString(line) {
			suggestions = append(suggestions, models.StyleSuggestion{
				Line:     lineNum + 1,
				Column:   0,
				Message:  "Inconsistent naming convention detected",
				Rule:     "naming-convention",
				Severity: "INFO",
			})
		}

	}
	return suggestions
}
