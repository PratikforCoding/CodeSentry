package utils

import (
	"github.com/PratikforCoding/CodeSentry/internal/models"
	"log"
	"regexp"
	"strings"
	"unicode"
)

type Parser struct {
	code     string
	tokens   []models.Token
	position int
	line     int
	column   int
}

func NewParser(code string) *Parser {
	log.Println("Debug: creating tokenizer")
	return &Parser{
		code:   code,
		tokens: make([]models.Token, 0),
		line:   1,
		column: 1,
	}
}

func (p *Parser) Tokenize() []models.Token {
	p.tokens = []models.Token{}
	runes := []rune(p.code)

	for i := 0; i < len(runes); i++ {
		char := runes[i]

		// Handle newlines
		if char == '\n' {
			p.line++
			p.column = 1
			continue
		}

		// Handle whitespace
		if unicode.IsSpace(char) {
			p.column++
			continue
		}

		// Handle comments
		if char == '/' && i+1 < len(runes) {
			// Single-line comment
			if runes[i+1] == '/' {
				start := i
				startCol := p.column
				for i < len(runes) && runes[i] != '\n' {
					i++
				}
				p.addTokenWithPosition(models.COMMENT, string(runes[start:i]), p.line, startCol)
				i-- // Adjust because outer loop will increment
				p.column = startCol + (i - start)
				continue
			}
			// Multi-line comment
			if runes[i+1] == '*' {
				start := i
				startLine := p.line
				startCol := p.column
				i += 2 // Skip /*

				for i+1 < len(runes) {
					if runes[i] == '*' && runes[i+1] == '/' {
						i += 2 // Skip */
						break
					}
					if runes[i] == '\n' {
						p.line++
						p.column = 1
					} else {
						p.column++
					}
					i++
				}
				p.addTokenWithPosition(models.COMMENT, string(runes[start:i]), startLine, startCol)
				i-- // Adjust for outer loop
				continue
			}
		}

		// Handle Python # comments
		if char == '#' {
			start := i
			startCol := p.column
			for i < len(runes) && runes[i] != '\n' {
				i++
			}
			p.addTokenWithPosition(models.COMMENT, string(runes[start:i]), p.line, startCol)
			i-- // Adjust because outer loop will increment
			p.column = startCol + (i - start)
			continue
		}

		// Handle string literals
		if char == '"' || char == '\'' {
			quote := char
			start := i
			startCol := p.column
			i++ // Skip opening quote

			for i < len(runes) && runes[i] != quote {
				if runes[i] == '\\' && i+1 < len(runes) {
					i += 2 // Skip escaped character
				} else {
					if runes[i] == '\n' {
						p.line++
						p.column = 1
					} else {
						p.column++
					}
					i++
				}
			}

			if i < len(runes) {
				i++ // Skip closing quote
			}

			p.addTokenWithPosition(models.LITERAL, string(runes[start:i]), p.line, startCol)
			i-- // Adjust for outer loop
			continue
		}

		// Handle numeric literals
		if unicode.IsDigit(char) || (char == '.' && i+1 < len(runes) && unicode.IsDigit(runes[i+1])) {
			start := i
			startCol := p.column

			// Handle different number formats
			if char == '0' && i+1 < len(runes) {
				next := runes[i+1]
				if next == 'x' || next == 'X' {
					// Hexadecimal
					i += 2
					for i < len(runes) && (unicode.IsDigit(runes[i]) ||
						(runes[i] >= 'a' && runes[i] <= 'f') ||
						(runes[i] >= 'A' && runes[i] <= 'F')) {
						i++
					}
				} else if next == 'b' || next == 'B' {
					// Binary
					i += 2
					for i < len(runes) && (runes[i] == '0' || runes[i] == '1') {
						i++
					}
				} else if next == 'o' || next == 'O' {
					// Octal
					i += 2
					for i < len(runes) && runes[i] >= '0' && runes[i] <= '7' {
						i++
					}
				}
			}

			if i == start || (i == start+1 && char == '0') {
				// Regular decimal number
				hasDot := false
				for i < len(runes) && (unicode.IsDigit(runes[i]) || runes[i] == '.') {
					if runes[i] == '.' {
						if hasDot {
							break // Second dot, stop parsing
						}
						hasDot = true
					}
					i++
				}

				// Handle scientific notation
				if i < len(runes) && (runes[i] == 'e' || runes[i] == 'E') {
					i++
					if i < len(runes) && (runes[i] == '+' || runes[i] == '-') {
						i++
					}
					for i < len(runes) && unicode.IsDigit(runes[i]) {
						i++
					}
				}
			}

			p.addTokenWithPosition(models.LITERAL, string(runes[start:i]), p.line, startCol)
			i-- // Adjust for outer loop
			p.column = startCol + (i - start)
			continue
		}

		// Handle identifiers and keywords
		if unicode.IsLetter(char) || char == '_' {
			start := i
			startCol := p.column

			for i < len(runes) && (unicode.IsLetter(runes[i]) ||
				unicode.IsDigit(runes[i]) || runes[i] == '_') {
				i++
			}

			word := string(runes[start:i])

			if p.isKeyword(word) {
				p.addTokenWithPosition(models.KEYWORD, word, p.line, startCol)
			} else {
				p.addTokenWithPosition(models.IDENTIFIER, word, p.line, startCol)
			}

			i-- // Adjust for outer loop
			p.column = startCol + len(word)
			continue
		}

		// Handle operators and delimiters
		startCol := p.column
		if p.isOperator(char) {
			// Check for multi-character operators
			if i+1 < len(runes) {
				twoChar := string(runes[i : i+2])
				if p.isMultiCharOperator(twoChar) {
					p.addTokenWithPosition(models.OPERATOR, twoChar, p.line, startCol)
					i++ // Extra increment for two-char operator
					p.column += 2
					continue
				}
			}
			p.addTokenWithPosition(models.OPERATOR, string(char), p.line, startCol)
		} else if p.isDelimiter(char) {
			p.addTokenWithPosition(models.DELIMITER, string(char), p.line, startCol)
		} else {
			p.addTokenWithPosition(models.UNKNOWN, string(char), p.line, startCol)
		}

		p.column++
	}

	return p.tokens
}

