package custom

import (
	"testing"
)

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
		{
			name: "integer value",
			err: ValidationError{
				Field:   "age",
				Message: "Age must be positive",
				Code:    1002,
				Value:   -5,
			},
			expectedMsg: "Validation error on field 'age': Age must be positive (code: 1002, value: -5)",
		},
		{
			name: "float value",
			err: ValidationError{
				Field:   "price",
				Message: "Price cannot be negative",
				Code:    1003,
				Value:   -10.50,
			},
			expectedMsg: "Validation error on field 'price': Price cannot be negative (code: 1003, value: -10.5)",
		},
		{
			name: "nil value",
			err: ValidationError{
				Field:   "data",
				Message: "Data cannot be nil",
				Code:    1004,
				Value:   nil,
			},
			expectedMsg: "Validation error on field 'data': Data cannot be nil (code: 1004, value: <nil>)",
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

func TestValidationError_AsError(t *testing.T) {
	err := &ValidationError{
		Field:   "email",
		Message: "Invalid email format",
		Code:    2001,
		Value:   "invalid-email",
	}

	// Test that it can be used as an error interface
	var e error = err
	if e.Error() == "" {
		t.Error("ValidationError should implement error interface")
	}

	expectedMsg := "Validation error on field 'email': Invalid email format (code: 2001, value: invalid-email)"
	if e.Error() != expectedMsg {
		t.Errorf("ValidationError as error interface = %v; want %v", e.Error(), expectedMsg)
	}
}

func TestValidationError_Fields(t *testing.T) {
	err := ValidationError{
		Field:   "password",
		Message: "Password too short",
		Code:    3001,
		Value:   "123",
	}

	if err.Field != "password" {
		t.Errorf("ValidationError.Field = %v; want %v", err.Field, "password")
	}
	if err.Message != "Password too short" {
		t.Errorf("ValidationError.Message = %v; want %v", err.Message, "Password too short")
	}
	if err.Code != 3001 {
		t.Errorf("ValidationError.Code = %v; want %v", err.Code, 3001)
	}
	if err.Value != "123" {
		t.Errorf("ValidationError.Value = %v; want %v", err.Value, "123")
	}
}
