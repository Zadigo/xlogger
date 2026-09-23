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

	// logFiles, err := fileRedis.CollectFilesInFolder(serverConfig.LogServer.Logs.Folder)
	rootDir := ctx.Value("rootDir").(string)
	if rootDir == "" {
		ch <- fmt.Errorf("🔴 Root directory is not set in context")
		return
	}

	fileCollector := FileCollector{}
	logFiles, err := fileCollector.CollectFilesInFolder(rootDir, serverConfig.LogServer.Logs.Folder)
	if err != nil {
		ch <- fmt.Errorf("🔴 Could not get log files: %w", err)
		return
	}

	// Check the number of log files locally and those registered in Redis,
	// if the number defers, add the missing files to Redis

	fileRedis.SaveFiles(logFiles)

	log.Printf("📁 Found %d log files", len(logFiles))

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
	}
}
