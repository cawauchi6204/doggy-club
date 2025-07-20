# DoggyClub Backend - Testing Guide

This document provides comprehensive information about the testing strategy, test suites, and how to run tests for the DoggyClub backend.

## Table of Contents

- [Testing Strategy](#testing-strategy)
- [Test Structure](#test-structure)
- [Running Tests](#running-tests)
- [Test Types](#test-types)
- [Continuous Integration](#continuous-integration)
- [Test Data Management](#test-data-management)
- [Performance Testing](#performance-testing)
- [Troubleshooting](#troubleshooting)

## Testing Strategy

Our testing strategy follows the testing pyramid approach:

```
    /\
   /E2E\     ← Few, slow, high confidence
  /----\
 /Integration\ ← Some, medium speed, medium confidence  
/----------\
|   Unit   | ← Many, fast, low confidence
\----------/
```

### Test Coverage Goals

- **Unit Tests**: >90% coverage
- **Integration Tests**: >80% coverage of critical paths
- **E2E Tests**: All user workflows covered
- **Performance Tests**: Key endpoints benchmarked

## Test Structure

```
tests/
├── e2e/                    # End-to-end tests
│   └── main_test.go        # Complete user workflows
├── integration/            # Integration tests
│   └── api_integration_test.go  # API endpoint tests
├── performance/            # Performance and load tests
│   └── load_test.go        # Load testing suite
└── testdata/              # Test fixtures and data

pkg/
├── services/
│   ├── auth_service_test.go      # Unit tests for auth service
│   ├── user_service_test.go      # Unit tests for user service
│   └── ...
├── handlers/
│   ├── auth_handler_test.go      # Integration tests for handlers
│   └── ...
└── testutils/
    └── testutils.go              # Test utilities and helpers
```

## Running Tests

### Quick Start

```bash
# Run all unit and integration tests
make test

# Run all tests including E2E and load tests
make test-all

# Run comprehensive test suite with proper setup
./test-runner.sh
```

### Individual Test Types

```bash
# Unit tests only
make test-unit

# Integration tests only
make test-integration

# E2E tests (requires running application)
make test-e2e

# Load/performance tests
make test-load

# Benchmark tests
make test-benchmark
```

### Advanced Test Runner

The `test-runner.sh` script provides comprehensive testing with automatic service management:

```bash
# Full test suite with all services
./test-runner.sh

# Skip specific test types
./test-runner.sh --skip-e2e --skip-integration

# Include load tests (disabled by default)
./test-runner.sh --include-load

# Don't cleanup services after tests
./test-runner.sh --no-cleanup

# Show help
./test-runner.sh --help
```

### Test Environment Setup

```bash
# Start test databases
make test-env-up

# Run tests with environment
make test-full

# Stop test databases
make test-env-down
```

### Coverage Reports

```bash
# Generate coverage report
make test-coverage

# View coverage in browser
open coverage.html

# Text coverage report
go tool cover -func=coverage.out
```

## Test Types

### 1. Unit Tests

Unit tests test individual functions and methods in isolation.

**Location**: `pkg/*/test.go`
**Purpose**: Test business logic, data transformations, validations

```bash
# Run unit tests
make test-unit

# Run with coverage
go test -v -race -cover ./pkg/services/...
```

**Example**:
```go
func TestAuthService_Register(t *testing.T) {
    // Setup
    db := testutils.SetupTestDB(t)
    authService := NewAuthService(db, "test-secret")
    
    // Test
    req := RegisterRequest{
        Email:     "test@example.com",
        Password:  "SecurePassword123!",
        FirstName: "Test",
        LastName:  "User",
    }
    
    resp, err := authService.Register(req)
    
    // Assertions
    assert.NoError(t, err)
    assert.NotEmpty(t, resp.AccessToken)
    assert.Equal(t, "test@example.com", resp.User.Email)
}
```

### 2. Integration Tests

Integration tests verify that different components work together correctly.

**Location**: `tests/integration/`
**Purpose**: Test API endpoints, database interactions, service integrations

```bash
# Run integration tests
make test-integration

# With test database
TEST_DB_HOST=localhost TEST_DB_PORT=5433 make test-integration
```

**Example**:
```go
func (suite *APIIntegrationTestSuite) TestUserRegistrationFlow() {
    // Test complete registration flow
    registerData := map[string]interface{}{
        "email":    "test@example.com",
        "password": "SecurePassword123!",
        // ...
    }
    
    // Make HTTP request
    resp, err := http.Post(suite.server.URL+"/auth/register", 
        "application/json", bytes.NewBuffer(body))
    
    // Verify response
    suite.Equal(http.StatusCreated, resp.StatusCode)
    // Verify database state
    // Verify side effects
}
```

### 3. End-to-End Tests

E2E tests simulate real user interactions with the complete system.

**Location**: `tests/e2e/`
**Purpose**: Test complete user workflows, API interactions

```bash
# Start application first
make run &

# Run E2E tests
make test-e2e
```

**Features Tested**:
- User registration and authentication
- Profile management
- Dog profile CRUD operations
- Post creation and interactions
- Social features (likes, comments)
- Encounter detection
- Gift system
- Notifications

### 4. Performance Tests

Performance tests measure system performance under various load conditions.

**Location**: `tests/performance/`
**Purpose**: Load testing, benchmark testing, performance regression detection

```bash
# Run performance tests
make test-load

# Run benchmarks
make test-benchmark

# CPU profiling
make profile-cpu

# Memory profiling
make profile-mem
```

**Test Scenarios**:
- Health endpoint load (1000 concurrent requests)
- User profile reads (high frequency)
- Post creation (write-heavy workload)
- Mixed workload simulation
- Memory leak detection

## Continuous Integration

### GitHub Actions

Our CI pipeline (`.github/workflows/ci.yml`) runs:

1. **Code Quality**:
   - Format checking
   - Linting (golangci-lint)
   - Security scanning (gosec)
   - Go vet

2. **Testing**:
   - Unit tests with coverage
   - Integration tests with test databases
   - Performance benchmarks (on main branch)

3. **Build**:
   - Cross-platform builds
   - Docker image builds

4. **Deployment**:
   - Automated deployment on main branch

### Local CI Simulation

```bash
# Run CI pipeline locally
make ci

# Quick CI checks
make ci-quick
```

## Test Data Management

### Test Database

- Separate test database (`doggyclub_test`)
- Automatic migration before tests
- Cleanup after each test suite
- Isolated test transactions

### Test Fixtures

Located in `tests/testdata/`:
- User profiles
- Dog profiles  
- Posts and comments
- Sample images
- Configuration files

### Test Utilities

The `pkg/testutils/` package provides:
- Database setup/teardown
- Test data factories
- HTTP test helpers
- Assertion helpers
- Mock services

```go
// Database utilities
db := testutils.SetupTestDB(t)
defer testutils.CleanupTestDB(t, db)

// Test data factories
user := testutils.CreateTestUser(t, db)
dog := testutils.CreateTestDog(t, db, user.ID)

// HTTP utilities
router := testutils.SetupTestRouter()
req := testutils.CreateAuthRequest("GET", "/users/me", nil, token)
```

## Performance Testing

### Load Testing

Our load tests simulate realistic user behavior:

```go
func TestMixedWorkload(t *testing.T) {
    // 60% read operations (profiles, posts)
    // 30% interaction operations (likes, comments)  
    // 10% write operations (new posts, profile updates)
}
```

### Benchmarking

```bash
# Benchmark specific functions
go test -bench=BenchmarkAuthService -benchmem ./pkg/services/

# Compare benchmarks
go test -bench=. -count=5 ./pkg/services/ > old.txt
# Make changes
go test -bench=. -count=5 ./pkg/services/ > new.txt
benchcmp old.txt new.txt
```

### Performance Metrics

We track:
- **Response Time**: p50, p95, p99 percentiles
- **Throughput**: Requests per second
- **Error Rate**: Percentage of failed requests
- **Resource Usage**: CPU, memory, database connections
- **Database Performance**: Query time, connection pool usage

### Performance Goals

- **API Response Time**: p95 < 200ms
- **Database Queries**: < 100ms average
- **Throughput**: > 1000 req/sec for reads, > 100 req/sec for writes
- **Error Rate**: < 0.1% under normal load
- **Memory Usage**: Stable over time (no leaks)

## Test Configuration

### Environment Variables

```bash
# Test database
TEST_DB_HOST=localhost
TEST_DB_PORT=5433
TEST_DB_NAME=doggyclub_test
TEST_DB_USER=postgres
TEST_DB_PASSWORD=password

# Test Redis
TEST_REDIS_HOST=localhost
TEST_REDIS_PORT=6380

# Test flags
SKIP_E2E_TESTS=true     # Skip E2E tests in CI
SKIP_LOAD_TESTS=true    # Skip load tests by default
TEST_TIMEOUT=300s       # Test timeout
```

### Docker Test Environment

```yaml
# docker-compose.test.yml
version: '3.8'
services:
  test-db:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: doggyclub_test
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    ports:
      - "5433:5432"
  
  test-redis:
    image: redis:7-alpine
    ports:
      - "6380:6379"
```

## Test Best Practices

### 1. Test Structure

```go
func TestFunction(t *testing.T) {
    // Arrange - Setup test data and environment
    
    // Act - Execute the function being tested
    
    // Assert - Verify the results
}
```

### 2. Table-Driven Tests

```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    bool
        wantErr bool
    }{
        {"valid email", "test@example.com", true, false},
        {"invalid email", "invalid", false, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ValidateEmail(tt.input)
            assert.Equal(t, tt.want, got)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### 3. Test Isolation

- Each test should be independent
- Use transactions for database tests
- Clean up resources after tests
- Don't rely on test execution order

### 4. Meaningful Test Names

```go
func TestAuthService_Register_WithValidData_ReturnsTokenAndUser(t *testing.T)
func TestAuthService_Register_WithExistingEmail_ReturnsError(t *testing.T)
func TestAuthService_Register_WithWeakPassword_ReturnsValidationError(t *testing.T)
```

### 5. Test Coverage

- Aim for high test coverage but focus on critical paths
- Test both happy path and error scenarios
- Test edge cases and boundary conditions
- Test concurrent access where applicable

## Troubleshooting

### Common Issues

1. **Database Connection Errors**
   ```bash
   # Check if test database is running
   docker ps | grep postgres
   
   # Check connection
   psql -h localhost -p 5433 -U postgres -d doggyclub_test
   ```

2. **Port Conflicts**
   ```bash
   # Check what's using the port
   lsof -i :5433
   
   # Kill process using port
   kill -9 <PID>
   ```

3. **Test Timeouts**
   ```bash
   # Increase timeout for slow tests
   go test -timeout=600s ./...
   ```

4. **Memory Issues**
   ```bash
   # Run with race detector
   go test -race ./...
   
   # Memory profiling
   go test -memprofile=mem.prof ./...
   ```

### Debugging Tests

```bash
# Verbose output
go test -v ./...

# Run specific test
go test -v -run=TestSpecificFunction ./...

# Debug with Delve
dlv test ./pkg/services -- -test.run=TestFunction
```

### Performance Debugging

```bash
# CPU profiling
go test -cpuprofile=cpu.prof -bench=. ./...
go tool pprof cpu.prof

# Memory profiling  
go test -memprofile=mem.prof -bench=. ./...
go tool pprof mem.prof

# Trace analysis
go test -trace=trace.out ./...
go tool trace trace.out
```

### Test Data Issues

```bash
# Reset test database
make test-env-down
make test-env-up

# Check test data
psql -h localhost -p 5433 -U postgres -d doggyclub_test -c "SELECT * FROM users;"

# Manual cleanup
psql -h localhost -p 5433 -U postgres -d doggyclub_test -c "TRUNCATE users CASCADE;"
```

### CI/CD Issues

1. **Test failures in CI but not locally**
   - Check environment differences
   - Verify all dependencies are installed
   - Check timing issues (add delays if needed)

2. **Flaky tests**
   - Identify race conditions
   - Add proper waits for async operations
   - Use deterministic test data

3. **Resource limits**
   - Reduce test parallelism
   - Optimize test data size
   - Use test caching

## Contributing to Tests

### Adding New Tests

1. **Unit Tests**: Add to the same package as the code being tested
2. **Integration Tests**: Add to `tests/integration/`
3. **E2E Tests**: Add to `tests/e2e/`
4. **Performance Tests**: Add to `tests/performance/`

### Test Naming Conventions

- File: `*_test.go`
- Function: `TestPackage_Function_Scenario_ExpectedResult`
- Benchmark: `BenchmarkPackage_Function`

### Required for New Features

- Unit tests for all new functions
- Integration tests for new API endpoints
- E2E tests for new user workflows
- Performance tests for critical paths
- Documentation updates

---

For more information about specific testing aspects, refer to the individual test files and comments within the codebase.