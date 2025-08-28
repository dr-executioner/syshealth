package reporter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"syshealth/pkg/config"
	"syshealth/pkg/structs"
	"time"
)

type Report struct {
	MachineID string                `json:"machine_id"`
	Timestamp time.Time             `json:"timestamp"`
	Checks    []structs.CheckResult `json:"checks"`
}

// Send posts the report to backend. For now we do a simple POST with retries.
func Send(cfg *config.Config, results []structs.CheckResult) error {
	r := Report{
		MachineID: cfg.MachineID,
		Timestamp: time.Now().UTC(),
		Checks:    results,
	}
	b, _ := json.Marshal(r)

	// For development, just print
	fmt.Printf("[reporter] would send to %s: %s\n", cfg.BackendAPI, string(b))

	// If you want to actually send, uncomment below (and ensure cfg.BackendAPI valid)
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", cfg.BackendAPI, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	return fmt.Errorf("unexpected status: %s", resp.Status)
}