func (p *Parser) addTokenWithPosition(tokenType models.TokenType, value string, line, column int) {
	token := models.Token{
		Type:  tokenType,
		Value: value,
		Line:  line,
		Col:   column,
	}
	p.tokens = append(p.tokens, token)
}

func (p *Parser) isMultiCharOperator(op string) bool {
	multiCharOps := []string{
		"==", "!=", "<=", ">=", "&&", "||", "++", "--", "+=", "-=",
		"*=", "/=", "%=", "&=", "|=", "^=", "<<", ">>", "**", "//",
		"::", "->", "=>", "??", "?.", "!!", "<-",
	}

	for _, mop := range multiCharOps {
		if op == mop {
			return true
		}
	}
	return false
}

func (p *Parser) addToken(tokenType models.TokenType, value string) {
	p.tokens = append(p.tokens, models.Token{
		Type:  tokenType,
		Value: value,
		Line:  p.line,
		Col:   p.column,
	})
}

func (p *Parser) isOperator(char rune) bool {
	operators := "+-*/%=<>!&|^~"
	return strings.ContainsRune(operators, char)
}

func (p *Parser) isDelimiter(char rune) bool {
	delimiters := "()[]{},.;:"
	return strings.ContainsRune(delimiters, char)
}

func CountLines(code string) (total, blank, comment int) {
	lines := strings.Split(code, "\n")

	// Remove the last empty line if it exists due to trailing \n
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	total = len(lines)
	inBlockComment := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			blank++
			continue
		}

		// Handle block comments
		if strings.Contains(trimmed, "/*") && !inBlockComment {
			inBlockComment = true
		}
		if strings.Contains(trimmed, "*/") && inBlockComment {
			inBlockComment = false
			// Check if there's any code after the closing comment
			afterComment := strings.Split(trimmed, "*/")
			if len(afterComment) > 1 && strings.TrimSpace(afterComment[1]) == "" {
				comment++
				continue
			}
		}
		if inBlockComment {
			comment++
			continue
		}

		// Handle single-line comments
		if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#") {
			comment++
		}
	}
	return total, blank, comment
}

func (p *Parser) GetTokensByType(tokenType models.TokenType) []models.Token {
	var result []models.Token
	for _, token := range p.tokens {
		if token.Type == tokenType {
			result = append(result, token)
		}
	}
	return result
}

