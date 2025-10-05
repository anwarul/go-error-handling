package formatted

import (
	"fmt"
)

func ValidateAge(age int) error {
	if age < 0 {
		return fmt.Errorf("invalid age: %d. Age cannot be negative", age)
	}
	if age > 130 {
		return fmt.Errorf("invalid age: %d. Age cannot be greater than 130", age)
	}
	return nil
}
