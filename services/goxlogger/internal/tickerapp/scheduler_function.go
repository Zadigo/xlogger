package tickerapp

import (
	"context"
	"fmt"
	"log"

	"github.com/Zadigo/goxlogger/internal/utils"
	"github.com/redis/go-redis/v9"
)

func logFileAnalyzer(ctx context.Context, ch chan<- error, serverConfig *utils.ServerConfig, redisClient *redis.Client) {
	fileRedis := NewFileRedis(ctx, redisClient)

	logFiles, err := fileRedis.CollectFilesInFolder(serverConfig.LogServer.Logs.Folder)

	if err != nil {
		ch <- fmt.Errorf("🔴 Could not get log files: %w", err)
		return
	}

	fileRedis.SaveFiles(logFiles)

	log.Printf("📁 Found %d log files\n", len(logFiles))

	for _, logFile := range logFiles {
		strLogs, err := fileRedis.ReadFile(logFile.Path, serverConfig)
		if err != nil {
			log.Printf("🔴 Could not read file %s: %s\n", logFile.Path, err)
			continue
		}

		if err = fileRedis.CacheLogs(logFile.Name, strLogs); err != nil {
			log.Printf("🔴 Could not cache content for file %s: %s\n", logFile.Name, err)
			continue
		}

		_, err = NewLogsRedis(ctx, redisClient).SaveTransform(strLogs)
		if err != nil {
			ch <- fmt.Errorf("🔴 Could not transform logs for file %s: %w", logFile.Name, err)
			continue
		}

		// logLines := make([]LogLine, 0, len(strLogs))

		// for _, value := range strLogs {
		// 	logLine := LogLine{RawLine: value}
		// 	result, err := logLine.ParseLine()

		// 	if err != nil {
		// 		select {
		// 		case ch <- err:
		// 		default:
		// 			log.Printf("🔴 Parse error (channel full): %s\n", err)
		// 		}
		// 	}

		// 	logLines = append(logLines, result)
		// }

		// logsRedis := NewLogsRedis(ctx, redisClient)
		// logsRedis.SaveLogs(logLines)
	}
}
