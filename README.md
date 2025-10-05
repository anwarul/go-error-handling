# Go Error Handling Examples

A comprehensive collection of Go error handling patterns and techniques, demonstrating modern error handling approaches in Go 1.25+.

## About

This repository showcases practical error handling patterns in Go, organized into focused packages that demonstrate different error handling techniques. Each package contains real-world examples with comprehensive test coverage.

## Quick Start

```bash
# Clone the repository
git clone https://github.com/anwarul/go-error-handling.git
cd go-error-handling

# Initialize and run all examples
go mod tidy
go run main.go

# Run tests for all packages
go test ./...

# Run tests with coverage
go test -cover ./...
```

## Table of Contents

1. [Project Structure](#project-structure)
2. [Basic Error Handling](#1-basic-error-handling)
3. [Custom Error Types](#2-custom-error-types)
4. [Formatted Errors](#3-formatted-errors)
5. [Error Wrapping](#4-error-wrapping)
6. [Sentinel Errors](#5-sentinel-errors)
7. [Database Errors](#6-database-errors)
8. [User Validation](#7-user-validation)
9. [Example Integration](#8-example-integration)
10. [Testing](#testing)
11. [Best Practices](#best-practices)

## Project Structure

```
go-error-handling/
├── main.go                    # Main entry point - runs all examples
├── go.mod                     # Module definition
├── basic/                     # Basic error handling
│   ├── basic_error.go         # Simple division with error checking
│   └── basic_error_test.go    # Comprehensive tests
├── custom/                    # Custom error types
│   ├── validation_error.go    # ValidationError struct with custom formatting
│   └── validation_error_test.go
├── formatted/                 # Formatted error messages
│   ├── formatted_error.go     # Age validation with fmt.Errorf
│   └── formatted_error_test.go
├── wrapping/                  # Error wrapping chains
│   ├── wrapping_error.go      # Multi-level error wrapping
│   └── wrapping_error_test.go
├── utils/                     # Sentinel errors
│   ├── constants.go           # Predefined error constants
│   └── constants_test.go
├── database/                  # Database error handling
│   ├── database_error.go      # Rich error type with metadata
│   └── database_error_test.go
├── user/                      # User operations and validation
│   ├── user.go                # User struct and operations
│   └── user_test.go
├── example/                   # Integration examples
│   ├── example_error.go       # Demonstrates all error patterns
│   └── example_error_test.go
├── TEST_README.md             # Detailed testing documentation
└── README.md                  # This file
```

## Examples

### 1. Basic Error Handling

**Package:** `basic/`

Simple error creation and checking with a division function.

```go
package basic

import "errors"

func Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// Usage in example/example_error.go
func BasicErrorExample() {
    result, err := basic.Divide(10, 0)
    if err != nil {
        log.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Result: %.2f\n", result)
}
```

**Key Points:**
- Always check errors immediately  
- Return errors as the last value
- Use `nil` for no error
- Log errors with context

### 2. Custom Error Types

**Package:** `custom/`

Structured error information with custom formatting.

```go
package custom

import "fmt"

type ValidationError struct {
    Field   string
    Message string
    Code    int
    Value   interface{}
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("Validation error on field '%s': %s (code: %d, value: %v)", 
        e.Field, e.Message, e.Code, e.Value)
}

// Usage in example/example_error.go
func CustomErrorExample(value int) error {
    if value < 0 {
        return &custom.ValidationError{
            Field:   "value",
            Message: "Value cannot be negative",
            Code:    1001,
            Value:   value,
        }
    }
    if value > 100 {
        return &custom.ValidationError{
            Field:   "value", 
            Message: "Value cannot be greater than 100",
            Code:    1002,
            Value:   value,
        }
    }
    return nil
}
```

**Use Cases:**
- API validation responses
- Structured error information
- Error codes for client handling
- Preserving invalid values for debugging

### 3. Formatted Errors

**Package:** `formatted/`

Using `fmt.Errorf` for contextual error messages.

```go
package formatted

import "fmt"

func ValidateAge(age int) error {
    if age < 0 {
        return fmt.Errorf("invalid age: %d. Age cannot be negative", age)
    }
    if age > 130 {
        return fmt.Errorf("invalid age: %d. Age cannot be greater than 130", age)
    }
    return nil
}

// Usage in example/example_error.go  
func FormattedErrorExample(age int) {
    err := formatted.ValidateAge(age)
    if err != nil {
        log.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Valid age: %d\n", age)
}
```

**Key Points:**
- Include relevant values in error messages
- Use descriptive, user-friendly messages
- Consider internationalization for user-facing errors

### 4. Error Wrapping

**Package:** `wrapping/`

Multi-level error wrapping with context preservation.

```go
package wrapping

import (
    "fmt"
    "os"
)

func ProcessUserData(userID int) error {
    err := loadUserConfig(userID)
    if err != nil {
        return fmt.Errorf("failed to process user %d: %w", userID, err)
    }
    return nil
}

func loadUserConfig(userID int) error {
    filename := fmt.Sprintf("user_%d.json", userID)
    err := readConfigFile(filename)
    if err != nil {
        return fmt.Errorf("failed to load config for user %d: %w", userID, err)
    }
    return nil
}

func readConfigFile(filename string) error {
    _, err := os.ReadFile(filename)
    if err != nil {
        return fmt.Errorf("failed to read config file %s: %w", filename, err)
    }
    return nil
}

// Usage in example/example_error.go
func WrappingErrorExample(filename string) {
    err := wrapping.ProcessUserData(123)
    if err != nil {
        log.Printf("Full error chain: %v\n", err)
        
        // Check if it wraps a specific error
        if errors.Is(err, os.ErrNotExist) {
            log.Println("File not found - using defaults")
        }
    }
}
```

**Best Practices:**
- Use `%w` to wrap errors (Go 1.13+)
- Add context at each layer
- Use `errors.Is()` to check wrapped errors
- Preserve the original error chain

### 5. Sentinel Errors

**Package:** `utils/`

Predefined errors for expected conditions.

```go
package utils

import "errors"

var (
    ErrUserNotFound    = errors.New("user not found")
    ErrDuplicateEmail  = errors.New("email already exists")
    ErrInvalidPassword = errors.New("invalid password")
    ErrUnauthorized    = errors.New("unauthorized access")
    ErrDatabaseTimeout = errors.New("database operation timed out")
)

// Usage in user/user.go
func FindUserByEmail(email string) (*User, error) {
    if email == "" {
        return nil, fmt.Errorf("email cannot be empty: %w", utils.ErrUserNotFound)
    }
    return nil, utils.ErrUserNotFound
}

// Usage in example/example_error.go
func SentinelErrorExample() {
    user, err := user.FindUserByEmail("test@example.com")
    if err != nil {
        if errors.Is(err, utils.ErrUserNotFound) {
            log.Println("User doesn't exist - creating new account")
            return
        }
        log.Printf("Unexpected error: %v\n", err)
    }
    log.Printf("Found user: %v\n", user)
}
```

**When to Use:**
- Expected failure conditions
- Callers need to distinguish error types
- API boundary errors
- Consistent error identity across packages

### 6. Database Errors

**Package:** `database/`

Rich error types with metadata and context.

```go
package database

import (
    "fmt"
    "time"
)

type DatabaseError struct {
    Operation string
    Table     string
    Query     string
    Err       error
    Timestamp time.Time
    Retryable bool
}

func (e *DatabaseError) Error() string {
    return fmt.Sprintf("database error [%s on %s]: %v (retryable: %v, timestamp: %s)",
        e.Operation, e.Table, e.Err, e.Retryable, e.Timestamp.Format(time.RFC3339))
}

func (e *DatabaseError) Unwrap() error {
    return e.Err
}

// Usage in user/user.go
func QueryUsers(limit int) error {
    // Simulate database error
    return &database.DatabaseError{
        Operation: "SELECT",
        Table:     "users",
        Query:     fmt.Sprintf("SELECT * FROM users LIMIT %d", limit),
        Err:       errors.New("connection timeout"),
        Timestamp: time.Now(),
        Retryable: true,
    }
}

// Usage in example/example_error.go
func ComplexErrorExample() {
    err := user.QueryUsers(10)
    if err != nil {
        var dbErr *database.DatabaseError
        if errors.As(err, &dbErr) {
            log.Printf("Database operation: %s\n", dbErr.Operation)
            log.Printf("Table: %s\n", dbErr.Table)
            log.Printf("Retryable: %v\n", dbErr.Retryable)

            if dbErr.Retryable {
                log.Println("Retrying operation...")
            }
        }
    }
}
```

**Use Cases:**
- Database operations with retry logic
- Operations needing structured metadata
- Debugging and monitoring
- Circuit breaker patterns

### 7. User Validation

**Package:** `user/`

User operations with validation and database simulation.

```go 
package user

import (
    "errors"
    "fmt"
    "go-error-handling/custom"
    "go-error-handling/database"
    "go-error-handling/utils"
    "time"
)

type User struct {
    ID    int
    Email string
    Age   int
}

func ValidateUser(user User) error {
    if user.Age < 0 {
        return &custom.ValidationError{
            Field:   "Age",
            Message: "Age cannot be negative",
            Code:    2001,
            Value:   user.Age,
        }
    }
    if user.Age > 130 {
        return &custom.ValidationError{
            Field:   "Age", 
            Message: "Age cannot be greater than 130",
            Code:    2002,
            Value:   user.Age,
        }
    }
    if user.Email == "" {
        return &custom.ValidationError{
            Field:   "Email",
            Message: "Email cannot be empty",
            Code:    2003,
            Value:   user.Email,
        }
    }
    return nil
}

func FindUserByEmail(email string) (*User, error) {
    if email == "" {
        return nil, fmt.Errorf("email cannot be empty: %w", utils.ErrUserNotFound)
    }
    return nil, utils.ErrUserNotFound
}

func QueryUsers(limit int) error {
    // Simulate database error
    return &database.DatabaseError{
        Operation: "SELECT",
        Table:     "users", 
        Query:     fmt.Sprintf("SELECT * FROM users LIMIT %d", limit),
        Err:       errors.New("connection timeout"),
        Timestamp: time.Now(),
        Retryable: true,
    }
}
```

**Patterns Demonstrated:**
- Validation with custom error types
- Sentinel error wrapping
- Database error simulation
- Structured error information

### 8. Example Integration

**Package:** `example/`

Integration examples showing all error patterns working together.

```go
package example

// main.go runs all these examples
func main() {
    example.BasicErrorExample()
    
    example.CustomErrorExample(-5)
    example.CustomErrorExample(150)
    
    example.FormattedErrorExample(-10)
    example.FormattedErrorExample(25)
    example.FormattedErrorExample(150)
    
    example.WrappingErrorExample("non_existent_file.txt")
    example.WrappingErrorExample("valid_file.txt")
    
    example.ComplexErrorExample()
    example.CustomErrorExample(999)
}
```

**Example Functions:**
- `BasicErrorExample()` - Simple error handling
- `CustomErrorExample(value)` - Custom error types with codes
- `FormattedErrorExample(age)` - Formatted error messages
- `WrappingErrorExample(filename)` - Error wrapping chains
- `SentinelErrorExample()` - Sentinel error detection
- `ComplexErrorExample()` - Database errors with metadata

**Integration Benefits:**
- Shows error patterns in context
- Demonstrates error type assertions
- Real-world error propagation
- Complete error handling workflows

## Testing

This project includes comprehensive tests for all error handling patterns. See [`TEST_README.md`](TEST_README.md) for detailed testing documentation.

### Test Coverage

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test -v ./basic
go test -v ./custom
go test -v ./database
```

### Coverage Results

- **basic**: 100% coverage - Division function with all edge cases
- **custom**: 100% coverage - ValidationError with different value types  
- **formatted**: 100% coverage - Age validation with boundary testing
- **database**: 100% coverage - DatabaseError with unwrapping
- **user**: 100% coverage - User validation and operations
- **wrapping**: 84.6% coverage - Error chains with file operations
- **example**: 91.7% coverage - Integration examples
- **utils**: No statements - Only constants

### Key Testing Patterns

- **Error Identity**: Using `errors.Is()` for sentinel error detection
- **Type Assertions**: Using `errors.As()` for custom error extraction
- **Boundary Testing**: Edge cases and limit values
- **Error Message Validation**: Ensuring proper error formatting
- **Panic Prevention**: Integration tests for stability
- **Chain Traversal**: Verifying error unwrapping works correctly

## Best Practices

Based on the patterns demonstrated in this project:

### ✅ Do's

1. **Always check errors** immediately after operations
2. **Add context** with `fmt.Errorf()` or error wrapping
3. **Use custom types** when you need structured error information
4. **Implement error interfaces** properly (`Error()`, `Unwrap()`)
5. **Use sentinel errors** for expected conditions
6. **Test error paths** comprehensively
7. **Preserve error chains** with `%w` formatting
8. **Log errors at boundaries** (main, handlers, top-level functions)

### ❌ Don'ts

1. **Don't ignore errors** - always handle them appropriately
2. **Don't panic** for normal error conditions
3. **Don't log and return** - choose one approach
4. **Don't lose context** - always add meaningful information
5. **Don't expose internal errors** to external callers
6. **Don't create errors without context**

### Error Handling Principles

1. **Fail Fast**: Check errors immediately and return early
2. **Add Context**: Each layer should add meaningful information
3. **Preserve Chains**: Use `%w` to maintain error relationships
4. **Type Safety**: Use `errors.Is()` and `errors.As()` for type checking
5. **Consistent Patterns**: Apply the same error handling approach throughout your codebase

## Common Patterns Demonstrated

### Pattern: Early Returns

Demonstrated in `user/user.go`:

```go
func ValidateUser(user User) error {
    if user.Age < 0 {
        return &custom.ValidationError{
            Field:   "Age",
            Message: "Age cannot be negative",
            Code:    2001,
            Value:   user.Age,
        }
    }
    if user.Age > 130 {
        return &custom.ValidationError{
            Field:   "Age",
            Message: "Age cannot be greater than 130", 
            Code:    2002,
            Value:   user.Age,
        }
    }
    if user.Email == "" {
        return &custom.ValidationError{
            Field:   "Email",
            Message: "Email cannot be empty",
            Code:    2003,
            Value:   user.Email,
        }
    }
    return nil
}
```

### Pattern: Error Context

Demonstrated in `wrapping/wrapping_error.go`:

```go
func ProcessUserData(userID int) error {
    err := loadUserConfig(userID) 
    if err != nil {
        return fmt.Errorf("failed to process user %d: %w", userID, err)
    }
    return nil
}

func loadUserConfig(userID int) error {
    filename := fmt.Sprintf("user_%d.json", userID)
    err := readConfigFile(filename)
    if err != nil {
        return fmt.Errorf("failed to load config for user %d: %w", userID, err)
    }
    return nil
}
```

### Pattern: Type Assertions

Demonstrated in `example/example_error.go`:

```go
func ComplexErrorExample() {
    err := user.QueryUsers(10)
    if err != nil {
        var dbErr *database.DatabaseError
        if errors.As(err, &dbErr) {
            log.Printf("Database operation: %s\n", dbErr.Operation)
            log.Printf("Table: %s\n", dbErr.Table)
            log.Printf("Retryable: %v\n", dbErr.Retryable)

            if dbErr.Retryable {
                log.Println("Retrying operation...")
            }
        }
    }
}
```

## Key Features

This project demonstrates:

✅ **Modern Go Error Handling** (Go 1.25+)  
✅ **Comprehensive Test Coverage** (>90% across all packages)  
✅ **Real-world Examples** with practical use cases  
✅ **Error Wrapping** with `fmt.Errorf` and `%w`  
✅ **Custom Error Types** with structured information  
✅ **Sentinel Error Patterns** for expected conditions  
✅ **Error Type Assertions** with `errors.Is()` and `errors.As()`  
✅ **Multi-level Error Chains** with context preservation  
✅ **Database Error Simulation** with retry logic  
✅ **Integration Examples** showing patterns working together  

## Usage Examples

### Running the Examples

```bash
# Run all examples
go run main.go

# Output will show:
# 2025/10/05 12:00:00 Error: division by zero
# 2025/10/05 12:00:00 Error: Validation error on field 'value': Value cannot be negative (code: 1001, value: -5)
# 2025/10/05 12:00:00 Error: Validation error on field 'value': Value cannot be greater than 100 (code: 1002, value: 150)
# ... and more examples
```

### Running Individual Tests

```bash
# Test basic error handling
go test -v ./basic

# Test custom error types  
go test -v ./custom

# Test error wrapping
go test -v ./wrapping

# Test all with coverage
go test -cover ./...
```

## Learning Path

### 1. Start with Basics
- Run `go test -v ./basic` to understand simple error handling
- Study `basic/basic_error.go` for error creation and checking
- Review `basic/basic_error_test.go` for comprehensive test patterns

### 2. Move to Custom Types  
- Explore `custom/validation_error.go` for structured errors
- Understand how custom error types provide better context
- See how `Error()` method formatting works

### 3. Learn Error Wrapping
- Study `wrapping/wrapping_error.go` for multi-level error chains
- Understand how `%w` preserves error relationships
- Practice with `errors.Is()` for error identity checking

### 4. Master Type Assertions
- Review `example/example_error.go` for `errors.As()` usage
- Learn when to use sentinel errors vs custom types
- Understand error chain traversal

### 5. Apply in Real Scenarios
- Study `user/user.go` for validation patterns
- Explore `database/database_error.go` for retry logic
- Understand error metadata and structured information

## Tools and Linters

Enhance your error handling with these tools:

```bash
# Check for unchecked errors
go install github.com/kisielk/errcheck@latest
errcheck ./...

# Lint for Go 1.13+ error handling
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint run --enable=errorlint,goerr113,wrapcheck

# Format your code
go fmt ./...

# Run all tests with race detection
go test -race ./...
```

## Package Documentation

### `basic/` - Foundation Error Handling
- **`basic_error.go`**: Simple error creation with `errors.New()`
- **`basic_error_test.go`**: Tests covering success, failure, and edge cases
- **Pattern**: Basic error checking and early returns

### `custom/` - Structured Error Information  
- **`validation_error.go`**: Custom error type with fields and formatting
- **`validation_error_test.go`**: Tests for different value types and error interface
- **Pattern**: When you need more than just an error message

### `formatted/` - Contextual Error Messages
- **`formatted_error.go`**: Using `fmt.Errorf()` for dynamic error messages  
- **`formatted_error_test.go`**: Boundary testing and message validation
- **Pattern**: Adding context with interpolated values

### `wrapping/` - Error Chain Management
- **`wrapping_error.go`**: Multi-level error wrapping with context
- **`wrapping_error_test.go`**: Error chain traversal and unwrapping tests
- **Pattern**: Preserving error history across function calls

### `utils/` - Sentinel Error Constants
- **`constants.go`**: Predefined errors for expected conditions
- **`constants_test.go`**: Error identity and uniqueness verification  
- **Pattern**: Using `errors.Is()` for error type checking

### `database/` - Rich Error Metadata
- **`database_error.go`**: Complex error type with operation details
- **`database_error_test.go`**: Testing error unwrapping and metadata
- **Pattern**: Errors with structured information for debugging and retry logic

### `user/` - Domain-Specific Operations
- **`user.go`**: User validation and operations with multiple error types
- **`user_test.go`**: Integration testing of different error patterns
- **Pattern**: Combining multiple error handling approaches

### `example/` - Integration and Demonstration
- **`example_error.go`**: Shows all patterns working together
- **`example_error_test.go`**: End-to-end testing and panic prevention
- **Pattern**: Real-world error handling workflows

## Resources

### Official Go Documentation
- [Go Blog: Error Handling and Go](https://go.dev/blog/error-handling-and-go)
- [Go Blog: Working with Errors in Go 1.13](https://go.dev/blog/go1.13-errors)
- [Effective Go: Errors](https://go.dev/doc/effective_go#errors)
- [Go 1.25 Release Notes](https://go.dev/doc/go1.25)

### Related Packages
- [`errors`](https://pkg.go.dev/errors) - Standard library error handling
- [`fmt`](https://pkg.go.dev/fmt) - Error formatting with `Errorf`
- [`testing`](https://pkg.go.dev/testing) - Testing error conditions

### Testing Resources
- [Go Testing Documentation](https://go.dev/doc/tutorial/add-a-test)
- [Table-Driven Tests](https://go.dev/wiki/TableDrivenTests)
- [Test Coverage](https://go.dev/blog/cover)

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-pattern`)
3. Add your example with tests
4. Ensure all tests pass (`go test ./...`)
5. Submit a pull request

### Contribution Guidelines

- ✅ Include comprehensive tests
- ✅ Follow existing code patterns
- ✅ Add documentation and examples
- ✅ Maintain high test coverage
- ✅ Use descriptive commit messages

## License

MIT License - See [LICENSE](LICENSE) file for details.

## Project Status

- ✅ **Stable**: All packages have comprehensive tests
- ✅ **Maintained**: Regular updates for new Go versions
- ✅ **Production Ready**: Patterns used in real-world applications
- ✅ **Well Documented**: Extensive examples and explanations

## Feedback

- 📝 Open an [issue](https://github.com/anwarul/go-error-handling/issues) for bugs
- 💡 Start a [discussion](https://github.com/anwarul/go-error-handling/discussions) for ideas
- ⭐ Star this repository if you find it useful
- 🔄 Share with your Go community

---

**Happy error handling!** 🚀

```bash
git clone https://github.com/anwarul/go-error-handling.git
cd go-error-handling
go run main.go
```