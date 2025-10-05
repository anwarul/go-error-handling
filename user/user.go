package user

import (
	"errors"
	"fmt"
	"go-error-handling/custom"
	"go-error-handling/database"
	"go-error-handling/utils"
	"time"
)

type ValidationError = custom.ValidationError

type User struct {
	ID    int
	Email string
	Age   int
}

func ValidateUser(user User) error {
	if user.Age < 0 {
		return &ValidationError{
			Field:   "Age",
			Message: "Age cannot be negative",
			Code:    2001,
			Value:   user.Age,
		}
	}
	if user.Age > 130 {
		return &ValidationError{
			Field:   "Age",
			Message: "Age cannot be greater than 130",
			Code:    2002,
			Value:   user.Age,
		}
	}
	if user.Email == "" {
		return &ValidationError{
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
