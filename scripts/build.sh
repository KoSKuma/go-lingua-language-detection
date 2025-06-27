#!/bin/bash

set -e

echo "Building Rust library..."
cd rust
cargo build --release
cd ..

echo "Generating C header..."
cd rust
cbindgen --config cbindgen.toml --crate lingua-ffi --output target/release/lingua_ffi.h
cd ..

echo "Building Go application..."
go build -o lingua-detector cmd/lingua-detector/main.go

echo "Build completed successfully!"
echo "Run './lingua-detector' to test the application." 
