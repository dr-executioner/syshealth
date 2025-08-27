package config

import "time"

type Config struct {
	Interval       time.Duration
	StateFilePath  string
	BackendAPI     string
	MachineID      string
	RunOnceForDemo bool
}

func Default() *Config {
	return &Config{
		Interval:       1 * time.Minute,
		StateFilePath:  "../logs/system_health.json",
		BackendAPI:     "http://127.0.0.1:8080/api/report",
		MachineID:      "dev-machine-001",
		RunOnceForDemo: false,
	}
}
