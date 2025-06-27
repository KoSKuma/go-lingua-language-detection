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
