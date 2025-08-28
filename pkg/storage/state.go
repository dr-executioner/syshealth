package storage

import (
	"encoding/json"
	"os"
)

// SaveState writes the current state to the given path
func SaveState(path string, data map[string]interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0644)
}

// LoadState reads the state from the given path
func LoadState(path string) (map[string]interface{}, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var state map[string]interface{}
	if err := json.Unmarshal(bytes, &state); err != nil {
		return nil, err
	}
	return state, nil
}

// EnsureStateFile makes sure the state file exists, creating it if missing
func EnsureStateFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.WriteFile(path, []byte("{}"), 0644)
	}
	return nil
}
