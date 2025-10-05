package utils

import (
	"errors"
	"fmt"
	"testing"
)

func TestSentinelErrors_Identity(t *testing.T) {
	// Test that sentinel errors maintain their identity for errors.Is comparison
	tests := []struct {
		name string
		err  error
	}{
		{"ErrUserNotFound", ErrUserNotFound},
		{"ErrDuplicateEmail", ErrDuplicateEmail},
		{"ErrInvalidPassword", ErrInvalidPassword},
		{"ErrUnauthorized", ErrUnauthorized},
		{"ErrDatabaseTimeout", ErrDatabaseTimeout},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test errors.Is with the same error
			if !errors.Is(tt.err, tt.err) {
				t.Errorf("errors.Is(%v, %v) should be true", tt.err, tt.err)
			}

			// Test that wrapped versions are still detected
			wrapped := fmt.Errorf("operation failed: %w", tt.err)
			if !errors.Is(wrapped, tt.err) {
				t.Errorf("errors.Is(wrapped_error, %v) should be true", tt.err)
			}
		})
	}
}

func TestSentinelErrors_Messages(t *testing.T) {
	tests := []struct {
		name            string
		err             error
		expectedMessage string
	}{
		{"ErrUserNotFound", ErrUserNotFound, "user not found"},
		{"ErrDuplicateEmail", ErrDuplicateEmail, "email already exists"},
		{"ErrInvalidPassword", ErrInvalidPassword, "invalid password"},
		{"ErrUnauthorized", ErrUnauthorized, "unauthorized access"},
		{"ErrDatabaseTimeout", ErrDatabaseTimeout, "database operation timed out"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.expectedMessage {
				t.Errorf("%s.Error() = %v; want %v", tt.name, tt.err.Error(), tt.expectedMessage)
			}
		})
	}
}

func TestSentinelErrors_Uniqueness(t *testing.T) {
	// Test that different sentinel errors are not equal
	sentinelErrors := []error{
		ErrUserNotFound,
		ErrDuplicateEmail,
		ErrInvalidPassword,
		ErrUnauthorized,
		ErrDatabaseTimeout,
	}

	for i, err1 := range sentinelErrors {
		for j, err2 := range sentinelErrors {
			if i != j && errors.Is(err1, err2) {
				t.Errorf("Different sentinel errors should not be equal: %v vs %v", err1, err2)
			}
		}
	}
}

func TestSentinelErrors_InWrappedChain(t *testing.T) {
	// Test complex error chains
	baseErr := ErrUserNotFound
	level1 := fmt.Errorf("database layer: %w", baseErr)
	level2 := fmt.Errorf("service layer: %w", level1)
	level3 := fmt.Errorf("handler layer: %w", level2)

	if !errors.Is(level3, ErrUserNotFound) {
		t.Error("deeply wrapped ErrUserNotFound should still be detectable")
	}

	// Test that it doesn't match other sentinel errors
	if errors.Is(level3, ErrDuplicateEmail) {
		t.Error("wrapped ErrUserNotFound should not match ErrDuplicateEmail")
	}
}
