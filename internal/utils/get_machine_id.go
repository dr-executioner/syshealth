package util

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// GetMachineID returns a unique, stable machine identifier
func GetMachineID() string {
	// Try Linux machine-id
	if id, err := os.ReadFile("/etc/machine-id"); err == nil {
		return strings.TrimSpace(string(id))
	}

	// Try macOS serial number
	if output, err := exec.Command("system_profiler", "SPHardwareDataType").CombinedOutput(); err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "Serial Number") {
				return strings.TrimSpace(strings.Split(line, ":")[1])
			}
		}
	}

	// Try Windows UUID
	if output, err := exec.Command("wmic", "csproduct", "get", "UUID").CombinedOutput(); err == nil {
		lines := strings.Split(string(output), "\n")
		if len(lines) > 1 {
			return strings.TrimSpace(lines[1])
		}
	}

	// Fallback: generate a UUID and persist
	home, _ := os.UserHomeDir()
	idFile := filepath.Join(home, ".system_machine_id")

	if data, err := os.ReadFile(idFile); err == nil {
		return strings.TrimSpace(string(data))
	}

	newID := uuid.NewString()
	_ = os.WriteFile(idFile, []byte(newID), 0644)
	return newID
}
