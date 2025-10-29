#!/bin/bash
set -e

echo "🔍 Debugging Docker build..."

# Check if go.mod exists
if [ ! -f go.mod ]; then
    echo "❌ go.mod not found. Initializing..."
    go mod init myapp
fi

# Generate go.sum if missing
if [ ! -f go.sum ]; then
    echo "📝 Generating go.sum..."
    go mod tidy
fi

# Test Go module download locally
echo "🧪 Testing Go module download locally..."
go mod download

# Build with detailed output
echo "🐳 Building Docker image..."
docker build --progress=plain --no-cache -t myapp-debug .

echo "✅ Build completed successfully!"
