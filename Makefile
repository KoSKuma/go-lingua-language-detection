.PHONY: build clean test run

# Build the Rust library and then the Go application
build: build-rust build-go

# Build the Rust library
build-rust:
	cd rust && cargo build --release
	cp rust/target/release/liblingua_ffi.* rust/target/release/

# Build the Go application
build-go: build-rust
	go build -o lingua-detector cmd/lingua-detector/main.go

# Clean build artifacts
clean:
	cd rust && cargo clean
	rm -f lingua-detector
	rm -f rust/target/release/liblingua_ffi.*

# Run tests
test:
	go test -v ./pkg/lingua

# Install dependencies
deps:
	go mod tidy
	cd rust && cargo build

# Development build (debug mode)
dev: build-rust
	cd rust && cargo build
	go build -o lingua-detector cmd/lingua-detector/main.go 

# Run the built artifact
run:
	./lingua-detector
