#!/bin/bash

set -e

echo "Building Go backend..."

# Set build variables
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')

# Auto-detect platform if not set
if [ -z "$GOOS" ]; then
    GOOS=$(go env GOOS)
fi
if [ -z "$GOARCH" ]; then
    GOARCH=$(go env GOARCH)
fi

# Create bin directory if it doesn't exist
mkdir -p bin

# Build binary
echo "Building for ${GOOS}/${GOARCH}..."
CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build \
    -ldflags="-w -s -X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME" \
    -o bin/rent-receipt-server \
    ./cmd/server

echo "Backend built successfully: bin/rent-receipt-server"
echo "Version: $VERSION"
echo "Build time: $BUILD_TIME"
