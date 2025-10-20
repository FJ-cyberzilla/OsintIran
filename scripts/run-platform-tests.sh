#!/bin/bash
# scripts/run-platform-tests.sh

#!/bin/bash
set -e

echo "🚀 Starting Platform Integration Tests"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Test configuration
export NODE_ENV=test
export TEST_TIMEOUT=120000

# Check dependencies
echo -e "${BLUE}Checking dependencies...${NC}"
if ! command -v node &> /dev/null; then
    echo -e "${RED}Node.js is required but not installed.${NC}"
    exit 1
fi

if ! command -v npm &> /dev/null; then
    echo -e "${RED}npm is required but not installed.${NC}"
    exit 1
fi

# Install test dependencies if needed
echo -e "${BLUE}Installing test dependencies...${NC}"
cd tests/integration
npm install

# Create test environment
echo -e "${BLUE}Setting up test environment...${NC}"
docker-compose -f docker-compose.test.yml up -d
sleep 30

# Wait for services
echo -e "${BLUE}Waiting for test services...${NC}"
until curl -s http://localhost:8080/api/v1/health > /dev/null; do
    sleep 5
done

# Run platform tests
echo -e "${BLUE}Running platform integration tests...${NC}"
node platforms/platform-test-runner.js

# Capture exit code
TEST_EXIT_CODE=$?

# Cleanup
echo -e "${BLUE}Cleaning up test environment...${NC}"
docker-compose -f docker-compose.test.yml down

# Report results
if [ $TEST_EXIT_CODE -eq 0 ]; then
    echo -e "${GREEN}🎉 All platform integration tests passed!${NC}"
else
    echo -e "${RED}❌ Some platform integration tests failed${NC}"
fi

exit $TEST_EXIT_CODE