func (p *Parser) GetTokensByValue(value string) []models.Token {
	var result []models.Token
	for _, token := range p.tokens {
		if strings.EqualFold(token.Value, value) {
			result = append(result, token)
		}
	}
	return result
}

func (p *Parser) CountTokenType(tokenType models.TokenType) int {
	return len(p.GetTokensByType(tokenType))
}

func (p *Parser) CountKeyword(keyword string) int {
	count := 0
	keywords := p.GetTokensByType(models.KEYWORD)
	for _, token := range keywords {
		if strings.EqualFold(token.Value, keyword) {
			count++
		}
	}
	return count
}

func (p *Parser) GetComplexityTokens() []models.Token {
	complexityKeywords := []string{
		"if", "else", "elif", "for", "while", "switch", "case",
		"try", "catch", "except", "finally", "break", "continue", "return",
	}

	var result []models.Token
	for _, keyword := range complexityKeywords {
		result = append(result, p.GetTokensByValue(keyword)...)
	}

	// Add operators that increase complexity
	operators := p.GetTokensByType(models.OPERATOR)
	for _, op := range operators {
		if op.Value == "&&" || op.Value == "||" || op.Value == "?" {
			result = append(result, op)
		}
	}

	return result
}

func (p *Parser) GetFunctionTokens(language models.Language) []models.Token {
	var result []models.Token
	seenPositions := make(map[int]bool)

	langKey := strings.ToLower(string(language))
	patterns, ok := FunctionPatterns[langKey]
	if !ok {
		// Unsupported language, return empty slice
		return result
	}

	for _, pattern := range patterns {
		matches := pattern.FindAllStringSubmatchIndex(p.code, -1)
		for _, match := range matches {
			start := match[0]
			end := match[1]

			if seenPositions[start] {
				continue
			}
			seenPositions[start] = true

			line, col := p.getLineAndColumn(start)
			tokenValue := p.code[start:end]

			token := models.Token{
				Type:  models.KEYWORD,
				Value: tokenValue,
				Line:  line,
				Col:   col,
			}
			result = append(result, token)
		}
	}

	return result
}

func contains(slice []models.Token, item models.Token) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func (p *Parser) GetSecurityRiskyTokens() map[string][]models.Token {
	result := make(map[string][]models.Token)

	// Check each security pattern against the code
	for category, pattern := range SecurityPatterns {
		matches := pattern.FindAllStringSubmatch(p.code, -1)
		for _, match := range matches {
			if len(match) > 0 {
				// Find the position of this match in the code
				matchStart := strings.Index(p.code, match[0])
				if matchStart != -1 {
					line, col := p.getLineAndColumn(matchStart)
					token := models.Token{
						Type:  models.IDENTIFIER, // or KEYWORD depending on the match
						Value: match[0],
						Line:  line,
						Col:   col,
					}
					result[category] = append(result[category], token)
				}
			}
		}
	}

	// Enhanced token-based detection for additional context
	additionalRiskyPatterns := map[string][]string{
		"sql_operations":     {"select", "insert", "update", "delete", "create", "drop", "alter", "union", "where", "from"},
		"eval_functions":     {"eval", "exec", "system", "shell_exec", "popen", "subprocess", "Runtime.getRuntime"},
		"file_operations":    {"open", "read", "write", "file", "fopen", "fread", "fwrite", "include", "require"},
		"network_operations": {"http", "request", "socket", "connect", "curl", "wget", "fetch", "XMLHttpRequest"},
		"crypto_operations":  {"encrypt", "decrypt", "hash", "md5", "sha1", "aes", "des", "rsa"},
		"auth_operations":    {"login", "authenticate", "authorize", "session", "token", "jwt", "oauth"},
		"deserialization":    {"pickle", "yaml", "json", "serialize", "unserialize", "ObjectInputStream"},
		"command_execution":  {"cmd", "command", "shell", "bash", "powershell", "execute"},
	}

	for category, keywords := range additionalRiskyPatterns {
		for _, keyword := range keywords {
			tokens := p.GetTokensByValue(keyword)
			for _, token := range tokens {
				// Check context to determine if it's actually risky
				if p.isRiskyContext(token, keyword) {
					result[category] = append(result[category], token)
				}
			}
		}
	}

	return result
}

