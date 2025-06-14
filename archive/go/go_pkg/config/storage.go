package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// writeJSONFile marshals data to JSON and writes it to the specified file,
// overwriting any existing file with the same name
func WriteJSONFile(filePath string, data interface{}) error {

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling to JSON: %w", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("error writing file %s: %w", filePath, err)
	}

	return nil
}
