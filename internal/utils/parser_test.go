package utils_test

import (
	"github.com/PratikforCoding/CodeSentry/internal/models"
	"github.com/PratikforCoding/CodeSentry/internal/utils"
	"testing"
)

func TestParser_TokenizeBasic(t *testing.T) {
	code := `package main
// This is a comment
func main() {
    var x = 42
    var s = "hello"
    if x > 10 {
        // do something
    }
}`
	parser := utils.NewParser(code)
	tokens := parser.Tokenize()

	if len(tokens) == 0 {
		t.Fatal("Expected tokens, got none")
	}

	// Check some expected token types and values
	foundFunc := false
	foundIf := false
	foundLiteral := false
	foundComment := false

	for _, token := range tokens {
		switch token.Value {
		case "func":
			if token.Type != models.KEYWORD {
				t.Errorf("Expected 'func' to be KEYWORD, got %v", token.Type)
			}
			foundFunc = true
		case "if":
			if token.Type != models.KEYWORD {
				t.Errorf("Expected 'if' to be KEYWORD, got %v", token.Type)
			}
			foundIf = true
		case "\"hello\"":
			if token.Type != models.LITERAL {
				t.Errorf("Expected string literal token, got %v", token.Type)
			}
			foundLiteral = true
		}
		if token.Type == models.COMMENT {
			foundComment = true
		}
	}

	if !foundFunc {
		t.Error("Did not find 'func' keyword token")
	}
	if !foundIf {
		t.Error("Did not find 'if' keyword token")
	}
	if !foundLiteral {
		t.Error("Did not find string literal token")
	}
	if !foundComment {
		t.Error("Did not find comment token")
	}
}

func TestCountLines(t *testing.T) {
	code := `
// Comment line

func main() {
    // Another comment
    println("Hello")
}
/* Block comment
   continues here
*/
`
	total, blank, comment := utils.CountLines(code)

	if total != 10 {
		t.Errorf("Expected total lines 9, got %d", total)
	}
	if blank != 2 {
		t.Errorf("Expected blank lines 2, got %d", blank)
	}
	if comment != 5 {
		t.Errorf("Expected comment lines 4, got %d", comment)
	}
}

func TestGetTokensByType(t *testing.T) {
	code := `if (x == 10) { return x }`
	parser := utils.NewParser(code)
	parser.Tokenize()

	keywords := parser.GetTokensByType(models.KEYWORD)
	if len(keywords) == 0 {
		t.Error("Expected to find keyword tokens")
	}

	operators := parser.GetTokensByType(models.OPERATOR)
	if len(operators) == 0 {
		t.Error("Expected to find operator tokens")
	}

	delimiters := parser.GetTokensByType(models.DELIMITER)
	if len(delimiters) == 0 {
		t.Error("Expected to find delimiter tokens")
	}
}

func TestCountKeyword(t *testing.T) {
	code := `if (x == 10) { if (y == 20) { return x } }`
	parser := utils.NewParser(code)
	parser.Tokenize()

	count := parser.CountKeyword("if")
	if count != 2 {
		t.Errorf("Expected 2 'if' keywords, got %d", count)
	}
}

func TestIsRiskyContext(t *testing.T) {
	code := `eval(userInput)`
	parser := utils.NewParser(code)
	parser.Tokenize()

	tokens := parser.GetTokensByValue("eval")
	if len(tokens) == 0 {
		t.Fatal("Expected to find 'eval' token")
	}

	risky := parser.IsRiskyContext(tokens[0], "eval")
	if !risky {
		t.Error("Expected 'eval' token to be risky context")
	}
}
