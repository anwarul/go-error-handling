package database

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestDatabaseError_Error(t *testing.T) {
	baseErr := errors.New("connection timeout")
	timestamp := time.Date(2023, 10, 5, 10, 30, 0, 0, time.UTC)

	err := &DatabaseError{
		Operation: "SELECT",
		Table:     "users",
		Query:     "SELECT * FROM users WHERE id = ?",
		Err:       baseErr,
		Timestamp: timestamp,
		Retryable: true,
	}

	result := err.Error()

	// Check that all components are in the error message
	expectedComponents := []string{
		"database error",
		"SELECT",
		"users",
		"connection timeout",
		"retryable: true",
		"2023-10-05T10:30:00Z",
	}

	for _, component := range expectedComponents {
		if !strings.Contains(result, component) {
			t.Errorf("DatabaseError.Error() should contain '%s', got: %s", component, result)
		}
	}
}

func TestDatabaseError_Unwrap(t *testing.T) {
	baseErr := errors.New("original error")
	dbErr := &DatabaseError{
		Operation: "INSERT",
		Table:     "products",
		Err:       baseErr,
		Timestamp: time.Now(),
		Retryable: false,
	}

	unwrapped := dbErr.Unwrap()
	if unwrapped != baseErr {
		t.Errorf("DatabaseError.Unwrap() = %v; want %v", unwrapped, baseErr)
	}

	// Test that errors.Is works with wrapped error
	if !errors.Is(dbErr, baseErr) {
		t.Error("errors.Is should find the wrapped error")
	}
}

func TestDatabaseError_Fields(t *testing.T) {
	baseErr := errors.New("test error")
	timestamp := time.Now()

	dbErr := &DatabaseError{
		Operation: "UPDATE",
		Table:     "orders",
		Query:     "UPDATE orders SET status = ?",
		Err:       baseErr,
		Timestamp: timestamp,
		Retryable: true,
	}

	if dbErr.Operation != "UPDATE" {
		t.Errorf("DatabaseError.Operation = %v; want %v", dbErr.Operation, "UPDATE")
	}
	if dbErr.Table != "orders" {
		t.Errorf("DatabaseError.Table = %v; want %v", dbErr.Table, "orders")
	}
	if dbErr.Query != "UPDATE orders SET status = ?" {
		t.Errorf("DatabaseError.Query = %v; want %v", dbErr.Query, "UPDATE orders SET status = ?")
	}
	if dbErr.Err != baseErr {
		t.Errorf("DatabaseError.Err = %v; want %v", dbErr.Err, baseErr)
	}
	if !dbErr.Timestamp.Equal(timestamp) {
		t.Errorf("DatabaseError.Timestamp = %v; want %v", dbErr.Timestamp, timestamp)
	}
	if !dbErr.Retryable {
		t.Errorf("DatabaseError.Retryable = %v; want %v", dbErr.Retryable, true)
	}
}

func TestDatabaseError_RetryableFlag(t *testing.T) {
	tests := []struct {
		name      string
		retryable bool
	}{
		{"retryable error", true},
		{"non-retryable error", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbErr := &DatabaseError{
				Operation: "DELETE",
				Table:     "temp",
				Err:       errors.New("test"),
				Timestamp: time.Now(),
				Retryable: tt.retryable,
			}

			if dbErr.Retryable != tt.retryable {
				t.Errorf("DatabaseError.Retryable = %v; want %v", dbErr.Retryable, tt.retryable)
			}

			// Check that the retryable flag appears in the error message
			errMsg := dbErr.Error()
			expectedText := fmt.Sprintf("retryable: %v", tt.retryable)
			if !strings.Contains(errMsg, expectedText) {
				t.Errorf("Error message should contain '%s', got: %s", expectedText, errMsg)
			}
		})
	}
}

func TestUnwramp_Function(t *testing.T) {
	// Test with an error that implements Unwrap
	baseErr := errors.New("base error")
	dbErr := &DatabaseError{
		Operation: "SELECT",
		Table:     "test",
		Err:       baseErr,
		Timestamp: time.Now(),
		Retryable: false,
	}

	unwrapped := Unwramp(dbErr)
	if unwrapped != baseErr {
		t.Errorf("Unwramp(dbErr) = %v; want %v", unwrapped, baseErr)
	}

	// Test with an error that doesn't implement Unwrap
	simpleErr := errors.New("simple error")
	result := Unwramp(simpleErr)
	if result != nil {
		t.Errorf("Unwramp(simpleErr) = %v; want nil", result)
	}
}

func TestDatabaseError_AsError(t *testing.T) {
	dbErr := &DatabaseError{
		Operation: "INSERT",
		Table:     "users",
		Err:       errors.New("constraint violation"),
		Timestamp: time.Now(),
		Retryable: false,
	}

	// Test that it can be used as an error interface
	var e error = dbErr
	if e.Error() == "" {
		t.Error("DatabaseError should implement error interface")
	}

	// Test errors.As functionality
	var target *DatabaseError
	if !errors.As(e, &target) {
		t.Error("errors.As should be able to extract DatabaseError")
	}
	if target.Operation != "INSERT" {
		t.Errorf("Extracted DatabaseError.Operation = %v; want %v", target.Operation, "INSERT")
	}
}
