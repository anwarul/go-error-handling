package utils

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrDuplicateEmail  = errors.New("email already exists")
	ErrInvalidPassword = errors.New("invalid password")
	ErrUnauthorized    = errors.New("unauthorized access")
	ErrDatabaseTimeout = errors.New("database operation timed out")
)
