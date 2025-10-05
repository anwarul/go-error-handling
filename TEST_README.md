# Go Error Handling - Test Suite

This project includes comprehensive tests for all error handling patterns and packages.

## Test Coverage

- **basic**: 100% coverage - Tests division function with various inputs and error cases
- **custom**: 100% coverage - Tests custom ValidationError type
- **database**: 100% coverage - Tests DatabaseError type with unwrapping
- **example**: 91.7% coverage - Tests all example functions for panics and logic
- **formatted**: 100% coverage - Tests formatted error messages with edge cases
- **user**: 100% coverage - Tests user validation and database operations
- **utils**: No statements - Only contains sentinel error constants
- **wrapping**: 84.6% coverage - Tests error wrapping chains with real file operations

## Running Tests

### Run all tests:
```bash
go test ./...
```

### Run tests with coverage:
```bash
go test -cover ./...
```

### Run tests with verbose output:
```bash
go test -v ./...
```

### Run specific package tests:
```bash
go test -v ./basic
go test -v ./custom
go test -v ./database
```

### Generate coverage report:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Test Patterns Covered

### 1. Basic Error Handling (`basic/`)
- ✅ Successful operations with various inputs
- ✅ Division by zero error cases
- ✅ Edge cases (divide by 1, divide zero by number)
- ✅ Error message validation

### 2. Custom Error Types (`custom/`)
- ✅ Error message formatting with different value types
- ✅ Error interface implementation
- ✅ Field validation and access

### 3. Formatted Errors (`formatted/`)
- ✅ Valid age validation (boundary values)
- ✅ Negative age error cases
- ✅ Too old age error cases
- ✅ Error message content validation

### 4. Database Errors (`database/`)
- ✅ DatabaseError creation and field validation
- ✅ Error message formatting
- ✅ Error unwrapping with `errors.Is()` and `errors.As()`
- ✅ Retryable flag testing
- ✅ Custom unwrap function testing

### 5. User Package (`user/`)
- ✅ User validation with ValidationError
- ✅ Email validation (empty and not found cases)
- ✅ QueryUsers database error simulation
- ✅ Error chain verification with `errors.Is()` and `errors.As()`

### 6. Utils Package (`utils/`)
- ✅ Sentinel error identity verification
- ✅ Error message consistency
- ✅ Uniqueness of different sentinel errors
- ✅ Wrapped error chain detection

### 7. Wrapping Package (`wrapping/`)
- ✅ Error chain creation and traversal
- ✅ File operation error wrapping
- ✅ Multiple levels of error context
- ✅ `os.ErrNotExist` detection in wrapped chains
- ✅ Temporary file testing for success cases

### 8. Example Functions (`example/`)
- ✅ Panic prevention tests for all example functions
- ✅ Custom error validation with specific codes
- ✅ Integration testing of error detection logic
- ✅ End-to-end workflow testing

## Key Testing Principles Applied

1. **Error Identity**: Tests use `errors.Is()` to verify error identity in chains
2. **Type Assertions**: Tests use `errors.As()` to extract specific error types
3. **Boundary Testing**: Tests validate edge cases and boundary conditions
4. **Error Message Validation**: Tests verify error messages contain expected information
5. **Panic Prevention**: Integration tests ensure functions don't panic
6. **Chain Traversal**: Tests verify error unwrapping works correctly

## Test Naming Conventions

- `TestFunctionName_Scenario` - Main test function pattern
- `TestFunctionName_Success` - Happy path scenarios
- `TestFunctionName_ErrorCase` - Specific error conditions
- `TestFunctionName_DoesNotPanic` - Panic prevention tests

## Mock and Stub Strategies

- **File Operations**: Uses temporary files for success cases, non-existent files for error cases
- **Database Operations**: Simulates errors without real database connections
- **User Operations**: Returns predefined errors for testing error handling paths

This comprehensive test suite ensures that all error handling patterns in the project work correctly and maintain their behavior during refactoring.
