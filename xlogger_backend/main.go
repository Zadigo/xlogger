package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/Zadigo/xlogger_backend/internal"
	"github.com/Zadigo/xlogger_backend/internal/backend"
	"github.com/Zadigo/xlogger_backend/internal/handlers"
	"github.com/go-co-op/gocron"
)

type RedisCachedLogs struct {
	logs []string
}

func main() {
	serverConfig := backend.GetServerConfig()
	scheduler := gocron.NewScheduler(time.UTC)

	// Redis
	redisClient := backend.NewRedisClient(serverConfig)

	// Goroutune to run the scheduler
	// in the background
	go func() {
		_, err := scheduler.Cron(serverConfig.Config.Analysis).Do(func() {
			// Get all the log files in the folder
			path, err := filepath.Abs(serverConfig.Config.LogsFolder.Name)
			logFiles, err := backend.FilePathWalkDir(path)
			if err != nil {
				log.Fatal("🔴 Could not get log files")
			}

			fmt.Printf("📁 Found %d log files\n", len(logFiles))

			for i, logFile := range logFiles {
				extension := filepath.Ext(logFile)
				if extension != ".log" {
					logFiles = append(logFiles[:i], logFiles[i+1:]...)
					log.Printf("⚠️ Skipping file %s with unsupported extension %s\n", logFile, extension)
					continue
				}
			}

			for _, filePath := range logFiles {
				logs, err := internal.ReadFile(filePath, serverConfig)
				if err != nil {
					log.Printf("🔴 Could not read file %s: %s\n", filePath, err)
					continue
				}

				// Cache the log files in Redis as sets
				cmd := redisClient.SAdd(context.Background(), fmt.Sprintf("logs:%s", filepath.Base(filePath)), RedisCachedLogs{
					logs: logs,
				})

				if cmd.Err() != nil {
					log.Printf("🔴 Could not cache logs for file %s: %s\n", filePath, cmd.Err())
					continue
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
			}
		})

		if err != nil {
			log.Fatal("Could not schedule job")
		}

		scheduler.StartAsync()
	}()

	// Keep the main function running
	http.HandleFunc("/ws/live", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.LiveWsHandler(w, r, redisClient)
	}))
	log.Println("🚀 Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
