package example

import (
	"errors"
	"fmt"
	"go-error-handling/custom"
	"go-error-handling/database"
	"go-error-handling/user"
	"go-error-handling/utils"
	"go-error-handling/wrapping"
	"os"
	"testing"
)

func TestBasicErrorExample_DoesNotPanic(t *testing.T) {
	// This test ensures BasicErrorExample doesn't panic
	// Since it uses log.Printf, we can't easily capture output,
	// but we can ensure it completes without panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("BasicErrorExample() panicked: %v", r)
		}
	}()

	BasicErrorExample()
}

func TestCustomErrorExample_NegativeValue(t *testing.T) {
	testValues := []int{-1, -5, -10, -100}

	for _, value := range testValues {
		t.Run(fmt.Sprintf("negative_%d", value), func(t *testing.T) {
			err := CustomErrorExample(value)

			if err == nil {
				t.Errorf("CustomErrorExample(%d) expected error but got none", value)
			}

			var validationErr *custom.ValidationError
			if !errors.As(err, &validationErr) {
				t.Errorf("Expected ValidationError, got %T", err)
			}

			if validationErr.Field != "value" {
				t.Errorf("Expected field 'value', got '%s'", validationErr.Field)
			}
			if validationErr.Code != 1001 {
				t.Errorf("Expected code 1001, got %d", validationErr.Code)
			}
			if validationErr.Value != value {
				t.Errorf("Expected value %d, got %v", value, validationErr.Value)
			}
		})
	}
}

func TestCustomErrorExample_TooLargeValue(t *testing.T) {
	testValues := []int{101, 150, 200, 1000}

	for _, value := range testValues {
		t.Run(fmt.Sprintf("too_large_%d", value), func(t *testing.T) {
			err := CustomErrorExample(value)

			if err == nil {
				t.Errorf("CustomErrorExample(%d) expected error but got none", value)
			}

			var validationErr *custom.ValidationError
			if !errors.As(err, &validationErr) {
				t.Errorf("Expected ValidationError, got %T", err)
			}

			if validationErr.Field != "value" {
				t.Errorf("Expected field 'value', got '%s'", validationErr.Field)
			}
			if validationErr.Code != 1002 {
				t.Errorf("Expected code 1002, got %d", validationErr.Code)
			}
			if validationErr.Value != value {
				t.Errorf("Expected value %d, got %v", value, validationErr.Value)
			}
		})
	}
}

func TestCustomErrorExample_ValidValue(t *testing.T) {
	testValues := []int{0, 1, 50, 99, 100}

	for _, value := range testValues {
		t.Run(fmt.Sprintf("valid_%d", value), func(t *testing.T) {
			err := CustomErrorExample(value)

			if err != nil {
				t.Errorf("CustomErrorExample(%d) expected no error but got: %v", value, err)
			}
		})
	}
}

func TestFormattedErrorExample_DoesNotPanic(t *testing.T) {
	testAges := []int{-10, 25, 150}

	for _, age := range testAges {
		t.Run(fmt.Sprintf("age_%d", age), func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("FormattedErrorExample(%d) panicked: %v", age, r)
				}
			}()

			FormattedErrorExample(age)
		})
	}
}

func TestWrappingErrorExample_DoesNotPanic(t *testing.T) {
	testFilenames := []string{"non_existent_file.txt", "valid_file.txt"}

	for _, filename := range testFilenames {
		t.Run(filename, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("WrappingErrorExample(%s) panicked: %v", filename, r)
				}
			}()

			WrappingErrorExample(filename)
		})
	}
}

func TestWrappingErrorExample_ErrorDetection(t *testing.T) {
	// Since WrappingErrorExample calls wrapping.ProcessUserData(123),
	// we know it will always fail with os.ErrNotExist
	// This test verifies the error detection logic works

	// We can't directly test the log output, but we can test the underlying logic
	err := wrapping.ProcessUserData(123)
	if err == nil {
		t.Error("wrapping.ProcessUserData(123) should return an error")
	}

	// The function should detect os.ErrNotExist in the error chain
	if !errors.Is(err, os.ErrNotExist) {
		t.Error("wrapping.ProcessUserData error should wrap os.ErrNotExist")
	}
}

func TestSentinelErrorExample_DoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("SentinelErrorExample() panicked: %v", r)
		}
	}()

	SentinelErrorExample()
}

func TestSentinelErrorExample_ErrorDetection(t *testing.T) {
	// Test the underlying logic used in SentinelErrorExample
	_, err := user.FindUserByEmail("test@example.com")

	if err == nil {
		t.Error("user.FindUserByEmail should return an error")
	}

	// Should detect ErrUserNotFound
	if !errors.Is(err, utils.ErrUserNotFound) {
		t.Error("user.FindUserByEmail error should be or wrap ErrUserNotFound")
	}
}

func TestComplexErrorExample_DoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ComplexErrorExample() panicked: %v", r)
		}
	}()

	ComplexErrorExample()
}

func TestComplexErrorExample_ErrorTypeAssertion(t *testing.T) {
	// Test the underlying logic used in ComplexErrorExample
	err := user.QueryUsers(10)

	if err == nil {
		t.Error("user.QueryUsers should return an error")
	}

	// Should be able to extract DatabaseError
	var dbErr *database.DatabaseError
	if !errors.As(err, &dbErr) {
		t.Errorf("user.QueryUsers error should be DatabaseError, got %T", err)
	}

	// Verify the DatabaseError properties
	if dbErr.Operation != "SELECT" {
		t.Errorf("Expected operation 'SELECT', got '%s'", dbErr.Operation)
	}
	if dbErr.Table != "users" {
		t.Errorf("Expected table 'users', got '%s'", dbErr.Table)
	}
	if !dbErr.Retryable {
		t.Error("Database error should be retryable")
	}
}

func TestAllExampleFunctions_Integration(t *testing.T) {
	// Integration test to ensure all example functions can run together
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Integration test panicked: %v", r)
		}
	}()

	// Run all the example functions like in main.go
	BasicErrorExample()

	CustomErrorExample(-5)
	CustomErrorExample(150)

	FormattedErrorExample(-10)
	FormattedErrorExample(25)
	FormattedErrorExample(150)

	WrappingErrorExample("non_existent_file.txt")
	WrappingErrorExample("valid_file.txt")

	ComplexErrorExample()
	CustomErrorExample(999)
}