func (p *Parser) AnalyzeNestingDepth(language models.Language) int {
	maxDepth := 0
	currentDepth := 0
	if language == "python" {
		return p.analyzeIndentationDepth()
	}
	for _, token := range p.tokens {
		if token.Type == models.DELIMITER {
			switch token.Value {
			case "{", "(", "[":
				currentDepth++
				if currentDepth > maxDepth {
					maxDepth = currentDepth
				}
			case "}", ")", "]":
				if currentDepth > 0 {
					currentDepth--
				}
			}
		}
	}

	return maxDepth
}

func (p *Parser) analyzeIndentationDepth() int {
	lines := strings.Split(p.code, "\n")
	maxDepth := 0

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue // Skip empty lines
		}

		// Count leading spaces/tabs
		indentLevel := 0
		for _, char := range line {
			if char == ' ' {
				indentLevel++
			} else if char == '\t' {
				indentLevel += 4 // Assume tab = 4 spaces
			} else {
				break
			}
		}

		// Convert to logical nesting depth (assuming 4 spaces per level)
		depth := indentLevel / 4
		if depth > maxDepth {
			maxDepth = depth
		}
	}

	return maxDepth
}

func (p *Parser) GetLanguageScore(lang models.Language) int {
	score := 0
	keywords := p.GetTokensByType(models.KEYWORD)
	langKeywords := GetLanguageKeywords()[lang]

	for _, token := range keywords {
		for _, langKeyword := range langKeywords {
			if strings.EqualFold(token.Value, langKeyword) {
				score++
				break
			}
		}
	}

	return score
}

func (p *Parser) getLineAndColumn(pos int) (int, int) {
	line := 1
	col := 1
	for i := 0; i < pos && i < len(p.code); i++ {
		if p.code[i] == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}
	return line, col
}

/*
	func (p *Parser) isFunctionDeclaration(token models.Token) bool {
		// Find the token in the original code and check surrounding context
		codeLines := strings.Split(p.code, "\n")
		if token.Line <= len(codeLines) {
			line := codeLines[token.Line-1]

			// Check if the line contains function declaration patterns
			for _, pattern := range FunctionPatterns {
				if pattern.MatchString(line) {
					return true
				}
			}

			// Check if followed by identifier and parentheses
			return strings.Contains(line, token.Value) &&
				strings.Contains(line, "(") &&
				strings.Index(line, token.Value) < strings.Index(line, "(")
		}
		return false
	}
*/
func (p *Parser) isRiskyContext(token models.Token, keyword string) bool {
	// Get the line containing the token
	codeLines := strings.Split(p.code, "\n")
	if token.Line <= len(codeLines) {
		line := codeLines[token.Line-1]

		// Check for risky patterns in the context
		riskyContexts := []string{
			"\\+",                                       // String concatenation
			"\\$",                                       // Variable interpolation
			"input", "request", "param", "argv", "args", // User input
			"=", ":=", // Assignment
			"\\(", // Function call
		}

		for _, context := range riskyContexts {
			if matched, _ := regexp.MatchString("(?i)"+keyword+".*"+context, line); matched {
				return true
			}
			if matched, _ := regexp.MatchString("(?i)"+context+".*"+keyword, line); matched {
				return true
			}
		}
	}
	return false
}

func (p *Parser) GetVulnerabilityScore() int {
	score := 0
	riskyTokens := p.GetSecurityRiskyTokens()

	// Weight different categories
	weights := map[string]int{
		"sql_injection":          10,
		"xss_vulnerability":      8,
		"command_injection":      10,
		"eval_usage":             9,
		"path_traversal":         7,
		"hardcoded_password":     6,
		"hardcoded_api_key":      6,
		"unsafe_deserialization": 8,
		"weak_crypto":            5,
		"auth_bypass":            9,
	}

	for category, tokens := range riskyTokens {
		if weight, exists := weights[category]; exists {
			score += len(tokens) * weight
		} else {
			score += len(tokens) * 3 // Default weight
		}
	}

	return score
}

