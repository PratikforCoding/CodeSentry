package services_test

import (
	"strings"
	"testing"

	"github.com/PratikforCoding/CodeSentry/internal/models"
	"github.com/PratikforCoding/CodeSentry/internal/services"
	"github.com/PratikforCoding/CodeSentry/internal/utils"
	"github.com/stretchr/testify/assert"
)

// Helper: Create a Parser with given code and tokenize
func newParserWithCode(t *testing.T, code string) *utils.Parser {
	parser := utils.NewParser(code)
	parser.Tokenize()
	return parser
}

func TestIsRiskyContext(t *testing.T) {
	code := `userInput = input("Enter data") + "unsafe"`
	parser := newParserWithCode(t, code)

	// Simulate token at line 1, col 1 with keyword "input"
	token := models.Token{
		Value: "input",
		Line:  1,
		Col:   strings.Index(code, "input") + 1,
	}

	assert.True(t, parser.IsRiskyContext(token, "input"), "Should detect risky context with input and +")
	assert.False(t, parser.IsRiskyContext(token, "safeKeyword"), "Should not detect risky context for unrelated keyword")
}

func TestGetVulnerabilityScore(t *testing.T) {
	code := `
		eval("dangerous")
		SELECT * FROM users WHERE username = 'admin'
		password = "1234"
	`
	parser := newParserWithCode(t, code)

	score := parser.GetVulnerabilityScore()
	assert.GreaterOrEqual(t, score, 10, "Score should reflect detected vulnerabilities")
}

func TestGetVulnerabilityByCategory(t *testing.T) {
	code := `
		eval("dangerous")
		SELECT * FROM users WHERE username = 'admin'
		password = "1234"
	`
	parser := newParserWithCode(t, code)

	vulns := parser.GetVulnerabilityByCategory()
	assert.Contains(t, vulns["critical"], "eval_usage", "Should detect eval_usage as critical")
	assert.Contains(t, vulns["critical"], "sql_injection", "Should detect sql_injection as critical")
	assert.Contains(t, vulns["medium"], "hardcoded_password", "Should detect hardcoded_password as medium")
}

func TestGetJavaMethodSignatures(t *testing.T) {
	code := `
		public class Test {
			public void foo() {
			}
			private static int bar(int x) throws Exception {
			}
			@Deprecated
			public String baz() {
			}
		}
	`
	parser := newParserWithCode(t, code)

	signatures := parser.GetJavaMethodSignatures()
	assert.NotEmpty(t, signatures, "Should detect Java method signatures")
	assert.Contains(t, strings.Join(signatures, " "), "public void foo()", "Should include foo method signature")
	assert.Contains(t, strings.Join(signatures, " "), "private static int bar(int x)", "Should include bar method signature")
	assert.Contains(t, strings.Join(signatures, " "), "public String baz()", "Should include baz method signature")
}

func TestDetectSecurityHotspots(t *testing.T) {
	code := `
		eval("dangerous")
		SELECT * FROM users WHERE username = 'admin'
		password = "1234"
	`
	parser := newParserWithCode(t, code)

	hotspots := parser.DetectSecurityHotspots()
	assert.NotEmpty(t, hotspots, "Should detect security hotspots")

	var foundEval, foundSQL, foundPassword bool
	for _, hs := range hotspots {
		switch hs.Category {
		case "eval_usage":
			foundEval = true
			assert.Equal(t, "CRITICAL", hs.Severity)
		case "sql_injection":
			foundSQL = true
			assert.Equal(t, "CRITICAL", hs.Severity)
		case "hardcoded_password":
			foundPassword = true
			assert.Equal(t, "MEDIUM", hs.Severity)
		}
	}
	assert.True(t, foundEval, "Should find eval_usage hotspot")
	assert.True(t, foundSQL, "Should find sql_injection hotspot")
	assert.True(t, foundPassword, "Should find hardcoded_password hotspot")
}

func TestSecurityAnalyzer_AnalyzeSecurity(t *testing.T) {
	analyzer := services.NewSecurityAnalyzer()

	code := `
		eval("dangerous")
		SELECT * FROM users WHERE username = 'admin'
		password = "1234"
		open("/etc/passwd")
		http.get("http://example.com")
	`

	issues := analyzer.AnalyzeSecurity(code)
	assert.NotEmpty(t, issues, "Should detect security issues")

	var foundSQLOp, foundEval, foundFile, foundNetwork bool
	for _, issue := range issues {
		switch issue.Type {
		case "sql_operations":
			foundSQLOp = true
			assert.Contains(t, issue.Description, "SQL operation detected")
		case "eval_functions":
			foundEval = true
			assert.Contains(t, issue.Description, "Dynamic code execution")
		case "file_operations":
			foundFile = true
			assert.Contains(t, issue.Description, "File operation detected")
		case "network_operations":
			foundNetwork = true
			assert.Contains(t, issue.Description, "Network operation detected")
		}
	}
	assert.True(t, foundSQLOp, "Should detect SQL operation")
	assert.True(t, foundEval, "Should detect eval function")
	assert.True(t, foundFile, "Should detect file operation")
	assert.True(t, foundNetwork, "Should detect network operation")
}

func TestSecurityAnalyzer_DescriptionsAndSeverities(t *testing.T) {
	analyzer := services.NewSecurityAnalyzer()

	// Test token security description and severity
	desc := analyzer.GetTokenSecurityDescription("sql_operations", "SELECT")
	assert.Contains(t, desc, "SQL operation detected")
	sev := analyzer.GetTokenSecuritySeverity("sql_operations")
	assert.Equal(t, "HIGH", sev)

	desc = analyzer.GetTokenSecurityDescription("eval_functions", "eval")
	assert.Contains(t, desc, "Dynamic code execution")
	sev = analyzer.GetTokenSecuritySeverity("eval_functions")
	assert.Equal(t, "CRITICAL", sev)

	desc = analyzer.GetTokenSecurityDescription("unknown", "token")
	assert.Contains(t, desc, "Security-sensitive operation detected")
	sev = analyzer.GetTokenSecuritySeverity("unknown")
	assert.Equal(t, "MEDIUM", sev)

	// Test general security description and severity
	desc = analyzer.GetSecurityDescription("sql_injection")
	assert.Contains(t, desc, "SQL injection vulnerability")
	sev = analyzer.GetSecuritySeverity("sql_injection")
	assert.Equal(t, "HIGH", sev)

	desc = analyzer.GetSecurityDescription("hardcoded_password")
	assert.Contains(t, desc, "Hardcoded password")
	sev = analyzer.GetSecuritySeverity("hardcoded_password")
	assert.Equal(t, "CRITICAL", sev)

	desc = analyzer.GetSecurityDescription("unknown_issue")
	assert.Contains(t, desc, "Security issue detected")
	sev = analyzer.GetSecuritySeverity("unknown_issue")
	assert.Equal(t, "MEDIUM", sev)
}
