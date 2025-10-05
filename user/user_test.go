package user

import (
	"errors"
	"fmt"
	"go-error-handling/custom"
	"go-error-handling/database"
	"go-error-handling/utils"
	"strings"
	"testing"
)

func TestValidateUser_Success(t *testing.T) {
	validUsers := []User{
		{ID: 1, Email: "test@example.com", Age: 25},
		{ID: 2, Email: "user@domain.org", Age: 0},
		{ID: 3, Email: "senior@company.com", Age: 130},
		{ID: 4, Email: "a@b.c", Age: 65},
	}

	for _, user := range validUsers {
		t.Run(fmt.Sprintf("user_id_%d", user.ID), func(t *testing.T) {
			err := ValidateUser(user)
			if err != nil {
				t.Errorf("ValidateUser(%+v) returned unexpected error: %v", user, err)
			}
		})
	}
}

func TestValidateUser_NegativeAge(t *testing.T) {
	invalidUsers := []User{
		{ID: 1, Email: "test@example.com", Age: -1},
		{ID: 2, Email: "user@domain.org", Age: -10},
		{ID: 3, Email: "negative@test.com", Age: -100},
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

			if validationErr.Field != "Age" {
				t.Errorf("Expected field 'Age', got '%s'", validationErr.Field)
			}
			if validationErr.Code != 2001 {
				t.Errorf("Expected code 2001, got %d", validationErr.Code)
			}
			if validationErr.Value != user.Age {
				t.Errorf("Expected value %d, got %v", user.Age, validationErr.Value)
			}
		})
	}
}

func TestValidateUser_TooOldAge(t *testing.T) {
	invalidUsers := []User{
		{ID: 1, Email: "test@example.com", Age: 131},
		{ID: 2, Email: "user@domain.org", Age: 150},
		{ID: 3, Email: "ancient@test.com", Age: 200},
	}

	for _, user := range invalidUsers {
		t.Run(fmt.Sprintf("too_old_age_%d", user.Age), func(t *testing.T) {
			err := ValidateUser(user)
			if err == nil {
				t.Errorf("ValidateUser(%+v) expected error but got none", user)
			}

			var validationErr *custom.ValidationError
			if !errors.As(err, &validationErr) {
				t.Errorf("Expected ValidationError, got %T", err)
			}

			if validationErr.Field != "Age" {
				t.Errorf("Expected field 'Age', got '%s'", validationErr.Field)
			}
			if validationErr.Code != 2002 {
				t.Errorf("Expected code 2002, got %d", validationErr.Code)
			}
			if validationErr.Value != user.Age {
				t.Errorf("Expected value %d, got %v", user.Age, validationErr.Value)
			}
		})
	}
}

func TestValidateUser_EmptyEmail(t *testing.T) {
	user := User{ID: 1, Email: "", Age: 25}

	err := ValidateUser(user)
	if err == nil {
		t.Error("ValidateUser with empty email expected error but got none")
	}

	var validationErr *custom.ValidationError
	if !errors.As(err, &validationErr) {
		t.Errorf("Expected ValidationError, got %T", err)
	}

	if validationErr.Field != "Email" {
		t.Errorf("Expected field 'Email', got '%s'", validationErr.Field)
	}
	if validationErr.Code != 2003 {
		t.Errorf("Expected code 2003, got %d", validationErr.Code)
	}
	if validationErr.Value != "" {
		t.Errorf("Expected empty string value, got %v", validationErr.Value)
	}
}

func TestFindUserByEmail_EmptyEmail(t *testing.T) {
	user, err := FindUserByEmail("")

	if user != nil {
		t.Error("FindUserByEmail with empty email should return nil user")
	}
	if err == nil {
		t.Error("FindUserByEmail with empty email should return error")
	}

	// Check that it wraps ErrUserNotFound
	if !errors.Is(err, utils.ErrUserNotFound) {
		t.Error("FindUserByEmail error should wrap ErrUserNotFound")
	}

	// Check error message
	if !strings.Contains(err.Error(), "email cannot be empty") {
		t.Errorf("Error message should mention empty email, got: %s", err.Error())
	}
}

func TestFindUserByEmail_NotFound(t *testing.T) {
	testEmails := []string{
		"test@example.com",
		"user@domain.org",
		"nonexistent@test.com",
	}

	for _, email := range testEmails {
		t.Run(email, func(t *testing.T) {
			user, err := FindUserByEmail(email)

			if user != nil {
				t.Errorf("FindUserByEmail(%s) should return nil user", email)
			}
			if err == nil {
				t.Errorf("FindUserByEmail(%s) should return error", email)
			}

			// Should return ErrUserNotFound
			if !errors.Is(err, utils.ErrUserNotFound) {
				t.Errorf("FindUserByEmail(%s) should return ErrUserNotFound", email)
			}
		})
	}
}

func TestQueryUsers_ReturnsError(t *testing.T) {
	testLimits := []int{1, 10, 100, 1000}

	for _, limit := range testLimits {
		t.Run(fmt.Sprintf("limit_%d", limit), func(t *testing.T) {
			err := QueryUsers(limit)

			if err == nil {
				t.Errorf("QueryUsers(%d) should return error", limit)
			}

			// Should return DatabaseError
			var dbErr *database.DatabaseError
			if !errors.As(err, &dbErr) {
				t.Errorf("QueryUsers(%d) should return DatabaseError, got %T", limit, err)
			}

			// Check DatabaseError fields
			if dbErr.Operation != "SELECT" {
				t.Errorf("Expected operation 'SELECT', got '%s'", dbErr.Operation)
			}
			if dbErr.Table != "users" {
				t.Errorf("Expected table 'users', got '%s'", dbErr.Table)
			}
			if !dbErr.Retryable {
				t.Error("Database error should be retryable")
			}

			// Check that query contains the limit
			expectedQuery := fmt.Sprintf("SELECT * FROM users LIMIT %d", limit)
			if dbErr.Query != expectedQuery {
				t.Errorf("Expected query '%s', got '%s'", expectedQuery, dbErr.Query)
			}

			// Check that it unwraps to the connection timeout error
			unwrapped := dbErr.Unwrap()
			if unwrapped == nil {
				t.Error("DatabaseError should wrap an underlying error")
			}
			if unwrapped.Error() != "connection timeout" {
				t.Errorf("Expected underlying error 'connection timeout', got '%s'", unwrapped.Error())
			}
		})
	}
}
