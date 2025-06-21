package services

import (
	"github.com/PratikforCoding/CodeSentry/internal/models"
	"github.com/PratikforCoding/CodeSentry/internal/utils"
	"strings"
)

type LanguageDetector struct{}

func NewLanguageDetector() *LanguageDetector {
	return &LanguageDetector{}
}

func (ld *LanguageDetector) DetectLanguage(code string) models.Language {
	parser := utils.NewParser(code)
	_ = parser.Tokenize()

	scores := make(map[models.Language]int)

	if strings.Contains(code, "package main") || strings.Contains(code, "func main()") {
		scores[models.Go] += 15
	}

	if strings.Contains(code, "function ") || strings.Contains(code, "const ") || strings.Contains(code, "let ") {
		scores[models.JavaScript] += 10
	}

	if strings.Contains(code, "def ") || strings.Contains(code, "import ") {
		scores[models.Python] += 10
	}

	if strings.Contains(code, "public class") || strings.Contains(code, "public static void main") {
		scores[models.Java] += 15
	}

	if strings.Contains(code, "SELECT") || strings.Contains(code, "INSERT") {
		scores[models.SQL] += 15
	}

	languages := []models.Language{models.Go, models.JavaScript, models.Python, models.Java}
	for _, lang := range languages {
		tokenScore := parser.GetLanguageScore(lang)
		scores[lang] += tokenScore
	}

	maxScore := 0
	detectedLang := models.Unknown

	for lang, score := range scores {
		if score > maxScore {
			maxScore = score
			detectedLang = lang
		}
	}

	return detectedLang
}

func (ld *LanguageDetector) buildKeywordLanguageMap() map[string][]models.Language {
	keywordMap := make(map[string][]models.Language)
	keywords := utils.GetLanguageKeywords()

	for lang, langKeywords := range keywords {
		for _, keyword := range langKeywords {
			keywordMap[keyword] = append(keywordMap[keyword], lang)
		}
	}
	return keywordMap
}

func (ld *LanguageDetector) cleanWord(word string) string {
	cleaned := strings.ToLower(strings.Trim(word, ".,;:(){}[]\"'<>!?"))
	return cleaned
}

func (ld *LanguageDetector) getBestLanguage(scores map[models.Language]int) models.Language {
	var bestLang models.Language
	maxScore := 0

	for lang, score := range scores {
		if score > maxScore {
			maxScore = score
			bestLang = lang
		}
	}
	return bestLang
}
