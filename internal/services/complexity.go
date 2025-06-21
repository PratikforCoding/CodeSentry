package services

import (
	"fmt"
	"github.com/PratikforCoding/CodeSentry/internal/utils"
)

type ComplexityAnalyzer struct{}

func NewComplexityAnalyzer() *ComplexityAnalyzer {
	return &ComplexityAnalyzer{}
}

func (ca *ComplexityAnalyzer) AnalyzeComplexity(code string) int {
	fmt.Println("🔍 DEBUG: Creating parser...")
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

func (ca *ComplexityAnalyzer) CalculateNestingDepth(code string) int {
	parser := utils.NewParser(code)
	parser.Tokenize()
	return parser.AnalyzeNestingDepth()
}

func (ca *ComplexityAnalyzer) CountFunctions(code string) int {
	parser := utils.NewParser(code)
	parser.Tokenize()

	functionTokens := parser.GetFunctionTokens()
	return len(functionTokens)
}
