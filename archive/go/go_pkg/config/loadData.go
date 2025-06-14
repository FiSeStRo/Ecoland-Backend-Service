package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func LoadJsonDataFromFileStorage[T any](name string) (T, error) {

	content, err := os.ReadFile(strings.Join([]string{"../../file_storage/", name}, ""))
	var zero T
	if err != nil {
		return zero, fmt.Errorf("could not load file: %w", err)
	}

	var data T
	if err := json.Unmarshal(content, &data); err != nil {
		return zero, fmt.Errorf("could not read data: %w", err)
	}

	return data, nil
}
