package services

import (
	"github.com/PratikforCoding/CodeSentry/internal/models"
	"github.com/PratikforCoding/CodeSentry/internal/utils"
)

type ComplexityAnalyzer struct{}

func NewComplexityAnalyzer() *ComplexityAnalyzer {
	return &ComplexityAnalyzer{}
}

func (ca *ComplexityAnalyzer) AnalyzeComplexity(code string) int {
	parser := utils.NewParser(code)
	_ = parser.Tokenize()
	complexity := 1

	complexityTokens := parser.GetComplexityTokens()

	for _, token := range complexityTokens {
		switch token.Value {
		case "if", "elif", "else if":
			complexity++
		case "for", "while", "do":
			complexity++
		case "switch":
			complexity++
		case "case":
			complexity++
		case "try", "catch", "except", "finally":
			complexity++
		case "&&", "||":
			complexity++
		case "?":
			complexity++
		}
	}

	return complexity
}

func (ca *ComplexityAnalyzer) CalculateNestingDepth(code string, language models.Language) int {
	parser := utils.NewParser(code)
	parser.Tokenize()
	return parser.AnalyzeNestingDepth(language)
}

func (ca *ComplexityAnalyzer) CountFunctions(code string, language models.Language) int {
	parser := utils.NewParser(code)
	parser.Tokenize()

	functionTokens := parser.GetFunctionTokens(language)
	return len(functionTokens)
}
