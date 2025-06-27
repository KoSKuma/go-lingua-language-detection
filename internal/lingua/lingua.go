package lingua

/*
#cgo LDFLAGS: -L${SRCDIR}/../../rust/target/release -llingua_ffi
#cgo CFLAGS: -I${SRCDIR}/../../rust/target/release
#include "lingua_ffi.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"strings"
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

// LanguageResult represents a detected language with its confidence score
type LanguageResult struct {
	Language   string
	Confidence float64
}

// DetectMultipleLanguages detects all languages above the given confidence threshold
func (ld *LanguageDetector) DetectMultipleLanguages(text string, threshold float64) []LanguageResult {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	cResult := C.detect_multiple_languages(cText, C.double(threshold))
	defer C.free_string(cResult)

	result := C.GoString(cResult)

	if result == "no_languages_above_threshold" {
		return []LanguageResult{}
	}

	return parseLanguageResults(result)
}

// DetectTopLanguages detects the top N most likely languages
func (ld *LanguageDetector) DetectTopLanguages(text string, topN int) []LanguageResult {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))

	cResult := C.detect_top_languages(cText, C.int(topN))
	defer C.free_string(cResult)

	result := C.GoString(cResult)

	if result == "no_languages_detected" {
		return []LanguageResult{}
	}

	return parseLanguageResults(result)
}

// parseLanguageResults parses the comma-separated language:confidence string
func parseLanguageResults(result string) []LanguageResult {
	var results []LanguageResult

	parts := strings.Split(result, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Handle language names that might contain colons (like "Chinese:Traditional")
		lastColonIndex := strings.LastIndex(part, ":")
		if lastColonIndex == -1 {
			continue
		}

		language := part[:lastColonIndex]
		confidenceStr := part[lastColonIndex+1:]

		var confidence float64
		if _, err := fmt.Sscanf(confidenceStr, "%f", &confidence); err == nil {
			results = append(results, LanguageResult{
				Language:   language,
				Confidence: confidence,
			})
		}
	}

	return results
}
