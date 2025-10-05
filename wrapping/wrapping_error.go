package wrapping

import (
	"fmt"
	"os"
)

func ProcessUserData(userID int) error {
	err := loadUserConfig(userID)
	if err != nil {
		return fmt.Errorf("failed to process user %d: %w", userID, err)
	}
	return nil
}

func loadUserConfig(userID int) error {
	filename := fmt.Sprintf("user_%d.json", userID)
	err := readConfigFile(filename)
	if err != nil {
		return fmt.Errorf("failed to load config for user %d: %w", userID, err)
	}
	return nil
}

func readConfigFile(filename string) error {
	_, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", filename, err)
	}
	return nil
}
