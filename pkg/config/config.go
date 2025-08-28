package config

import (
	"os"
	"path/filepath"
	util "syshealth/internal/utils"
	"time"
)

type Config struct {
	Interval       time.Duration
	StateFilePath  string
	BackendAPI     string
	MachineID      string
	RunOnceForDemo bool
}

func Default() *Config {
	home, _ := os.UserHomeDir()
	stateFile := filepath.Join(home, ".system_health.json")

	return &Config{
		Interval:       15 * time.Minute,
		StateFilePath:  stateFile,
		BackendAPI:     util.GetAPIURL(),
		MachineID:      util.GetMachineID(),
		RunOnceForDemo: false,
	}
}
