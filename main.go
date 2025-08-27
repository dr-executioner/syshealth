package main

import (
	"fmt"
	"syshealth/pkg/config"
	"syshealth/pkg/scheduler"
	"time"
)

func main() {
	fmt.Println("monitor-agent starting...")
	cfg := config.Default()

	scheduler.Start(cfg)

	for {
		time.Sleep(24 * time.Hour)
	}
}
