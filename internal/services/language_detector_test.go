package services_test

import (
	"testing"

	"github.com/PratikforCoding/CodeSentry/internal/models"
	"github.com/PratikforCoding/CodeSentry/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestDetectLanguage_Go(t *testing.T) {
	ld := services.NewLanguageDetector()

	code := `
		package main

		import "fmt"

		func main() {
			fmt.Println("Hello, Go!")
		}
	`

	lang := ld.DetectLanguage(code)
	assert.Equal(t, models.Go, lang, "Should detect Go language")
}

func TestDetectLanguage_JavaScript(t *testing.T) {
	ld := services.NewLanguageDetector()

	code := `
		function greet() {
			console.log("Hello, JavaScript!");
		}
		const x = 10;
	`

	lang := ld.DetectLanguage(code)
	assert.Equal(t, models.JavaScript, lang, "Should detect JavaScript language")
}

func TestDetectLanguage_Python(t *testing.T) {
	ld := services.NewLanguageDetector()

	code := `
		def greet():
			print("Hello, Python!")

		import sys
	`

	lang := ld.DetectLanguage(code)
	assert.Equal(t, models.Python, lang, "Should detect Python language")
}

func TestDetectLanguage_Java(t *testing.T) {
	ld := services.NewLanguageDetector()

	code := `
		public class HelloWorld {
			public static void main(String[] args) {
				System.out.println("Hello, Java!");
			}
		}
	`

	lang := ld.DetectLanguage(code)
	assert.Equal(t, models.Java, lang, "Should detect Java language")
}

func TestDetectLanguage_SQL(t *testing.T) {
	ld := services.NewLanguageDetector()

	code := `
		SELECT * FROM users WHERE id = 1;
		INSERT INTO users (name, age) VALUES ('Alice', 30);
	`

	lang := ld.DetectLanguage(code)
	assert.Equal(t, models.SQL, lang, "Should detect SQL language")
}

func TestDetectLanguage_Unknown(t *testing.T) {
	ld := services.NewLanguageDetector()

	code := `
		This is some random text without code keywords.
	`

	lang := ld.DetectLanguage(code)
	assert.Equal(t, models.Unknown, lang, "Should detect Unknown language")
}
