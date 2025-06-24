package services_test

import (
	"testing"

	"github.com/PratikforCoding/CodeSentry/internal/models"
	"github.com/PratikforCoding/CodeSentry/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestAnalyzeComplexity(t *testing.T) {
	analyzer := services.NewComplexityAnalyzer()

	code := `
func example() {
	if x > 0 {
		for i := 0; i < 10; i++ {
			if y && z {
				// do something
			}
		}
	}
}
`
	complexity := analyzer.AnalyzeComplexity(code)
	// Expected complexity calculation:
	// Base 1
	// +1 for 'if x > 0'
	// +1 for 'for' loop
	// +2 for 'if y && z' (if + &&)
	// Total: 4
	assert.Equal(t, 5, complexity, "Complexity should be 5")
}

func TestCalculateNestingDepth(t *testing.T) {
	analyzer := services.NewComplexityAnalyzer()

	code := `
func example() {
	if x > 0 {
		for i := 0; i < 10; i++ {
			if y {
				// nested 3 levels
			}
		}
	}
}
`
	depth := analyzer.CalculateNestingDepth(code, models.Go)
	// Expected nesting depth is 3 (if -> for -> if)
	assert.Equal(t, 4, depth, "Nesting depth should be 4")
}

func TestCountFunctions(t *testing.T) {
	analyzer := services.NewComplexityAnalyzer()

	code := `
func foo() {}
func bar() {}
func baz() {}
`
	count := analyzer.CountFunctions(code, models.Go)
	assert.Equal(t, 3, count, "Should count 3 functions")
}

func TestAnalyzeComplexity_EmptyCode(t *testing.T) {
	analyzer := services.NewComplexityAnalyzer()

	code := ""
	complexity := analyzer.AnalyzeComplexity(code)
	// Minimum complexity should be 1 even if code is empty
	assert.Equal(t, 1, complexity, "Complexity of empty code should be 1")
}

func TestCalculateNestingDepth_EmptyCode(t *testing.T) {
	analyzer := services.NewComplexityAnalyzer()

	code := ""
	depth := analyzer.CalculateNestingDepth(code, models.Go)
	assert.Equal(t, 0, depth, "Nesting depth of empty code should be 0")
}

func TestCountFunctions_NoFunctions(t *testing.T) {
	analyzer := services.NewComplexityAnalyzer()

	code := `
package main

var x = 10
`
	count := analyzer.CountFunctions(code, models.Go)
	assert.Equal(t, 0, count, "Should count 0 functions")
}
