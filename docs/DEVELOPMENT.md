# Development Guide

## Quick Start for Development

### Prerequisites Check
```bash
# Check if Rust is available
cargo --version

# Check if Go is available
go version

# Check if cbindgen is installed
cbindgen --version
```

### Initial Setup
```bash
# Set up Rust environment
source $HOME/.cargo/env

# Install cbindgen if not already installed
cargo install cbindgen

# Build the project
make build
```

## Development Workflow

### 1. Adding New Languages

**Step 1: Update Rust FFI**
Edit `src/lib.rs` and add new languages to the `languages` vector:

```rust
let languages = vec![
    Language::English,
    Language::Spanish,
    // ... existing languages ...
    Language::Indonesian,
    Language::Malay,
    Language::Thai,
    Language::Vietnamese,
    Language::Tagalog,
    // Add your new language here
    Language::YourNewLanguage,
];
```

**Step 2: Rebuild**
```bash
cargo build --release
cbindgen --config cbindgen.toml --crate lingua-ffi --output target/release/lingua_ffi.h
```

**Step 3: Add Tests**
Add test cases to `lingua_test.go`:

```go
func TestLanguageDetection_NewLanguage(t *testing.T) {
    detector := NewLanguageDetector()
    tests := []struct {
        text     string
        expected string
    }{
        {"Your test sentence in the new language", "YourNewLanguage"},
    }
    // ... test implementation
}
```

### 2. Adding New FFI Functions

**Step 1: Add Rust Function**
In `src/lib.rs`:

```rust
#[no_mangle]
pub extern "C" fn detect_languages_with_scores(text: *const c_char) -> *mut c_char {
    let c_str = unsafe {
        assert!(!text.is_null());
        CStr::from_ptr(text)
    };
    
    let text_str = match c_str.to_str() {
        Ok(s) => s,
        Err(_) => return CString::new("error: invalid utf8").unwrap().into_raw(),
    };

    let languages = vec![/* your languages */];
    let detector = LanguageDetectorBuilder::from_languages(&languages).build();
    
    let results = detector.compute_language_confidence_values(text_str);
    
    // Format results as JSON or custom format
    let result = format!("{:?}", results);
    CString::new(result).unwrap().into_raw()
}
```

**Step 2: Regenerate Header**
```bash
cbindgen --config cbindgen.toml --crate lingua-ffi --output target/release/lingua_ffi.h
```

**Step 3: Add Go Wrapper**
In `lingua.go`:

```go
// DetectLanguagesWithScores returns all detected languages with confidence scores
func (ld *LanguageDetector) DetectLanguagesWithScores(text string) string {
    cText := C.CString(text)
    defer C.free(unsafe.Pointer(cText))
    
    cResult := C.detect_languages_with_scores(cText)
    defer C.free_string(cResult)
    
    return C.GoString(cResult)
}
```

### 3. Performance Optimization

**Rust Side:**
- Use `LanguageDetectorBuilder::from_languages()` with only needed languages
- Consider caching the detector instance
- Use `compute_language_confidence_values()` for better accuracy

**Go Side:**
- Reuse `LanguageDetector` instances
- Use goroutines for batch processing (but be careful with CGO)
- Profile with `go test -bench=.`

### 4. Error Handling Patterns

**Rust FFI Error Handling:**
```rust
#[no_mangle]
pub extern "C" fn safe_detect_language(text: *const c_char) -> *mut c_char {
    if text.is_null() {
        return CString::new("error: null pointer").unwrap().into_raw();
    }
    
    let c_str = unsafe { CStr::from_ptr(text) };
    let text_str = match c_str.to_str() {
        Ok(s) => s,
        Err(_) => return CString::new("error: invalid utf8").unwrap().into_raw(),
    };
    
    if text_str.is_empty() {
        return CString::new("error: empty text").unwrap().into_raw();
    }
    
    // ... rest of implementation
}
```

**Go Error Handling:**
```go
func (ld *LanguageDetector) SafeDetectLanguage(text string) (string, error) {
    if text == "" {
        return "", errors.New("empty text")
    }
    
    cText := C.CString(text)
    defer C.free(unsafe.Pointer(cText))
    
    cResult := C.safe_detect_language(cText)
    defer C.free_string(cResult)
    
    result := C.GoString(cResult)
    if strings.HasPrefix(result, "error:") {
        return "", errors.New(result)
    }
    
    return result, nil
}
```

