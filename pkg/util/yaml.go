package util

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// WriteYAML encodes struct to file as YAML.
func WriteYAML(value interface{}, file string) error {
	b, err := yaml.Marshal(value)
	if err != nil {
		return fmt.Errorf("error encoding YAML: %w", err)
	}

	return os.WriteFile(file, b, 0644)
}
