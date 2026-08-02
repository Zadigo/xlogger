package logic

import (
	"context"
	"fmt"
	"log"

	"github.com/Zadigo/goxlogger/internal/models"
	"github.com/redis/go-redis/v9"
)

func logFileAnalyzer(ctx context.Context, ch chan<- error, serverConfig *models.ServerConfig, redisClient *redis.Client) {
	fileRedis := NewFileRedis(ctx, redisClient)

	logFiles, err := fileRedis.GetLocalLogs(serverConfig.YamlConfig.LogServer.Logs.Folder)

	if err != nil {
		ch <- fmt.Errorf("🔴 Could not get log files: %w", err)
		return
	}

	fileRedis.SaveFiles(logFiles)

	log.Printf("📁 Found %d log files\n", len(logFiles))

	for _, logFile := range logFiles {
		logs, err := fileRedis.ReadFile(logFile.Path, serverConfig)
		if err != nil {
			log.Printf("🔴 Could not read file %s: %s\n", logFile.Path, err)
			continue
		}

		if err = fileRedis.CacheLogs(logFile.Name, logs); err != nil {
			log.Printf("🔴 Could not cache content for file %s: %s\n", logFile.Name, err)
			continue
		}

		logLines := make([]LogLine, 0, len(logs))

		for _, value := range logs {
			logLine := LogLine{RawLine: value}
			result, err := logLine.ParseLine()

			if err != nil {
				select {
				case ch <- err:
				default:
					log.Printf("🔴 Parse error (channel full): %s\n", err)
				}
			}

			logLines = append(logLines, result)
		}

		logsRedis := NewLogsRedis(ctx, redisClient)
		logsRedis.SaveLogs(logLines)
	}
}
