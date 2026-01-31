#!/bin/bash

set -e

echo "Building Go backend..."

# Set build variables
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GOOS=${GOOS:-linux}
GOARCH=${GOARCH:-amd64}

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
