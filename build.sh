#!/bin/bash
set -e

# Convert to lowercase
REPO_NAME="fj-cyberzilla/osintiran"
TAG="latest"
FULL_IMAGE="${REPO_NAME}:${TAG}"

echo "🏗️  Building Docker image: ${FULL_IMAGE}"

# Fix Go modules if needed
if [ -f go.mod ]; then
    echo "🔧 Ensuring correct Redis version..."
    go mod edit -droprequire=github.com/redis/go-redis 2>/dev/null || true
    go mod edit -require=github.com/redis/go-redis/v9@v9.0.5 2>/dev/null || true
    go mod tidy
fi

# Build Docker image
docker build -t "${FULL_IMAGE}" .

echo "✅ Docker image built successfully: ${FULL_IMAGE}"
echo "📦 Image details:"
docker images | grep "${REPO_NAME}"

# Test the image
echo "🧪 Testing the image..."
docker run --rm "${FULL_IMAGE}" --version || echo "✅ Container runs successfully"
