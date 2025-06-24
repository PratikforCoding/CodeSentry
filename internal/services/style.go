package services

import (
	"regexp"
	"strings"

	"github.com/PratikforCoding/CodeSentry/internal/models"
	"github.com/PratikforCoding/CodeSentry/internal/utils"
)

type StyleAnalyzer struct{}

func NewStyleAnalyzer() *StyleAnalyzer {
	return &StyleAnalyzer{}
}

var identifierRegex = regexp.MustCompile(`\b[a-zA-Z_][a-zA-Z0-9_]*\b`)

func (sa *StyleAnalyzer) AnalyzeStyle(code string) []models.StyleSuggestion {
	var suggestions []models.StyleSuggestion
	lines := strings.Split(code, "\n")

	// Step 1: Scan entire code for identifiers and detect naming styles
	hasSnakeCase := false
	hasCamelCase := false

	for _, line := range lines {
		ids := identifierRegex.FindAllString(line, -1)
		for _, id := range ids {
			if isSnakeCase(id) {
				hasSnakeCase = true
			} else if isCamelCase(id) {
				hasCamelCase = true
			}
		}
	}

	// Step 2: Analyze each line, add style suggestions
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

		// Step 3: If both naming styles are present, flag mixed naming on all lines containing identifiers
		if hasSnakeCase && hasCamelCase {
			ids := identifierRegex.FindAllString(line, -1)
			for _, id := range ids {
				if isSnakeCase(id) || isCamelCase(id) {
					suggestions = append(suggestions, models.StyleSuggestion{
						Line:     lineNum + 1,
						Column:   0,
						Message:  "Mixed naming convention detected",
						Rule:     "naming-convention",
						Severity: "INFO",
					})
					break // avoid duplicate suggestions for the same line
				}
			}
		}
	}

	return suggestions
}

func isSnakeCase(s string) bool {
	return strings.Contains(s, "_") && strings.ToLower(s) == s
}

func isCamelCase(s string) bool {
	if strings.Contains(s, "_") {
		return false
	}
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			return true
		}
	}
	return false
}
