package main

import (
	"log"
	"os"
	"path/filepath"
	"syshealth/pkg/config"
	"syshealth/pkg/scheduler"
	"syshealth/pkg/storage"
	"time"
)

func main() {
	cfg := config.Default()

	logFile := filepath.Join(os.TempDir(), "syshealth-agent.log")
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println("SysHealth Agent starting...")

	if err := storage.EnsureStateFile(cfg.StateFilePath); err != nil {
		log.Fatal(err)
	}
	scheduler.StartScheduler(cfg)

	for {
		time.Sleep(24 * time.Hour)
	}
}