func (p *Parser) GetVulnerabilityByCategory() map[string]map[string][]models.Token {
	riskyTokens := p.GetSecurityRiskyTokens()

	result := map[string]map[string][]models.Token{
		"critical": make(map[string][]models.Token),
		"high":     make(map[string][]models.Token),
		"medium":   make(map[string][]models.Token),
		"low":      make(map[string][]models.Token),
	}

	// Categorize by risk level
	riskLevels := map[string]string{
		"sql_injection":          "critical",
		"command_injection":      "critical",
		"eval_usage":             "critical",
		"xss_vulnerability":      "high",
		"path_traversal":         "high",
		"unsafe_deserialization": "high",
		"auth_bypass":            "high",
		"hardcoded_password":     "medium",
		"hardcoded_api_key":      "medium",
		"weak_crypto":            "medium",
		"error_disclosure":       "low",
		"debug_info":             "low",
	}

	for category, tokens := range riskyTokens {
		if level, exists := riskLevels[category]; exists {
			result[level][category] = tokens
		} else {
			result["low"][category] = tokens
		}
	}

	return result
}

func (p *Parser) GetJavaMethodSignatures() []string {
	var signatures []string

	// Enhanced patterns for different Java method types
	javaPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?m)^\s*(public|private|protected|static|\s)*\s*([a-zA-Z_][a-zA-Z0-9_<>]*)\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\([^)]*\)\s*(\{|throws)`),
		regexp.MustCompile(`(?m)^\s*(public|private|protected)?\s*([A-Z][a-zA-Z0-9_]*)\s*\([^)]*\)\s*\{`),                                                         // Constructors
		regexp.MustCompile(`(?m)^\s*@[A-Za-z]+\s*\n\s*(public|private|protected|static|\s)*\s*([a-zA-Z_][a-zA-Z0-9_<>]*)\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\([^)]*\)`), // Annotated methods
	}

	for _, pattern := range javaPatterns {
		matches := pattern.FindAllString(p.code, -1)
		for _, match := range matches {
			// Clean up the match
			cleaned := strings.TrimSpace(strings.ReplaceAll(match, "\n", " "))
			if cleaned != "" {
				signatures = append(signatures, cleaned)
			}
		}
	}

	return signatures
}

func (p *Parser) DetectSecurityHotspots() []SecurityHotspot {
	var hotspots []SecurityHotspot

	for category, pattern := range SecurityPatterns {
		matches := pattern.FindAllStringIndex(p.code, -1)
		for _, match := range matches {
			if len(match) >= 2 {
				line, col := p.getLineAndColumn(match[0])
				hotspot := SecurityHotspot{
					Category:    category,
					Line:        line,
					Column:      col,
					Severity:    p.getSeverity(category),
					Description: p.getDescription(category),
					CodeSnippet: p.getCodeSnippet(line),
				}
				hotspots = append(hotspots, hotspot)
			}
		}
	}

	return hotspots
}

type SecurityHotspot struct {
	Category    string
	Line        int
	Column      int
	Severity    string
	Description string
	CodeSnippet string
}

func (p *Parser) getSeverity(category string) string {
	severityMap := map[string]string{
		"sql_injection":          "CRITICAL",
		"command_injection":      "CRITICAL",
		"eval_usage":             "CRITICAL",
		"xss_vulnerability":      "HIGH",
		"path_traversal":         "HIGH",
		"unsafe_deserialization": "HIGH",
		"hardcoded_password":     "MEDIUM",
		"hardcoded_api_key":      "MEDIUM",
		"weak_crypto":            "MEDIUM",
		"debug_info":             "LOW",
	}

	if severity, exists := severityMap[category]; exists {
		return severity
	}
	return "LOW"
}

func (p *Parser) getDescription(category string) string {
	descriptions := map[string]string{
		"sql_injection":          "Potential SQL injection vulnerability through string concatenation",
		"command_injection":      "Potential command injection through unsafe command execution",
		"eval_usage":             "Use of eval() function can lead to code injection",
		"xss_vulnerability":      "Potential XSS vulnerability through unsafe DOM manipulation",
		"path_traversal":         "Potential path traversal vulnerability",
		"unsafe_deserialization": "Unsafe deserialization can lead to remote code execution",
		"hardcoded_password":     "Hardcoded password found in source code",
		"hardcoded_api_key":      "Hardcoded API key found in source code",
		"weak_crypto":            "Use of weak cryptographic algorithm",
		"debug_info":             "Debug information may leak sensitive data",
	}

	if desc, exists := descriptions[category]; exists {
		return desc
	}
	return "Security issue detected"
}

func (p *Parser) getCodeSnippet(line int) string {
	lines := strings.Split(p.code, "\n")
	if line > 0 && line <= len(lines) {
		return strings.TrimSpace(lines[line-1])
	}
	return ""
}
