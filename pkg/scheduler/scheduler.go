package scheduler

import (
	"fmt"
	"syshealth/pkg/checks"
	"syshealth/pkg/config"
	"syshealth/pkg/reporter"
	"syshealth/pkg/storage"
	"time"
)

func Start(cfg *config.Config) {
	fmt.Println("scheduler: starting loop, interval:", cfg.Interval)

	state, err := storage.LoadState(cfg.StateFilePath)
	if err != nil {
		fmt.Println("error loading state:", err)
		state = make(storage.State)
	}

	runOnce := func() {
		results := []checks.CheckResult{}
		results = append(results, checks.DiskEncryption())
		results = append(results, checks.OSUpdates())
		results = append(results, checks.Antivirus())
		results = append(results, checks.SleepTimeout())

		changed := false
		newState := make(storage.State)
		for _, r := range results {
			newState[r.Name] = r
			old, ok := state[r.Name]
			if !ok || old.OK != r.OK || old.Detail != r.Detail {
				changed = true
			}
		}

		if changed {
			fmt.Println("changes detected, sending report")
			if err := reporter.Send(cfg, results); err != nil {
				fmt.Println("report send error:", err)
			} else {
				// save new state
				if err := storage.SaveState(cfg.StateFilePath, newState); err != nil {
					fmt.Println("error saving state:", err)
				}
				state = newState
			}
		} else {
			fmt.Println("no changes detected")
		}
	}

	// first run immediately
	runOnce()

	if cfg.RunOnceForDemo {
		fmt.Println("run once mode, exiting")
		return
	}

	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			runOnce()
		}
	}
}
