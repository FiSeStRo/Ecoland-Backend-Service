package config

import (
	"bufio"
	"os"
	"strings"
)

// LoadEnv loads environment variables from a file
func LoadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if len(line) == 0 || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}

		// Parse key=value pair
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // Skip malformed lines
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Set environment variable
		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}

	return scanner.Err()
}
