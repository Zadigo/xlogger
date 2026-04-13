package main

import (
	"log"
	"time"

	"github.com/Zadigo/xlogger_backend/internal"
	"github.com/Zadigo/xlogger_backend/internal/backend"
	"github.com/go-co-op/gocron"
)

func main() {
	serverConfig := backend.GetServerConfig()
	scheduler := gocron.NewScheduler(time.UTC)

	_, err := scheduler.Cron(serverConfig.Config.Analysis).Do(func() {
		logs, err := internal.ReadFile("example2.log", serverConfig)
		if err != nil {
			log.Fatal("🔴 Could not read file")
		}

		for _, value := range logs {
			logLine := backend.LogLine{RawLine: value}
			result, err := logLine.ParseLine()
			if err == nil {
				log.Printf("🟢 %s %s %s %d %t\n", result.RemoteAddress, result.Method, result.Path, result.StatusCode, result.IsSuccess)
			} else {
				log.Printf("🔴 Could not parse line: %s\n", value)
			}
		}
	})

	if err != nil {
		log.Fatal("Could not schedule job")
	}
}
