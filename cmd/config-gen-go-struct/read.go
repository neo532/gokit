package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func readConfig(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	switch filepath.Ext(path) {
	case ".json":
		return decodeJSON(data)
	case ".ini":
		return decodeINI(data)
	default:
		return decodeYAML(data)
	}
}

func decodeYAML(data []byte) (map[string]any, error) {
	return yamlDecode(data)
}

func decodeJSON(data []byte) (map[string]any, error) {
	return jsonDecode(data)
}

func decodeINI(data []byte) (map[string]any, error) {
	return iniDecode(data)
}

// formatError wraps parse errors with the format name.
func formatError(format string, err error) error {
	return fmt.Errorf("%s: %w", format, err)
}
