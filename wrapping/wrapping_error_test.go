package wrapping

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestProcessUserData_FileNotFound(t *testing.T) {
	// Test with a user ID that will generate a non-existent file
	userID := 999

	err := ProcessUserData(userID)

	if err == nil {
		t.Errorf("ProcessUserData(%d) expected error but got none", userID)
	}

	// Check that the error wraps os.ErrNotExist
	if !errors.Is(err, os.ErrNotExist) {
		t.Error("ProcessUserData error should wrap os.ErrNotExist")
	}

	// Check the error chain contains information about the user ID
	errMsg := err.Error()
	expectedComponents := []string{
		"failed to process user",
		"999",
		"failed to load config for user",
		"failed to read config file",
		"user_999.json",
	}

	for _, component := range expectedComponents {
		if !strings.Contains(errMsg, component) {
			t.Errorf("Error message should contain '%s', got: %s", component, errMsg)
		}
	}
}

func TestProcessUserData_ErrorChain(t *testing.T) {
	userIDs := []int{1, 42, 123, 500}

	for _, userID := range userIDs {
		t.Run(fmt.Sprintf("user_%d", userID), func(t *testing.T) {
			err := ProcessUserData(userID)

			if err == nil {
				t.Errorf("ProcessUserData(%d) expected error but got none", userID)
			}

			// Test the error chain depth
			var currentErr error = err
			depth := 0

			for currentErr != nil {
				depth++
				if depth > 10 { // Prevent infinite loop
					t.Error("Error chain too deep, possible infinite loop")
					break
				}

				// Try to unwrap
				if unwrapper, ok := currentErr.(interface{ Unwrap() error }); ok {
					currentErr = unwrapper.Unwrap()
				} else {
					break
				}
			}

			// Should have multiple levels in the error chain
			if depth < 3 {
				t.Errorf("Expected at least 3 levels in error chain, got %d", depth)
			}
		})
	}
}

func TestLoadUserConfig_ErrorPropagation(t *testing.T) {
	userID := 123
	err := loadUserConfig(userID)

	if err == nil {
		t.Error("loadUserConfig should return error for non-existent file")
	}

	// Check that it wraps the file reading error
	if !errors.Is(err, os.ErrNotExist) {
		t.Error("loadUserConfig error should wrap os.ErrNotExist")
	}

	// Check error message format
	errMsg := err.Error()
	expectedComponents := []string{
		"failed to load config for user",
		"123",
		"failed to read config file",
		"user_123.json",
	}

	for _, component := range expectedComponents {
		if !strings.Contains(errMsg, component) {
			t.Errorf("Error message should contain '%s', got: %s", component, errMsg)
		}
	}
}

func TestReadConfigFile_Success(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test_config_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Write some test data
	testData := `{"name": "test", "value": 123}`
	if _, err := tmpFile.WriteString(testData); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}

	// Test reading the file
	err = readConfigFile(tmpFile.Name())
	if err != nil {
		t.Errorf("readConfigFile with valid file should not return error, got: %v", err)
	}
}

func TestReadConfigFile_FileNotFound(t *testing.T) {
	nonExistentFile := "definitely_does_not_exist_12345.json"

	err := readConfigFile(nonExistentFile)

	if err == nil {
		t.Error("readConfigFile with non-existent file should return error")
	}

	// Check that it wraps os.ErrNotExist
	if !errors.Is(err, os.ErrNotExist) {
		t.Error("readConfigFile error should wrap os.ErrNotExist")
	}

	// Check error message format
	errMsg := err.Error()
	expectedComponents := []string{
		"failed to read config file",
		nonExistentFile,
	}

	for _, component := range expectedComponents {
		if !strings.Contains(errMsg, component) {
			t.Errorf("Error message should contain '%s', got: %s", component, errMsg)
		}
	}
}

func TestErrorChain_UnwrapBehavior(t *testing.T) {
	err := ProcessUserData(999)

	// Test that we can traverse the error chain
	var levels []string
	currentErr := err

	for currentErr != nil {
		levels = append(levels, currentErr.Error())

		// Try to unwrap
		if unwrapper, ok := currentErr.(interface{ Unwrap() error }); ok {
			currentErr = unwrapper.Unwrap()
		} else {
			break
		}
	}

	// Should have multiple levels
	if len(levels) < 3 {
		t.Errorf("Expected at least 3 error levels, got %d: %v", len(levels), levels)
	}

	// First level should mention processing user
	if !strings.Contains(levels[0], "failed to process user") {
		t.Errorf("First error level should mention processing user, got: %s", levels[0])
	}

	// Last level should be the original file system error
	lastLevel := levels[len(levels)-1]
	if !strings.Contains(lastLevel, "no such file or directory") &&
		!strings.Contains(lastLevel, "cannot find the file") &&
		!strings.Contains(lastLevel, "system cannot find the file") {
		t.Errorf("Last error level should be file system error, got: %s", lastLevel)
	}
}
