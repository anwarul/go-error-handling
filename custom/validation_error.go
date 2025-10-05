package custom

import "fmt"

type ValidationError struct {
	Field   string
	Message string
	Code    int
	Value   interface{}
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation error on field '%s': %s (code: %d, value: %v)", e.Field, e.Message, e.Code, e.Value)
}
