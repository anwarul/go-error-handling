# Contributing to Go Error Handling Examples

Thanks for your interest in contributing! This repository aims to collect real-world, battle-tested error handling patterns.

## What We're Looking For

- **Production-proven patterns** - not theoretical examples
- **Clear use cases** - explain when and why to use each pattern
- **Complete examples** - code that actually runs
- **Anti-patterns** - common mistakes you've made (we all have)
- **Performance insights** - when patterns matter for performance
- **Migration guides** - how you evolved your error handling

## What We're Not Looking For

- Academic exercises without real-world context
- Overly complex patterns that don't add value
- Frameworks that wrap standard error handling
- Patterns that violate Go conventions

## How to Contribute

### 1. Reporting Issues

Found a bug or unclear example?

- Open an issue with a clear title
- Describe what you expected vs what you got
- Include Go version and OS if relevant
- Bonus: include a minimal reproduction

### 2. Suggesting Patterns

Have a pattern that's worked well for you?

- Open an issue first to discuss
- Explain the problem it solves
- Share your production experience with it
- Mention any trade-offs or gotchas

### 3. Submitting Code

Ready to add an example?

1. **Fork the repository**
   ```bash
   git clone https://github.com/anwarul/go-error-handling.git
   cd go-error-handling
   ```

2. **Create a feature branch**
   ```bash
   git checkout -b feature/your-pattern-name
   ```

3. **Add your example**
   - Create a new package directory or extend an existing one (`basic/`, `custom/`, `formatted/`, `wrapping/`, `utils/`, `database/`, `user/`, `example/`)
   - Include comprehensive tests that demonstrate the pattern
   - Add integration examples in the `example/` package
   - Update main README.md if needed

4. **Follow the style guide** (see below)

5. **Test your code**
   ```bash
   # Run all tests
   go test ./...
   
   # Run with coverage (aim for >90%)
   go test -cover ./...
   
   # Run with verbose output to see test details
   go test -v ./...
   
   # Check code quality
   go vet ./...
   go fmt ./...
   
   # If you have golangci-lint installed
   golangci-lint run
   ```

6. **Commit with clear messages**
   ```bash
   git commit -m "Add pattern: [brief description]"
   ```

7. **Push and create PR**
   ```bash
   git push origin feature/your-pattern-name
   ```

## Code Style Guide

### Error Messages

```go
// Good - lowercase, no punctuation
return errors.New("connection failed")
return fmt.Errorf("user %d not found", id)

// Bad - capitalized, with punctuation
return errors.New("Connection failed.")
return fmt.Errorf("User %d not found!", id)
```

### Error Wrapping

```go
// Good - always use %w when preserving the error
return fmt.Errorf("failed to process order: %w", err)

// Bad - using %v loses error chain
return fmt.Errorf("failed to process order: %v", err)
```

### Comments

```go
// Good - explains why, not what
// Retry with backoff because database connections
// often fail transiently under load
func retryOperation() error { }

// Bad - states the obvious
// This function retries an operation
func retryOperation() error { }
```

### Package Structure

Each package should follow this structure:

```go
// Package comment explaining the pattern
package yourpackage

import (
    // Standard library first
    "errors"
    "fmt"
    
    // Internal packages
    "go-error-handling/utils"
    "go-error-handling/custom"
)

// Clear function with explanation
func YourFunction() error {
    // Implementation demonstrating the pattern
    
    // Return appropriate error or nil
}
```

**Package Organization:**
- `basic/` - Foundation error handling patterns
- `custom/` - Custom error types with structured information  
- `formatted/` - Contextual error messages with `fmt.Errorf`
- `wrapping/` - Error chain management and context preservation
- `utils/` - Sentinel errors and constants
- `database/` - Rich error metadata for operations
- `user/` - Domain-specific validation and operations
- `example/` - Integration examples using all patterns

### Documentation

Every contribution needs:

1. **Package comment** - what problem does this solve?
2. **Function comments** - when to use this pattern
3. **Inline comments** - explain non-obvious decisions
4. **Comprehensive tests** - see existing `*_test.go` files for examples
5. **Integration example** - add usage to `example/example_error.go`
6. **Update main README.md** - add your pattern to the appropriate section

## Testing Guidelines

### Write Tests That Demonstrate Value

Follow the existing test patterns in the project:

```go
// Good - comprehensive testing like in wrapping/wrapping_error_test.go
func TestProcessUserData_FileNotFound(t *testing.T) {
    userID := 999
    
    err := ProcessUserData(userID)
    
    if err == nil {
        t.Errorf("ProcessUserData(%d) expected error but got none", userID)
    }

    // Check that the error wraps os.ErrNotExist
    if !errors.Is(err, os.ErrNotExist) {
        t.Error("ProcessUserData error should wrap os.ErrNotExist")
    }
}

// Good - table-driven tests like in custom/validation_error_test.go
func TestValidationError_Error(t *testing.T) {
    tests := []struct {
        name        string
        err         ValidationError
        expectedMsg string
    }{
        {
            name: "string value",
            err: ValidationError{
                Field:   "username",
                Message: "Username is required",
                Code:    1001,
                Value:   "",
            },
            expectedMsg: "Validation error on field 'username': Username is required (code: 1001, value: )",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := tt.err.Error()
            if result != tt.expectedMsg {
                t.Errorf("ValidationError.Error() = %v; want %v", result, tt.expectedMsg)
            }
        })
    }
}
```

