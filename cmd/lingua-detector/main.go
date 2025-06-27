package main

import (
	"fmt"
	"language-detection-go/internal/lingua"
)

func main() {
	detector := lingua.NewLanguageDetector()

	// Test with different languages
	testTexts := []string{
		"Hello, how are you today?",
		"Hola, ¿cómo estás hoy?",
		"Bonjour, comment allez-vous aujourd'hui?",
		"Hallo, wie geht es dir heute?",
		"Ciao, come stai oggi?",
		"こんにちは、今日はお元気ですか？",
		"안녕하세요, 오늘 어떠세요?",
		"你好，今天怎么样？",
	}

	fmt.Println("Language Detection Results:")
	fmt.Println("==========================")

	for _, text := range testTexts {
		language := detector.DetectLanguage(text)
		detectedLang, confidence := detector.DetectLanguageWithConfidence(text)

		fmt.Printf("Text: %s\n", text)
		fmt.Printf("Detected Language: %s\n", language)
		fmt.Printf("Language with Confidence: %s (%.3f)\n", detectedLang, confidence)
		fmt.Println("---")
	}
}
