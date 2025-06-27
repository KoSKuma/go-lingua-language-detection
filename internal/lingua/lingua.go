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
