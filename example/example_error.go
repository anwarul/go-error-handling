package example

import (
	"errors"
	"fmt"
	"go-error-handling/basic"
	"go-error-handling/custom"
	"go-error-handling/database"
	"go-error-handling/formatted"
	"go-error-handling/user"
	"go-error-handling/utils"
	"go-error-handling/wrapping"
	"log"
	"os"
)

// Example 1.1: Simple error creation and checking
func BasicErrorExample() {
	result, err := basic.Divide(10, 0)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Result: %.2f\n", result)
}

// Example 2.1: Using custom error types
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

// Example 2.1: Formatted error creation and checking
func FormattedErrorExample(age int) {
	err := formatted.ValidateAge(age)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Valid age: %d\n", age)
}

// Example 3.1: Error wrapping and unwrapping
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

// Example 4.1: Using sentinel errors for expected failures
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

// Example 5.1: Rich error types with metadata
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
