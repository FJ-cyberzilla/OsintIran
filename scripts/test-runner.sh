#!/bin/bash
# scripts/test-runner.sh

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Starting Comprehensive Test Suite${NC}"

# Function to run tests with timing
run_test() {
    local test_name=$1
    local test_command=$2
    local log_file="test-results/${test_name}.log"
    
    echo -e "${YELLOW}Running ${test_name}...${NC}"
    start_time=$(date +%s)
    
    mkdir -p test-results
    
    if eval $test_command > "$log_file" 2>&1; then
        end_time=$(date +%s)
        duration=$((end_time - start_time))
        echo -e "${GREEN}✅ ${test_name} passed (${duration}s)${NC}"
        return 0
    else
        end_time=$(date +%s)
        duration=$((end_time - start_time))
        echo -e "${RED}❌ ${test_name} failed (${duration}s)${NC}"
        echo -e "${YELLOW}See ${log_file} for details${NC}"
        return 1
    fi
}

# Create test environment
echo -e "${BLUE}Setting up test environment...${NC}"
docker-compose -f tests/docker-compose.test.yml up -d
sleep 30

# Wait for services to be ready
echo -e "${BLUE}Waiting for services to be ready...${NC}"
until curl -s http://localhost:8080/api/v1/health > /dev/null; do
    sleep 5
done

# Run test suites
failed_tests=0

# 1. Unit Tests
run_test "unit_backend" "cd backend && go test ./... -v" || ((failed_tests++))
run_test "unit_frontend" "cd frontend && npm test -- --coverage" || ((failed_tests++))

# 2. Integration Tests
run_test "integration_platforms" "cd tests/integration && npm test" || ((failed_tests++))
run_test "integration_proxies" "cd proxy-pool && go test ./tests/integration/..." || ((failed_tests++))

# 3. E2E Tests
run_test "e2e_tests" "cd tests/e2e && npx cypress run --headless" || ((failed_tests++))

# 4. Load Tests (if not in CI)
if [ -z "$CI" ]; then
    run_test "load_smoke" "cd tests/load && k6 run k6/smoke_test.js" || ((failed_tests++))
    run_test "load_stress" "cd tests/load && locust -f locustfile.py --headless -u 10 -r 1 -t 1m" || ((failed_tests++))
fi

# Generate test report
echo -e "${BLUE}Generating test report...${NC}"
generate_test_report() {
    cat > test-results/report.html << EOF
<!DOCTYPE html>
<html>
<head>
    <title>Test Execution Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .passed { color: green; }
        .failed { color: red; }
        .test-result { margin: 10px 0; padding: 10px; border-left: 4px solid; }
    </style>
</head>
<body>
    <h1>Test Execution Report</h1>
    <p>Generated on: $(date)</p>
    <div class="summary">
        <h2>Summary</h2>
        <p>Total Tests: $((failed_tests > 0 ? "❌" : "✅"))</p>
        <p>Failed Tests: $failed_tests</p>
    </div>
</body>
</html>
EOF
}

generate_test_report

# Cleanup
echo -e "${BLUE}Cleaning up test environment...${NC}"
docker-compose -f tests/docker-compose.test.yml down

# Final result
if [ $failed_tests -eq 0 ]; then
    echo -e "${GREEN}🎉 All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}❌ $failed_tests test suite(s) failed${NC}"
    exit 1
fi
