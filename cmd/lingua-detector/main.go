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
		// Thai + English mixed
		"สวัสดีครับ, good morning! How are you today?",
		// Thai + Myanmar mixed
		"สวัสดีครับ, မင်္ဂလာပါ။ နေကောင်းပါသလား။",
		// Pure Japanese
		"おはようございます。今日は良い天気ですね。",
		// Pure Indonesian
		"Selamat pagi! Bagaimana kabar Anda hari ini?",
		// Pure Malaysian
		"Selamat pagi! Bagaimana keadaan anda hari ini?",
		// Mixed Indonesian + English
		"Apa kabar? How are you doing today?",
		// Mixed Malaysian + English
		"Bagaimana keadaan? Are you feeling well?",
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

	// Test multiple language detection
	fmt.Println("\nMultiple Language Detection Results:")
	fmt.Println("====================================")

	mixedTexts := []string{
		"สวัสดีครับ, good morning! How are you today?", // Thai + English
		"สวัสดีครับ, မင်္ဂလာပါ။ နေကောင်းပါသလား။",       // Thai + Myanmar
		"Apa kabar? How are you doing today?",          // Indonesian + English
		"Hola, 你好! ¿Cómo estás?",                       // Spanish + Chinese + Spanish
		"Guten Tag, こんにちは! Wie geht es dir?",           // German + Japanese + German
	}

	for _, text := range mixedTexts {
		fmt.Printf("\nText: %s\n", text)

		// Test threshold-based detection (0.1 threshold)
		thresholdResults := detector.DetectMultipleLanguages(text, 0.1)
		fmt.Printf("Languages above 0.1 threshold: ")
		if len(thresholdResults) == 0 {
			fmt.Printf("None\n")
		} else {
			for i, result := range thresholdResults {
				if i > 0 {
					fmt.Printf(", ")
				}
				fmt.Printf("%s (%.3f)", result.Language, result.Confidence)
			}
			fmt.Printf("\n")
		}

		// Test top 3 languages
		topResults := detector.DetectTopLanguages(text, 3)
		fmt.Printf("Top 3 languages: ")
		for i, result := range topResults {
			if i > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("%s (%.3f)", result.Language, result.Confidence)
		}
		fmt.Printf("\n")
		fmt.Println("---")
	}
}
