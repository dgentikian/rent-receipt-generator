#!/bin/bash

set -e

echo "Building React frontend..."

# Install dependencies if node_modules doesn't exist
if [ ! -d "node_modules" ]; then
    echo "Installing dependencies..."
    npm ci
fi

# Build for production
echo "Running production build..."
npm run build

echo "Frontend built successfully: dist/"
