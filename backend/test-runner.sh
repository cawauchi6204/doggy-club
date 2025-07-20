#!/bin/bash

# DoggyClub Backend - Comprehensive Test Runner
# This script runs all tests in the correct order with proper setup and teardown

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
TEST_DB_NAME="doggyclub_test"
TEST_DB_USER="postgres"
TEST_DB_PASSWORD="password"
TEST_DB_HOST="localhost"
TEST_DB_PORT="5433"
TEST_REDIS_HOST="localhost"
TEST_REDIS_PORT="6380"
APP_PORT="8080"

# Test flags
RUN_UNIT=${RUN_UNIT:-true}
RUN_INTEGRATION=${RUN_INTEGRATION:-true}
RUN_E2E=${RUN_E2E:-true}
RUN_LOAD=${RUN_LOAD:-false}  # Disabled by default
CLEANUP=${CLEANUP:-true}

# Function to print colored output
print_status() {
    echo -e "${BLUE}[TEST-RUNNER]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if service is running
wait_for_service() {
    local host=$1
    local port=$2
    local service=$3
    local max_attempts=30
    local attempt=1

    print_status "Waiting for $service to be ready on $host:$port..."
    
    while [ $attempt -le $max_attempts ]; do
        if nc -z $host $port 2>/dev/null; then
            print_success "$service is ready!"
            return 0
        fi
        
        if [ $attempt -eq $max_attempts ]; then
            print_error "$service failed to start within timeout"
            return 1
        fi
        
        echo -n "."
        sleep 1
        ((attempt++))
    done
}

# Function to start test database
start_test_db() {
    print_status "Starting test database..."
    
    # Stop existing container if running
    docker stop doggyclub-test-db 2>/dev/null || true
    docker rm doggyclub-test-db 2>/dev/null || true
    
    # Start new container
    docker run -d --name doggyclub-test-db \
        -e POSTGRES_DB=$TEST_DB_NAME \
        -e POSTGRES_USER=$TEST_DB_USER \
        -e POSTGRES_PASSWORD=$TEST_DB_PASSWORD \
        -p $TEST_DB_PORT:5432 \
        postgres:15-alpine
    
    wait_for_service $TEST_DB_HOST $TEST_DB_PORT "PostgreSQL"
}

# Function to start test Redis
start_test_redis() {
    print_status "Starting test Redis..."
    
    # Stop existing container if running
    docker stop doggyclub-test-redis 2>/dev/null || true
    docker rm doggyclub-test-redis 2>/dev/null || true
    
    # Start new container
    docker run -d --name doggyclub-test-redis \
        -p $TEST_REDIS_PORT:6379 \
        redis:7-alpine
    
    wait_for_service $TEST_REDIS_HOST $TEST_REDIS_PORT "Redis"
}

# Function to start application for E2E tests
start_test_app() {
    print_status "Starting application for E2E tests..."
    
    # Kill existing process if running
    pkill -f "doggyclub-backend" || true
    sleep 2
    
    # Set test environment variables
    export ENV=test
    export PORT=$APP_PORT
    export DB_HOST=$TEST_DB_HOST
    export DB_PORT=$TEST_DB_PORT
    export DB_NAME=$TEST_DB_NAME
    export DB_USER=$TEST_DB_USER
    export DB_PASSWORD=$TEST_DB_PASSWORD
    export DB_SSL_MODE=disable
    export REDIS_HOST=$TEST_REDIS_HOST
    export REDIS_PORT=$TEST_REDIS_PORT
    export JWT_SECRET=test-jwt-secret-for-e2e-tests
    export ENABLE_ENCOUNTER_DETECTION=true
    export ENABLE_PUSH_NOTIFICATIONS=false
    export ENABLE_PREMIUM_FEATURES=true
    
    # Start application in background
    go run ./cmd/api &
    APP_PID=$!
    
    # Wait for application to be ready
    wait_for_service localhost $APP_PORT "Application"
    
    return $APP_PID
}

# Function to stop test services
cleanup_services() {
    if [ "$CLEANUP" = "true" ]; then
        print_status "Cleaning up test services..."
        
        # Stop application
        if [ ! -z "$APP_PID" ]; then
            kill $APP_PID 2>/dev/null || true
        fi
        pkill -f "doggyclub-backend" || true
        
        # Stop containers
        docker stop doggyclub-test-db doggyclub-test-redis 2>/dev/null || true
        docker rm doggyclub-test-db doggyclub-test-redis 2>/dev/null || true
        
        print_success "Cleanup completed"
    fi
}

# Function to run unit tests
run_unit_tests() {
    print_status "Running unit tests..."
    
    go test -v -race -timeout=300s \
        -coverprofile=coverage-unit.out \
        ./pkg/services/... \
        ./pkg/utils/... \
        ./pkg/middleware/...
    
    print_success "Unit tests completed"
}

# Function to run integration tests
run_integration_tests() {
    print_status "Running integration tests..."
    
    # Set environment variables for integration tests
    export TEST_DB_HOST=$TEST_DB_HOST
    export TEST_DB_PORT=$TEST_DB_PORT
    export TEST_DB_NAME=$TEST_DB_NAME
    export TEST_DB_USER=$TEST_DB_USER
    export TEST_DB_PASSWORD=$TEST_DB_PASSWORD
    export TEST_REDIS_HOST=$TEST_REDIS_HOST
    export TEST_REDIS_PORT=$TEST_REDIS_PORT
    
    go test -v -race -timeout=300s \
        -coverprofile=coverage-integration.out \
        ./tests/integration/...
    
    print_success "Integration tests completed"
}

# Function to run E2E tests
run_e2e_tests() {
    print_status "Running E2E tests..."
    
    go test -v -timeout=600s ./tests/e2e/...
    
    print_success "E2E tests completed"
}

# Function to run load tests
run_load_tests() {
    print_status "Running load tests..."
    
    # Set environment variables for load tests
    export TEST_DB_HOST=$TEST_DB_HOST
    export TEST_DB_PORT=$TEST_DB_PORT
    export TEST_DB_NAME=$TEST_DB_NAME
    export TEST_DB_USER=$TEST_DB_USER
    export TEST_DB_PASSWORD=$TEST_DB_PASSWORD
    export TEST_REDIS_HOST=$TEST_REDIS_HOST
    export TEST_REDIS_PORT=$TEST_REDIS_PORT
    
    go test -v -timeout=1200s ./tests/performance/...
    
    print_success "Load tests completed"
}

# Function to merge coverage reports
merge_coverage() {
    print_status "Merging coverage reports..."
    
    # Install gocovmerge if not present
    if ! command -v gocovmerge &> /dev/null; then
        go install github.com/wadey/gocovmerge@latest
    fi
    
    # Merge coverage files
    if [ -f coverage-unit.out ] && [ -f coverage-integration.out ]; then
        gocovmerge coverage-unit.out coverage-integration.out > coverage-merged.out
        go tool cover -html=coverage-merged.out -o coverage-merged.html
        go tool cover -func=coverage-merged.out
        
        print_success "Coverage reports merged: coverage-merged.html"
    elif [ -f coverage-unit.out ]; then
        cp coverage-unit.out coverage-merged.out
        go tool cover -html=coverage-merged.out -o coverage-merged.html
        print_success "Unit test coverage report: coverage-merged.html"
    fi
}

# Function to run code quality checks
run_quality_checks() {
    print_status "Running code quality checks..."
    
    # Format check
    print_status "Checking code formatting..."
    if ! go fmt ./...; then
        print_error "Code formatting issues found"
        return 1
    fi
    
    # Vet
    print_status "Running go vet..."
    if ! go vet ./...; then
        print_error "Go vet found issues"
        return 1
    fi
    
    # Lint (if available)
    if command -v golangci-lint &> /dev/null; then
        print_status "Running linter..."
        if ! golangci-lint run ./...; then
            print_warning "Linter found issues (continuing anyway)"
        fi
    else
        print_warning "golangci-lint not installed, skipping lint check"
    fi
    
    # Security scan (if available)
    if command -v gosec &> /dev/null; then
        print_status "Running security scan..."
        if ! gosec ./...; then
            print_warning "Security scan found issues (continuing anyway)"
        fi
    else
        print_warning "gosec not installed, skipping security scan"
    fi
    
    print_success "Code quality checks completed"
}

# Trap to ensure cleanup on exit
trap cleanup_services EXIT

# Main execution
main() {
    print_status "Starting DoggyClub Backend Test Suite"
    echo "====================================================="
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --skip-unit)
                RUN_UNIT=false
                shift
                ;;
            --skip-integration)
                RUN_INTEGRATION=false
                shift
                ;;
            --skip-e2e)
                RUN_E2E=false
                shift
                ;;
            --include-load)
                RUN_LOAD=true
                shift
                ;;
            --no-cleanup)
                CLEANUP=false
                shift
                ;;
            --help)
                echo "Usage: $0 [options]"
                echo "Options:"
                echo "  --skip-unit         Skip unit tests"
                echo "  --skip-integration  Skip integration tests"
                echo "  --skip-e2e          Skip E2E tests"
                echo "  --include-load      Include load tests"
                echo "  --no-cleanup        Don't cleanup services after tests"
                echo "  --help              Show this help"
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                exit 1
                ;;
        esac
    done
    
    # Check dependencies
    print_status "Checking dependencies..."
    if ! command -v docker &> /dev/null; then
        print_error "Docker is required but not installed"
        exit 1
    fi
    
    if ! command -v go &> /dev/null; then
        print_error "Go is required but not installed"
        exit 1
    fi
    
    if ! command -v nc &> /dev/null; then
        print_error "netcat (nc) is required but not installed"
        exit 1
    fi
    
    # Run code quality checks first
    run_quality_checks
    
    # Start test services
    start_test_db
    start_test_redis
    
    # Give services time to fully initialize
    sleep 3
    
    # Run tests in order
    TEST_FAILURES=0
    
    if [ "$RUN_UNIT" = "true" ]; then
        if ! run_unit_tests; then
            ((TEST_FAILURES++))
        fi
    fi
    
    if [ "$RUN_INTEGRATION" = "true" ]; then
        if ! run_integration_tests; then
            ((TEST_FAILURES++))
        fi
    fi
    
    if [ "$RUN_E2E" = "true" ]; then
        start_test_app
        APP_PID=$!
        sleep 5  # Give app time to initialize
        
        if ! run_e2e_tests; then
            ((TEST_FAILURES++))
        fi
    fi
    
    if [ "$RUN_LOAD" = "true" ]; then
        if ! run_load_tests; then
            ((TEST_FAILURES++))
        fi
    fi
    
    # Merge coverage reports
    merge_coverage
    
    # Summary
    echo ""
    echo "====================================================="
    if [ $TEST_FAILURES -eq 0 ]; then
        print_success "All tests passed! 🎉"
        echo ""
        echo "Test Summary:"
        [ "$RUN_UNIT" = "true" ] && echo "✅ Unit tests"
        [ "$RUN_INTEGRATION" = "true" ] && echo "✅ Integration tests"
        [ "$RUN_E2E" = "true" ] && echo "✅ E2E tests"
        [ "$RUN_LOAD" = "true" ] && echo "✅ Load tests"
        echo ""
        exit 0
    else
        print_error "Some tests failed! ($TEST_FAILURES test suite(s) failed)"
        echo ""
        echo "Failed test suites: $TEST_FAILURES"
        echo ""
        exit 1
    fi
}

# Run main function with all arguments
main "$@"