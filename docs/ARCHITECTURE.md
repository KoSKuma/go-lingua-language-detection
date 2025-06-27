# Architecture Documentation

## System Overview

The Language Detection system is a hybrid Rust-Go application that leverages the Lingua library for accurate language detection. The architecture uses Foreign Function Interface (FFI) to bridge Rust's performance and Go's ease of use.

## High-Level Architecture

```
┌─────────────────┐    CGO Interface    ┌─────────────────┐
│   Go Layer      │ ◄─────────────────► │   Rust Layer    │
│                 │                     │                 │
│ • lingua.go     │                     │ • src/lib.rs    │
│ • lingua_test.go│                     │ • Cargo.toml    │
│ • go.mod        │                     │ • build.rs      │
└─────────────────┘                     └─────────────────┘
         │                                       │
         │                                       │
         ▼                                       ▼
┌─────────────────┐                     ┌─────────────────┐
│   Build System  │                     │   Lingua Library│
│                 │                     │                 │
│ • Makefile      │                     │ • Language Models│
│ • cbindgen      │                     │ • Detection Algo│
└─────────────────┘                     └─────────────────┘
```

## Component Details

### 1. Rust FFI Layer (`src/lib.rs`)

**Purpose**: Wraps Lingua functionality in C-compatible functions for Go consumption.

**Key Components**:
- **Language Vector**: Defines supported languages for detection
- **FFI Functions**: Exported C-compatible functions
- **Memory Management**: Handles string conversion and cleanup

**Functions**:
```rust
// Primary detection function
detect_language(text: *const c_char) -> *mut c_char

// Detection with confidence scores
detect_language_with_confidence(text: *const c_char) -> *mut c_char

// Memory cleanup
free_string(ptr: *mut c_char)
```

**Data Flow**:
1. Receive C string pointer from Go
2. Convert to Rust string with UTF-8 validation
3. Create Lingua detector with configured languages
4. Perform language detection
5. Format result as C string
6. Return pointer to Go

### 2. Go Interface Layer (`lingua.go`)

**Purpose**: Provides a clean Go API that wraps the Rust FFI calls.

**Key Components**:
- **LanguageDetector Struct**: Main interface for language detection
- **CGO Integration**: Uses `import "C"` to call Rust functions
- **Memory Management**: Ensures proper cleanup with `defer` statements

**Methods**:
```go
// Basic language detection
DetectLanguage(text string) string

// Detection with confidence
DetectLanguageWithConfidence(text string) (string, float64)
```

**Data Flow**:
1. Convert Go string to C string
2. Call Rust FFI function
3. Convert result back to Go string
4. Clean up C memory allocations

### 3. Build System

**Components**:
- **Cargo.toml**: Rust dependencies and build configuration
- **build.rs**: Build script for cbindgen header generation
- **cbindgen.toml**: Configuration for C header generation
- **Makefile**: Orchestrates the complete build process

**Build Process**:
1. Cargo builds Rust library (`liblingua_ffi.dylib/.so`)
2. cbindgen generates C header (`lingua_ffi.h`)
3. Go compiler uses CGO to link against Rust library
4. Final Go binary includes Rust library

## Data Flow Diagrams

### Language Detection Flow

```
Go Application
     │
     ▼
┌─────────────┐
│ DetectLanguage │
└─────────────┘
     │
     ▼
┌─────────────┐    CGO Call    ┌─────────────┐
│ Go String   │ ─────────────► │ C String    │
└─────────────┘                └─────────────┘
     │                                │
     │                                ▼
     │                        ┌─────────────┐
     │                        │ Rust String │
     │                        └─────────────┘
     │                                │
     │                                ▼
     │                        ┌─────────────┐
     │                        │ Lingua      │
     │                        │ Detector    │
     │                        └─────────────┘
     │                                │
     │                                ▼
     │                        ┌─────────────┐
     │                        │ Language    │
     │                        │ Result      │
     │                        └─────────────┘
     │                                │
     ▼                                ▼
┌─────────────┐    CGO Return   ┌─────────────┐
│ Go Result   │ ◄────────────── │ C String    │
└─────────────┘                └─────────────┘
```

