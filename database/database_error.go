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

func Unwramp(err error) error {
	type unwrapper interface {
		Unwrap() error
	}
	if ue, ok := err.(unwrapper); ok {
		return ue.Unwrap()
	}
	return nil
}

func (e *DatabaseError) Unwrap() error {
	return e.Err
}