### Test Error Scenarios

Follow the comprehensive testing approach used in this project:

```go
// Example from user/user_test.go
func TestValidateUser_NegativeAge(t *testing.T) {
    invalidUsers := []User{
        {ID: 1, Email: "test@example.com", Age: -1},
        {ID: 2, Email: "user@domain.org", Age: -10},
    }

    for _, user := range invalidUsers {
        t.Run(fmt.Sprintf("negative_age_%d", user.Age), func(t *testing.T) {
            err := ValidateUser(user)
            if err == nil {
                t.Errorf("ValidateUser(%+v) expected error but got none", user)
            }

            var validationErr *custom.ValidationError
            if !errors.As(err, &validationErr) {
                t.Errorf("Expected ValidationError, got %T", err)
            }

            if validationErr.Code != 2001 {
                t.Errorf("Expected code 2001, got %d", validationErr.Code)
            }
        })
    }
}
```

### Required Test Coverage
- Aim for >90% coverage like the existing packages
- Test both success and failure cases
- Use `errors.Is()` and `errors.As()` for error type checking
- Include boundary value testing
- Test error message formatting

## Documentation Standards

### Integration with Existing Documentation

When adding new patterns:

1. **Update main README.md** - Add your pattern to the appropriate section
2. **Add to example package** - Include usage in `example/example_error.go`
3. **Update TEST_README.md** - Document your testing approach
4. **Follow existing structure** - Match the style of current packages

### Package Documentation Format

Each package should include:

```go
// Package comment explaining the error handling pattern demonstrated.
// 
// This package shows [specific pattern] which is useful when [use case].
// Key benefits include [benefit 1], [benefit 2], and [benefit 3].
//
// Example usage:
//   err := YourFunction()
//   if err != nil {
//       // Handle the error appropriately
//   }
package yourpackage
```

### Function Documentation

```go
// YourFunction demonstrates [pattern name] by [what it does].
// 
// This function returns [specific error type] when [condition].
// Use this pattern when [use case scenario].
//
// Example:
//   if err := YourFunction(); err != nil {
//       var customErr *YourErrorType
//       if errors.As(err, &customErr) {
//           // Handle custom error type
//       }
//   }
func YourFunction() error {
    // Implementation
}
```

### Code Comments

```go
// Good - explains the why
// Using a buffer because we expect high concurrency
// and want to avoid blocking goroutines
errCh := make(chan error, 100)

// Bad - explains the what (obvious from code)
// Makes an error channel with buffer of 100
errCh := make(chan error, 100)
```

## Performance Benchmarks

If your pattern has performance implications, include benchmarks:

```go
func BenchmarkErrorWrapping(b *testing.B) {
    for i := 0; i < b.N; i++ {
        err := errors.New("base error")
        _ = fmt.Errorf("wrapped: %w", err)
    }
}

func BenchmarkSentinelError(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = ErrNotFound
    }
}
```

Include results in your PR description:
```
```bash
BenchmarkErrorWrapping-8    5000000    245 ns/op
BenchmarkSentinelError-8    1000000000  0.25 ns/op
```

### Current Project Structure

When contributing, work within this structure:
```
go-error-handling/
├── main.go                    # Entry point running all examples
├── go.mod                     # Module: go-error-handling
├── basic/                     # Basic error handling
│   ├── basic_error.go
│   └── basic_error_test.go
├── custom/                    # Custom error types  
│   ├── validation_error.go
│   └── validation_error_test.go
├── formatted/                 # Formatted error messages
│   ├── formatted_error.go
│   └── formatted_error_test.go
├── wrapping/                  # Error wrapping chains
│   ├── wrapping_error.go
│   └── wrapping_error_test.go
├── utils/                     # Sentinel errors
│   ├── constants.go
│   └── constants_test.go
├── database/                  # Database error metadata
│   ├── database_error.go
│   └── database_error_test.go
├── user/                      # User operations
│   ├── user.go
│   └── user_test.go
├── example/                   # Integration examples
│   ├── example_error.go
│   └── example_error_test.go
├── TEST_README.md             # Testing documentation
└── README.md                  # Main documentation
```
```

## Commit Message Guidelines

Use clear, descriptive commit messages:

```bash
# Good
git commit -m "Add retry pattern with exponential backoff"
git commit -m "Fix race condition in concurrent error handling"
git commit -m "Update docs: clarify when to use sentinel errors"

# Bad
git commit -m "Update code"
git commit -m "Fix bug"
git commit -m "Changes"
```

## Review Process

1. **Automated checks** - must pass all tests and linters
2. **Code review** - maintainer will review for:
   - Correctness
   - Clarity
   - Real-world applicability
   - Documentation quality
3. **Revision** - address feedback if needed
4. **Merge** - once approved

## Questions?

- Open an issue for questions about contributing
- Tag issues with `question` or `discussion`
- Be respectful and constructive

## Code of Conduct

- Be respectful and inclusive
- Focus on the code, not the person
- Assume good intentions
- Share knowledge generously
- Learn from mistakes (yours and others')

## Recognition

All contributors will be:
- Listed in the repository contributors
- Credited in release notes
- Mentioned in pattern documentation

## Thank You!

Your contributions help the Go community write more reliable code. Every pattern, bug report, and suggestion makes a difference.

---

*These guidelines are based on 12 years of Go development. They'll evolve as we learn more. Suggestions welcome!*