.PHONY: build clean test

# Build the Rust library and then the Go application
build: build-rust build-go

# Build the Rust library
build-rust:
	cargo build --release
	cp target/release/liblingua_ffi.* target/release/

# Build the Go application
build-go: build-rust
	go build -o lingua-detector lingua.go

# Clean build artifacts
clean:
	cargo clean
	rm -f lingua-detector
	rm -f target/release/liblingua_ffi.*

# Run tests
test: build
	./lingua-detector

# Install dependencies
deps:
	go mod tidy
	cargo build

# Development build (debug mode)
dev: build-rust
	cargo build
	go build -o lingua-detector lingua.go 