### Memory Management Flow

```
Go Function Call
     │
     ▼
┌─────────────┐
│ C.CString() │ ← Allocate C memory
└─────────────┘
     │
     ▼
┌─────────────┐
│ defer C.free() │ ← Schedule cleanup
└─────────────┘
     │
     ▼
┌─────────────┐
│ Rust FFI    │ ← Use C string
└─────────────┘
     │
     ▼
┌─────────────┐
│ C.GoString() │ ← Convert result
└─────────────┘
     │
     ▼
┌─────────────┐
│ defer C.free_string() │ ← Cleanup result
└─────────────┘
     │
     ▼
┌─────────────┐
│ Return      │ ← Function returns
└─────────────┘
     │
     ▼
┌─────────────┐
│ Cleanup     │ ← defer statements execute
└─────────────┘
```

## Error Handling Architecture

### Rust Side Error Handling
```rust
// Input validation
if text.is_null() {
    return CString::new("error: null pointer").unwrap().into_raw();
}

// UTF-8 validation
let text_str = match c_str.to_str() {
    Ok(s) => s,
    Err(_) => return CString::new("error: invalid utf8").unwrap().into_raw(),
};

// Empty string handling
if text_str.is_empty() {
    return CString::new("error: empty text").unwrap().into_raw();
}
```

### Go Side Error Handling
```go
// Input validation
if text == "" {
    return "", errors.New("empty text")
}

// Result validation
result := C.GoString(cResult)
if strings.HasPrefix(result, "error:") {
    return "", errors.New(result)
}
```

## Performance Characteristics

### CGO Overhead
- Each CGO call has overhead (~50-100ns)
- String conversion between Go and C adds latency
- Memory allocation/deallocation for each call

### Optimization Strategies
1. **Batch Processing**: Process multiple texts in single Rust call
2. **Instance Reuse**: Reuse LanguageDetector instances
3. **Memory Pooling**: Pool C string allocations
4. **Concurrent Access**: Use goroutines with care (CGO is not goroutine-safe)

### Memory Usage
- Rust library: ~93MB (includes all language models)
- Per detection: ~1-10KB temporary allocations
- Go runtime: Minimal overhead

## Security Considerations

### Input Validation
- Null pointer checks in Rust
- UTF-8 validation for string conversion
- Empty string handling
- Maximum string length limits (if needed)

### Memory Safety
- Proper use of `CString::into_raw()` and `from_raw()`
- Always call `free_string()` for returned strings
- Use `defer` statements for cleanup in Go

## Extensibility Points

### Adding New Languages
1. Add `Language::Xxx` to languages vector in Rust
2. Rebuild Rust library
3. Regenerate C header
4. Add corresponding tests

### Adding New Functions
1. Add FFI function in Rust with proper C signature
2. Regenerate C header with cbindgen
3. Add Go wrapper method
4. Add tests and documentation

### Configuration Options
- Language selection
- Confidence thresholds
- Detection algorithms
- Performance tuning parameters

## Testing Architecture

### Test Categories
1. **Unit Tests**: Individual function testing
2. **Integration Tests**: End-to-end workflow testing
3. **Performance Tests**: Benchmark and stress testing
4. **Edge Case Tests**: Error conditions and boundary testing

### Test Data
- Single language samples
- Mixed language samples
- Edge cases (empty strings, very short text)
- Performance benchmarks

## Deployment Considerations

### Platform Support
- **macOS**: `liblingua_ffi.dylib`
- **Linux**: `liblingua_ffi.so`
- **Windows**: `lingua_ffi.dll`

### Dependencies
- Rust runtime (included in binary)
- C standard library
- No external language detection dependencies

### Distribution
- Single Go binary with embedded Rust library
- No separate library files needed
- Cross-compilation support via Cargo

This architecture provides a robust foundation for language detection while maintaining good performance and extensibility. 
