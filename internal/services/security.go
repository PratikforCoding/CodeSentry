package services

import (
	"github.com/PratikforCoding/CodeSentry/internal/models"
	"github.com/PratikforCoding/CodeSentry/internal/utils"
	"strings"
)

type SecurityAnalyzer struct{}

func NewSecurityAnalyzer() *SecurityAnalyzer {
	return &SecurityAnalyzer{}
}

func (sa *SecurityAnalyzer) AnalyzeSecurity(code string) []models.SecurityIssue {
	var issues []models.SecurityIssue

	parser := utils.NewParser(code)
	_ = parser.Tokenize()

	riskyTokens := parser.GetSecurityRiskyTokens()

	for category, tokenList := range riskyTokens {
		for _, token := range tokenList {
			issue := models.SecurityIssue{
				Type:        category,
				Line:        token.Line,
				Column:      token.Col,
				Description: sa.getTokenSecurityDescription(category, token.Value),
				Severity:    sa.getTokenSecuritySeverity(category),
				Pattern:     token.Value,
			}
			issues = append(issues, issue)
		}
	}

	issues = append(issues, sa.analyzeComplexPatterns(code)...)

	return issues
}

func (sa *SecurityAnalyzer) getTokenSecurityDescription(category, tokenValue string) string {
	descriptions := map[string]string{
		"sql_operations":     "SQL operation detected - ensure input is sanitized",
		"eval_functions":     "Dynamic code execution detected - potential code injection risk",
		"file_operations":    "File operation detected - validate file paths and permissions",
		"network_operations": "Network operation detected - ensure secure communication",
	}

	if desc, exists := descriptions[category]; exists {
		return desc + " (Token: " + tokenValue + ")"
	}
	return "Security-sensitive operation detected: " + tokenValue
}

func (sa *SecurityAnalyzer) getTokenSecuritySeverity(category string) string {
	severities := map[string]string{
		"sql_operations":     "HIGH",
		"eval_functions":     "CRITICAL",
		"file_operations":    "MEDIUM",
		"network_operations": "MEDIUM",
	}

	if severity, exists := severities[category]; exists {
		return severity
	}
	return "MEDIUM"
}

func (sa *SecurityAnalyzer) analyzeComplexPatterns(code string) []models.SecurityIssue {
	var issues []models.SecurityIssue
	lines := strings.Split(code, "\n")

	// Keep existing regex patterns for complex security issues
	for lineNum, line := range lines {
		for issueType, pattern := range utils.SecurityPatterns {
			if pattern.MatchString(line) {
				issue := models.SecurityIssue{
					Type:        issueType,
					Line:        lineNum + 1,
					Column:      0,
					Description: sa.getSecurityDescription(issueType),
					Severity:    sa.getSecuritySeverity(issueType),
					Pattern:     line,
				}
				issues = append(issues, issue)
			}
		}
	}

	return issues
}

func (sa *SecurityAnalyzer) getSecurityDescription(issueType string) string {
	descriptions := map[string]string{
		"sql_injection":          "Potential SQL injection vulnerability detected",
		"xss_vulnerability":      "Potential XSS vulnerability in DOM manipulation",
		"path_traversal":         "Path traversal vulnerability detected",
		"hardcoded_password":     "Hardcoded password found in source code",
		"hardcoded_api_key":      "Hardcoded API key or secret found",
		"eval_usage":             "Use of eval() function detected - potential code injection",
		"unsafe_deserialization": "Unsafe deserialization of user input",
	}

	if desc, exists := descriptions[issueType]; exists {
		return desc
	}
	return "Security issue detected"
}

func (sa *SecurityAnalyzer) getSecuritySeverity(issueType string) string {
	severities := map[string]string{
		"sql_injection":          "HIGH",
		"xss_vulnerability":      "HIGH",
		"path_traversal":         "HIGH",
		"hardcoded_password":     "CRITICAL",
		"hardcoded_api_key":      "HIGH",
		"eval_usage":             "MEDIUM",
		"unsafe_deserialization": "HIGH",
	}

	if severity, exists := severities[issueType]; exists {
		return severity
	}
	return "MEDIUM"
}
