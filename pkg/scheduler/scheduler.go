package scheduler

import (
	"encoding/json"
	"fmt"
	"time"

	"syshealth/pkg/checks"
	"syshealth/pkg/config"
	"syshealth/pkg/reporter"
	"syshealth/pkg/storage"
	"syshealth/pkg/structs"
)

func StartScheduler(cfg *config.Config) {
	fmt.Printf("scheduler: starting loop, interval: %s\n", cfg.Interval)

	for {
		results := []structs.CheckResult{
			checks.DiskEncryption(),
			checks.OSUpdates(),
			checks.SleepTimeout(),
			checks.Antivirus(),
		}

		current := map[string]interface{}{
			"checks": results,
		}

		// Load previous state from cfg.StateFilePath
		prev, _ := storage.LoadState(cfg.StateFilePath)

		// Compare states
		if !stateEqual(prev, current) {
			fmt.Println("changes detected, sending report")
			err := reporter.Send(cfg, results)
			if err == nil {
				_ = storage.SaveState(cfg.StateFilePath, current)
			}
		} else {
			fmt.Println("no changes, skipping report")
		}

		if cfg.RunOnceForDemo {
			break
		}

		time.Sleep(cfg.Interval)
	}
}

func stateEqual(a, b map[string]interface{}) bool {
	aj, _ := json.Marshal(a)
	bj, _ := json.Marshal(b)
	return string(aj) == string(bj)
}
