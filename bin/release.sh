#!/bin/bash

# Script to build and release redup

set -e

VERSION=${1:-$(git describe --tags --always --dirty)}
BUILD_DIR="build"

echo "Building redup version: $VERSION"

# Create build directory
mkdir -p $BUILD_DIR

# Build for multiple platforms
echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$VERSION" -o $BUILD_DIR/redup-linux-amd64 .

echo "Building for macOS..."
GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$VERSION" -o $BUILD_DIR/redup-darwin-amd64 .

echo "Building for Windows..."
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$VERSION" -o $BUILD_DIR/redup-windows-amd64.exe .

echo "Build complete! Binaries are in $BUILD_DIR/"
ls -la $BUILD_DIR/