## Testing Strategies

### 1. Unit Tests
```go
func TestLanguageDetector_EdgeCases(t *testing.T) {
    detector := NewLanguageDetector()
    
    // Test empty string
    result := detector.DetectLanguage("")
    if result != "unknown" {
        t.Errorf("Expected 'unknown' for empty string, got %s", result)
    }
    
    // Test very short text
    result = detector.DetectLanguage("a")
    // Short text might be ambiguous, so just check it doesn't panic
    _ = result
}
```

### 2. Benchmark Tests
```go
func BenchmarkLanguageDetection(b *testing.B) {
    detector := NewLanguageDetector()
    text := "This is a benchmark test for language detection performance."
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        detector.DetectLanguage(text)
    }
}
```

### 3. Integration Tests
```go
func TestLanguageDetection_Integration(t *testing.T) {
    detector := NewLanguageDetector()
    
    // Test with real-world text samples
    samples := map[string]string{
        "The quick brown fox jumps over the lazy dog.": "English",
        "El rápido zorro marrón salta sobre el perro perezoso.": "Spanish",
        // Add more real-world examples
    }
    
    for text, expected := range samples {
        result := detector.DetectLanguage(text)
        if result != expected {
            t.Errorf("Expected %s for '%s', got %s", expected, text, result)
        }
    }
}
```

## Troubleshooting

### Common Build Issues

**1. Header Not Found**
```bash
# Error: 'lingua_ffi.h' file not found
# Solution: Regenerate the header
cbindgen --config cbindgen.toml --crate lingua-ffi --output target/release/lingua_ffi.h
```

**2. Library Not Found**
```bash
# Error: cannot find -llingua_ffi
# Solution: Ensure library exists
ls -la target/release/liblingua_ffi.*
```

**3. CGO Compilation Errors**
```bash
# Error: cstdarg not found
# Solution: Ensure cbindgen generates C headers, not C++
# Check cbindgen.toml has: language = "C", cpp_compat = false
```

### Debugging Tips

**1. Check Generated Header**
```bash
# View the generated header to understand the interface
cat target/release/lingua_ffi.h
```

**2. Verify Library Symbols**
```bash
# On macOS
nm -D target/release/liblingua_ffi.dylib | grep detect

# On Linux
nm -D target/release/liblingua_ffi.so | grep detect
```

**3. Test Rust Functions Directly**
```bash
# Create a simple Rust test
cargo test
```

## Best Practices

### 1. Memory Management
- Always use `defer C.free_string()` in Go
- Check for null pointers in Rust
- Use `CString::into_raw()` and `from_raw()` properly

### 2. Error Handling
- Return meaningful error messages from Rust
- Handle UTF-8 conversion errors
- Validate input parameters

### 3. Performance
- Minimize CGO calls (they have overhead)
- Reuse detector instances
- Consider batch processing for multiple texts

### 4. Testing
- Test edge cases (empty strings, null pointers)
- Test with real-world text samples
- Include performance benchmarks
- Test mixed-language scenarios

## Future Enhancements

### 1. Concurrency Support
```go
// Consider using a pool of detectors for concurrent access
type LanguageDetectorPool struct {
    detectors chan *LanguageDetector
}

func NewLanguageDetectorPool(size int) *LanguageDetectorPool {
    pool := &LanguageDetectorPool{
        detectors: make(chan *LanguageDetector, size),
    }
    for i := 0; i < size; i++ {
        pool.detectors <- NewLanguageDetector()
    }
    return pool
}
```

### 2. Configuration Support
```go
type DetectorConfig struct {
    Languages []string
    MinConfidence float64
    MaxResults int
}
```

### 3. Batch Processing
```go
func (ld *LanguageDetector) DetectLanguagesBatch(texts []string) []string {
    results := make([]string, len(texts))
    for i, text := range texts {
        results[i] = ld.DetectLanguage(text)
    }
    return results
}
```

This development guide should help you continue building and improving the language detection system effectively. 
