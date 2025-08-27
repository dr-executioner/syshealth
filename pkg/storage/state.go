package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"syshealth/pkg/checks"
)

// State maps check name -> CheckResult
type State map[string]checks.CheckResult

// LoadState loads state from a path. If file missing, returns empty state
func LoadState(path string) (State, error) {
	expanded := filepath.Clean(path)
	if _, err := os.Stat(expanded); errors.Is(err, os.ErrNotExist) {
		return State{}, nil
	}
	b, err := os.ReadFile(expanded)
	if err != nil {
		return nil, err
	}
	var s State
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// SaveState writes the state to the given path (overwrites)
func SaveState(path string, s State) error {
	expanded := filepath.Clean(path)
	b, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(expanded, b, 0644)
}
