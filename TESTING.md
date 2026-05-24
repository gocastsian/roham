# Testing Guide

This document explains how to add, structure, and run tests in the roham project.

## Table of Contents

- [Running Tests](#running-tests)
- [Test Structure](#test-structure)
- [Adding New Tests](#adding-new-tests)
- [Test Types](#test-types)
- [Best Practices](#best-practices)

## Running Tests

### Quick Commands

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests for a specific package
go test ./userapp/service/user/...

# Run tests with verbose output
go test -v ./userapp/...

# Run a specific test
go test -v -run TestRegisterUser ./userapp/service/user/...

# Run tests with race detection
go test -race ./...

# Run tests with count (disable test caching)
go test -count=1 ./...

# Run E2E tests (requires services to be running)
make test-e2e
```

### Test Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out

# View coverage percentage
go test -cover ./...
```

## Test Structure

### Current Test Organization

The project follows a layered architecture, and tests mirror this structure:

```
roham/
├── userapp/
│   ├── delivery/http/
│   │   ├── handler.go
│   │   └── handler_test.go          # HTTP handler tests
│   ├── service/
│   │   ├── user/
│   │   │   ├── service.go
│   │   │   └── service_test.go      # Service layer tests
│   │   └── guard/
│   │       ├── service.go
│   │       └── service_test.go      # Service layer tests
│   ├── repository/
│   │   ├── user.go
│   │   └── user_test.go             # Repository tests (with sqlmock)
│   └── mockRepo/
│       └── userMock/                 # Mock implementations
│           └── user_repo_mock.go
├── filer/
│   └── e2e/                         # End-to-end tests
│       └── upload_test.go
└── pkg/
    ├── password/
    │   ├── hashpassword.go
    │   └── hashpassword_test.go     # Package tests
    └── paginate/sql/
        ├── query.go
        └── query_test.go            # Package tests
```

### Recommended Test Structure

1. **Unit Tests** (`*_test.go` in same package)
   - Test individual functions/methods in isolation
   - Use mocks for dependencies
   - Fast execution, no external dependencies

2. **Integration Tests** (same location, but marked with build tags)
   - Test interactions between components
   - May use test databases or containers
   - Use `//go:build integration` tag

3. **E2E Tests** (`e2e/` directory)
   - Test full workflows
   - Require running services
   - Located in `*/e2e/` directories

## Adding New Tests

### 1. Unit Test Example

```go
package user_test

import (
    "context"
    "testing"
    
    "github.com/gocastsian/roham/userapp/service/user"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestServiceMethod(t *testing.T) {
    // Arrange
    mockRepo := new(MockRepository)
    service := user.NewService(mockRepo, ...)
    
    // Act
    result, err := service.Method(context.Background(), input)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    mockRepo.AssertExpectations(t)
}
```

### 2. HTTP Handler Test Example

```go
package http_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
    e := echo.New()
    handler := NewHandler(...)
    
    body, _ := json.Marshal(requestData)
    req := httptest.NewRequest(http.MethodPost, "/endpoint", bytes.NewBuffer(body))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    ctx := e.NewContext(req, rec)
    
    err := handler.Handle(ctx)
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, rec.Code)
}
```

### 3. Repository Test with sqlmock

```go
package repository

import (
    "context"
    "testing"
    
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/assert"
)

func TestRepositoryMethod(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()
    
    repo := NewUserRepo(Config{}, db, logger)
    
    mock.ExpectPrepare("SELECT ...").
        ExpectQuery().
        WithArgs("arg1", "arg2").
        WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "test"))
    
    result, err := repo.Method(context.Background(), ...)
    
    assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}
```

### 4. Integration Test Example

Create a file with build tag: `//go:build integration`

```go
//go:build integration

package integration_test

import (
    "testing"
    "github.com/stretchr/testify/require"
)

func TestIntegrationFlow(t *testing.T) {
    // Requires actual database connection
    // Use testcontainers or docker-compose
}
```

Run with: `go test -tags=integration ./...`

## Test Types

### 1. Unit Tests
- **Location**: Same package as code (`*_test.go`)
- **Dependencies**: Mocks only
- **Speed**: Fast (< 1s per test)
- **Example**: `userapp/service/user/service_test.go`

### 2. Integration Tests
- **Location**: Same package with `//go:build integration` tag
- **Dependencies**: Test database, Redis, etc.
- **Speed**: Medium (1-10s per test)
- **Example**: Repository tests with real DB

### 3. E2E Tests
- **Location**: `*/e2e/` directories
- **Dependencies**: Full stack running
- **Speed**: Slow (> 10s per test)
- **Example**: `filer/e2e/upload_test.go`

## Best Practices

### 1. Test Naming
- Use descriptive names: `TestFunctionName_Scenario_ExpectedResult`
- Example: `TestRegisterUser_InvalidEmail_ReturnsError`

### 2. Test Organization
- Use table-driven tests for multiple scenarios:
```go
func TestMultipleScenarios(t *testing.T) {
    testCases := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {"success case", "valid", "result", false},
        {"error case", "invalid", "", true},
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

### 3. Test Helpers
- Use `testify/require` for critical assertions (stops on failure)
- Use `testify/assert` for non-critical assertions (continues on failure)
- Create helper functions for common setup:
```go
func setupTestService(t *testing.T) (*Service, *MockRepo) {
    mockRepo := new(MockRepository)
    service := NewService(mockRepo, ...)
    return service, mockRepo
}
```

### 4. Mocking
- Use `testify/mock` for interface mocks
- Keep mocks in `mockRepo/` directories
- Always call `AssertExpectations(t)` to verify mock calls

### 5. Test Data
- Use `t.TempDir()` for temporary files
- Use `t.Cleanup()` for cleanup operations
- Avoid hardcoded test data; use factories:
```go
func TestUser() user.User {
    return user.User{
        ID:       1,
        Username: "testuser",
        Email:    "test@example.com",
    }
}
```

### 6. Error Testing
- Test both success and error paths
- Verify error messages when relevant
- Use `assert.Error` and `assert.ErrorContains`

### 7. Coverage Goals
- Aim for >80% coverage on business logic
- Focus on critical paths first
- Don't obsess over 100% coverage

### 8. Test Isolation
- Each test should be independent
- Don't rely on test execution order
- Clean up after each test

## Test Utilities

### Common Test Helpers

See `testutil/` directory for shared test utilities:
- Database setup/teardown
- Mock factories
- HTTP test helpers
- Test data generators

## CI/CD Integration

Tests run automatically in CI/CD pipeline. See `.github/workflows/` for configuration.

## Troubleshooting

### Tests failing locally but passing in CI
- Check environment variables
- Verify database connections
- Check for race conditions: `go test -race`

### Slow tests
- Use `-short` flag to skip slow tests
- Check for unnecessary database calls
- Use mocks instead of real dependencies when possible

### Flaky tests
- Check for race conditions
- Verify test isolation
- Check for timing issues (use `time.Sleep` or `time.After`)







