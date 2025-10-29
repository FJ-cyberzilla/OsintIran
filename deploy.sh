#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}Starting deployment...${NC}"

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}Docker is not running. Please start Docker first.${NC}"
    exit 1
fi

# Build the image
echo -e "${YELLOW}Building Docker image...${NC}"
docker build -t my-app:latest .

# Stop existing container if running
if docker ps -a | grep -q my-app; then
    echo -e "${YELLOW}Stopping existing container...${NC}"
    docker stop my-app || true
    docker rm my-app || true
fi

# Run new container
echo -e "${YELLOW}Starting new container...${NC}"
docker run -d \
    --name my-app \
    --restart unless-stopped \
    -p 8080:8080 \
    -e ENVIRONMENT=production \
    -e LOG_LEVEL=info \
    my-app:latest

echo -e "${GREEN}Deployment completed successfully!${NC}"
echo -e "${YELLOW}Application is running on http://localhost:8080${NC}"

# Show logs
echo -e "${YELLOW}Showing container logs (Ctrl+C to exit)...${NC}"
docker logs -f my-app
