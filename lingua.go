package main

/*
#cgo LDFLAGS: -L${SRCDIR}/target/release -llingua_ffi
#cgo CFLAGS: -I${SRCDIR}/target/release
#include "lingua_ffi.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// LanguageDetector wraps the Rust Lingua library
type LanguageDetector struct{}

// NewLanguageDetector creates a new language detector instance
func NewLanguageDetector() *LanguageDetector {
	return &LanguageDetector{}
}

// DetectLanguage detects the language of the given text
func (ld *LanguageDetector) DetectLanguage(text string) string {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	cResult := C.detect_language(cText)
	defer C.free_string(cResult)

	return C.GoString(cResult)
}

// DetectLanguageWithConfidence detects the language with confidence score
func (ld *LanguageDetector) DetectLanguageWithConfidence(text string) (string, float64) {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	cResult := C.detect_language_with_confidence(cText)
	defer C.free_string(cResult)

	result := C.GoString(cResult)

	// Parse the result which is in format "Language:confidence"
	var language string
	var confidence float64
	fmt.Sscanf(result, "%s:%f", &language, &confidence)

	return language, confidence
}

func main() {
	detector := NewLanguageDetector()

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
