#!/bin/bash

# Set the base name of the binary
BINARY_NAME="reft"

# Create a directory for the releases
mkdir -p releases

# Build for Linux (x86_64)
echo "Building for Linux (x86_64)..."
GOOS=linux GOARCH=amd64 go build -o "releases/${BINARY_NAME}-linux-amd64"

# Build for macOS (ARM64)
echo "Building for macOS (ARM64)..."
GOOS=darwin GOARCH=arm64 go build -o "releases/${BINARY_NAME}-darwin-arm64"

echo "Build complete. Binaries are in the 'releases' directory."