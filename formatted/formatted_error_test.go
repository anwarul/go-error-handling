package formatted

import (
	"fmt"
	"strings"
	"testing"
)

func TestValidateAge_Success(t *testing.T) {
	validAges := []int{0, 1, 18, 25, 65, 100, 130}

	for _, age := range validAges {
		t.Run(fmt.Sprintf("age_%d", age), func(t *testing.T) {
			err := ValidateAge(age)
			if err != nil {
				t.Errorf("ValidateAge(%d) returned unexpected error: %v", age, err)
			}
		})
	}
}

func TestValidateAge_NegativeAge(t *testing.T) {
	invalidAges := []int{-1, -5, -100}

	for _, age := range invalidAges {
		t.Run(fmt.Sprintf("negative_age_%d", age), func(t *testing.T) {
			err := ValidateAge(age)
			if err == nil {
				t.Errorf("ValidateAge(%d) expected error but got none", age)
			}

			expectedSubstring := "Age cannot be negative"
			if !strings.Contains(err.Error(), expectedSubstring) {
				t.Errorf("ValidateAge(%d) error = %v; want error containing '%s'", age, err.Error(), expectedSubstring)
			}

			// Check that the age is included in the error message
			ageStr := fmt.Sprintf("%d", age)
			if !strings.Contains(err.Error(), ageStr) {
				t.Errorf("ValidateAge(%d) error message should contain the age value", age)
			}
		})
	}
}

func TestValidateAge_TooOldAge(t *testing.T) {
	invalidAges := []int{131, 150, 200, 1000}

	for _, age := range invalidAges {
		t.Run(fmt.Sprintf("too_old_age_%d", age), func(t *testing.T) {
			err := ValidateAge(age)
			if err == nil {
				t.Errorf("ValidateAge(%d) expected error but got none", age)
			}

			expectedSubstring := "Age cannot be greater than 130"
			if !strings.Contains(err.Error(), expectedSubstring) {
				t.Errorf("ValidateAge(%d) error = %v; want error containing '%s'", age, err.Error(), expectedSubstring)
			}

			// Check that the age is included in the error message
			ageStr := fmt.Sprintf("%d", age)
			if !strings.Contains(err.Error(), ageStr) {
				t.Errorf("ValidateAge(%d) error message should contain the age value", age)
			}
		})
	}
}

func TestValidateAge_BoundaryValues(t *testing.T) {
	tests := []struct {
		age         int
		shouldError bool
		description string
	}{
		{-1, true, "just below minimum"},
		{0, false, "minimum valid age"},
		{130, false, "maximum valid age"},
		{131, true, "just above maximum"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := ValidateAge(tt.age)
			if tt.shouldError && err == nil {
				t.Errorf("ValidateAge(%d) expected error but got none", tt.age)
			}
			if !tt.shouldError && err != nil {
				t.Errorf("ValidateAge(%d) expected no error but got: %v", tt.age, err)
			}
		})
	}
}

func TestValidateAge_ErrorMessages(t *testing.T) {
	tests := []struct {
		age              int
		expectedContains []string
	}{
		{
			age:              -10,
			expectedContains: []string{"invalid age", "-10", "Age cannot be negative"},
		},
		{
			age:              150,
			expectedContains: []string{"invalid age", "150", "Age cannot be greater than 130"},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("age_%d_message", tt.age), func(t *testing.T) {
			err := ValidateAge(tt.age)
			if err == nil {
				t.Errorf("ValidateAge(%d) expected error but got none", tt.age)
				return
			}

			errMsg := err.Error()
			for _, substring := range tt.expectedContains {
				if !strings.Contains(errMsg, substring) {
					t.Errorf("ValidateAge(%d) error message '%s' should contain '%s'", tt.age, errMsg, substring)
				}
			}
		})
	}
}
