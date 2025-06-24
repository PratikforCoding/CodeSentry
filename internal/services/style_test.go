package services_test

import (
	"testing"

	"github.com/PratikforCoding/CodeSentry/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestAnalyzeStyle_MaxLineLength(t *testing.T) {
	analyzer := services.NewStyleAnalyzer()

	longLine := ""
	for i := 0; i < 130; i++ {
		longLine += "a"
	}

	code := longLine + "\nshort line"

	suggestions := analyzer.AnalyzeStyle(code)

	assert.Len(t, suggestions, 1)
	assert.Equal(t, 1, suggestions[0].Line)
	assert.Equal(t, 120, suggestions[0].Column)
	assert.Equal(t, "Line exceeds maximum length of 120 characters", suggestions[0].Message)
	assert.Equal(t, "max-line-length", suggestions[0].Rule)
	assert.Equal(t, "WARNING", suggestions[0].Severity)
}

func TestAnalyzeStyle_TrailingWhitespace(t *testing.T) {
	analyzer := services.NewStyleAnalyzer()

	code := "line with trailing space \nclean line"

	suggestions := analyzer.AnalyzeStyle(code)

	assert.Len(t, suggestions, 1)
	assert.Equal(t, 1, suggestions[0].Line)
	assert.Equal(t, len([]string{"line with trailing space ", "clean line"}), suggestions[0].Column) // as per your code
	assert.Equal(t, "Trailing whitespace detected", suggestions[0].Message)
	assert.Equal(t, "no-trailing-whitespace", suggestions[0].Rule)
	assert.Equal(t, "INFO", suggestions[0].Severity)
}

func TestAnalyzeStyle_MixedIndentation(t *testing.T) {
	analyzer := services.NewStyleAnalyzer()

	code := "\t  mixed indentation line\nclean line"

	suggestions := analyzer.AnalyzeStyle(code)

	assert.Len(t, suggestions, 1)
	assert.Equal(t, 1, suggestions[0].Line)
	assert.Equal(t, 0, suggestions[0].Column)
	assert.Equal(t, "Mixed spaces and tabs for indentation", suggestions[0].Message)
	assert.Equal(t, "consistent-indentation", suggestions[0].Rule)
	assert.Equal(t, "WARNING", suggestions[0].Severity)
}

func TestAnalyzeStyle_NamingConvention(t *testing.T) {
	analyzer := services.NewStyleAnalyzer()

	code := "var someVariable = 1\nvar some_variable = 2"

	suggestions := analyzer.AnalyzeStyle(code)

	// Should detect naming convention issues on both lines
	assert.Len(t, suggestions, 2)
	for _, suggestion := range suggestions {
		assert.Equal(t, "Mixed naming convention detected", suggestion.Message)
		assert.Equal(t, "naming-convention", suggestion.Rule)
		assert.Equal(t, "INFO", suggestion.Severity)
	}
}

func TestAnalyzeStyle_NoIssues(t *testing.T) {
	analyzer := services.NewStyleAnalyzer()

	code := "cleanLine := 1\nanotherCleanLine := 2"

	suggestions := analyzer.AnalyzeStyle(code)

	assert.Empty(t, suggestions)
}
