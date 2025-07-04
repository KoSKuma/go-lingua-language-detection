package lingua

import (
	"fmt"
	"language-detection-go/internal/lingua"
	"testing"
)

func TestLanguageDetection(t *testing.T) {
	detector := lingua.NewLanguageDetector()

	tests := []struct {
		text     string
		expected string
	}{
		{"Hello, how are you?", "English"},
		{"Hola, ¿cómo estás?", "Spanish"},
		{"Bonjour, comment allez-vous?", "French"},
		{"Hallo, wie geht es dir?", "German"},
		{"Ciao, come stai?", "Portuguese"},
	}

	for _, test := range tests {
		t.Run(test.text, func(t *testing.T) {
			result := detector.DetectLanguage(test.text)
			if result != test.expected {
				t.Errorf("Expected %s, got %s for text: %s", test.expected, result, test.text)
			}
		})
	}
}

func TestLanguageDetectionWithConfidence(t *testing.T) {
	detector := lingua.NewLanguageDetector()

	text := "Hello, world!"
	language, confidence := detector.DetectLanguageWithConfidence(text)

	if language == "" {
		t.Error("Expected non-empty language")
	}

	if confidence < 0.0 || confidence > 1.0 {
		t.Errorf("Expected confidence between 0.0 and 1.0, got %f", confidence)
	}

	fmt.Printf("Detected: %s with confidence: %.3f\n", language, confidence)
}

func TestDetectMultipleLanguages(t *testing.T) {
	detector := lingua.NewLanguageDetector()

	// Test with mixed language text
	text := "สวัสดีครับ, good morning! How are you today?"
	results := detector.DetectMultipleLanguages(text, 0.1)

	if len(results) == 0 {
		t.Error("Expected at least one language above threshold")
	}

	// Check that all results have confidence >= threshold
	for _, result := range results {
		if result.Confidence < 0.1 {
			t.Errorf("Expected confidence >= 0.1, got %.3f for %s", result.Confidence, result.Language)
		}
	}

	fmt.Printf("Multiple languages detected: ")
	for i, result := range results {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("%s (%.3f)", result.Language, result.Confidence)
	}
	fmt.Printf("\n")
}

func TestDetectTopLanguages(t *testing.T) {
	detector := lingua.NewLanguageDetector()

	// Test with mixed language text
	text := "Apa kabar? How are you doing today?"
	results := detector.DetectTopLanguages(text, 3)

	if len(results) == 0 {
		t.Error("Expected at least one language")
	}

	if len(results) > 3 {
		t.Errorf("Expected at most 3 languages, got %d", len(results))
	}

	// Check that results are sorted by confidence (descending)
	for i := 1; i < len(results); i++ {
		if results[i-1].Confidence < results[i].Confidence {
			t.Errorf("Expected results to be sorted by confidence (descending)")
		}
	}

	fmt.Printf("Top 3 languages: ")
	for i, result := range results {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("%s (%.3f)", result.Language, result.Confidence)
	}
	fmt.Printf("\n")
}

// Example usage
func ExampleLanguageDetector_DetectLanguage() {
	detector := lingua.NewLanguageDetector()

	// Detect language
	language := detector.DetectLanguage("Hello, world!")
	fmt.Println(language)
	// Output: English
}

func ExampleLanguageDetector_DetectLanguageWithConfidence() {
	detector := lingua.NewLanguageDetector()

	// Detect language with confidence
	language, confidence := detector.DetectLanguageWithConfidence("Hola, mundo!")
	fmt.Printf("Language: %s, Confidence: %.3f\n", language, confidence)
	// Output: Language: Spanish:0.325, Confidence: 0.000
}

func ExampleLanguageDetector_DetectMultipleLanguages() {
	detector := lingua.NewLanguageDetector()

	// Detect multiple languages above threshold
	results := detector.DetectMultipleLanguages("สวัสดีครับ, good morning!", 0.1)
	for _, result := range results {
		fmt.Printf("%s (%.3f) ", result.Language, result.Confidence)
	}
	fmt.Println()
	// Output: English (0.707) Tagalog (0.142)
}

func ExampleLanguageDetector_DetectTopLanguages() {
	detector := lingua.NewLanguageDetector()

	// Detect top 3 languages
	results := detector.DetectTopLanguages("Apa kabar? How are you?", 3)
	for _, result := range results {
		fmt.Printf("%s (%.3f) ", result.Language, result.Confidence)
	}
	fmt.Println()
	// Output: Tagalog (0.204) English (0.200) Indonesian (0.182)
}

func TestLanguageDetection_SEA(t *testing.T) {
	detector := lingua.NewLanguageDetector()
	tests := []struct {
		text     string
		expected string
	}{
		{"Apa kabar? Semoga harimu menyenangkan.", "Indonesian"},
		{"Selamat pagi, bagaimana keadaanmu?", "Malay"},
		{"สวัสดีครับ วันนี้เป็นอย่างไรบ้าง", "Thai"},
		{"Chào bạn, hôm nay bạn thế nào?", "Vietnamese"},
		{"Kamusta ka? Sana ay maganda ang araw mo.", "Tagalog"},
	}
	for _, test := range tests {
		t.Run(test.text, func(t *testing.T) {
			result := detector.DetectLanguage(test.text)
			if result != test.expected {
				t.Errorf("Expected %s, got %s for text: %s", test.expected, result, test.text)
			}
		})
	}
}

func TestLanguageDetection_Mixed(t *testing.T) {
	detector := lingua.NewLanguageDetector()
	tests := []struct {
		text string
	}{
		{"Hello, apa kabar? How are you today?"},  // English + Indonesian
		{"Bonjour, selamat pagi! Comment ça va?"}, // French + Malay
		{"สวัสดีครับ, good morning!"},             // Thai + English
		{"Hola, 你好! ¿Cómo estás?"},                // Spanish + Chinese + Spanish
		{"Guten Tag, こんにちは! Wie geht es dir?"},    // German + Japanese + German
	}
	for _, test := range tests {
		t.Run(test.text, func(t *testing.T) {
			lang, conf := detector.DetectLanguageWithConfidence(test.text)
			if lang == "" || conf < 0.0 || conf > 1.0 {
				t.Errorf("Unexpected result for mixed input: %s (%.3f)", lang, conf)
			}
			fmt.Printf("Mixed input: %s -> %s (%.3f)\n", test.text, lang, conf)
		})
	}
}
