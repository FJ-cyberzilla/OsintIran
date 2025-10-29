#!/bin/bash
set -e

echo "🔧 Fixing Redis Go module version..."

# Check if go.mod exists
if [ ! -f go.mod ]; then
    echo "❌ go.mod not found. Please run 'go mod init your-module-name' first"
    exit 1
fi

# Remove incorrect redis require
echo "🗑️  Removing incorrect redis version..."
go mod edit -droprequire=github.com/redis/go-redis 2>/dev/null || true

# Add correct redis version
echo "📦 Adding correct redis version..."
go mod edit -require=github.com/redis/go-redis/v9@v9.0.5

# Update dependencies
echo "🔄 Updating dependencies..."
go mod tidy

# Verify
echo "✅ Verifying fix..."
if grep -q "github.com/redis/go-redis/v9 v9" go.mod; then
    echo "✅ Redis version fixed successfully!"
    echo "📋 Updated go.mod content:"
    grep redis go.mod
else
    echo "❌ Failed to fix redis version"
    exit 1
fi

# Test build
echo "🧪 Testing build..."
go mod download
echo "✅ All dependencies downloaded successfully!"

echo ""
echo "🚀 Now you can build your Docker image:"
echo "   docker build -t your-app ."
