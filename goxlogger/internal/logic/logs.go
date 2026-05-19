package logic

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/Zadigo/goxlogger/internal/models"
	"github.com/go-co-op/gocron"
	"github.com/redis/go-redis/v9"
)

type Logs struct {
	ctx       context.Context
	rootDir   string
	scheduler *gocron.Scheduler
	isStarted atomic.Bool
	debugMode bool
}

func (l *Logs) GetLogs() ([]string, error) {
	var files []string
	fullpath, err := filepath.Abs(l.rootDir)
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(fullpath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func (l *Logs) ReadFile(path string, serverConfig *models.ServerConfig) ([]string, error) {
	file, err := os.Open(path)

	var logs []string = make([]string, 0)
	if err != nil {
		log.Fatal("❌ Could not open file")
		return logs, err
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		logs = append(logs, line)
	}

	return logs, nil
}

func (l *Logs) StartServer(ctx context.Context, serverConfig *models.ServerConfig, redisClient *redis.Client) {
	l.isStarted.Store(true)

	ch := make(chan error, 1)

	go func() {
		_, err := l.scheduler.Cron(serverConfig.YamlConfig.LogServerConfig.Interval).Do(func() {
			logsRedis := LogRedis{ctx: l.ctx, redisClient: redisClient}

			// Get all the log files in the folder
			// path, err := filepath.Abs(serverConfig.Config.LogsFolder.Name)
			logFiles, err := l.GetLogs()
			if err != nil {
				ch <- fmt.Errorf("🔴 Could not get log files: %w", err)
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
				logs, err := l.ReadFile(filePath, serverConfig)
				if err != nil {
					log.Printf("🔴 Could not read file %s: %s\n", filePath, err)
					continue
				}

				// Cache the log files in Redis as sets
				// cmd := redisClient.SAdd(context.Background(), fmt.Sprintf("logs:%s", filepath.Base(filePath)), RedisCachedLogs{
				// 	logs: logs,
				// })

				// if cmd.Err() != nil {
				// 	log.Printf("🔴 Could not cache logs for file %s: %s\n", filePath, cmd.Err())
				// 	continue
				// }

				logLines := make([]LogLine, 0, len(logs))
				logsRedis.SaveLogs(logLines)

				for _, value := range logs {
					logLine := LogLine{RawLine: value}
					result, err := logLine.ParseLine()
					ch <- err
					logLines = append(logLines, result)

					// if err == nil {
					// 	log.Printf("🟢 %s %s %s %d %t\n", result.RemoteAddress, result.Method, result.Path, result.StatusCode, result.IsSuccess)
					// } else {
					// 	log.Printf("🔴 Could not parse line: %s\n", value)
					// }
				}
			}

			if l.debugMode {
				l.StopServer()
				log.Print("🟢 Log server stopped after one run in debug mode")
			}
		})

		ch <- err
	}()

	select {
	case err := <-ch:
		log.Printf("🔴 Log server error: %s\n", err)
	case <-ctx.Done():
		l.StopServer()
		log.Print("🟢 Log server stopped")
	}
}

func (l *Logs) StopServer() bool {
	if l.isStarted.Load() {
		l.scheduler.Stop()
		l.isStarted.Store(false)
		log.Print("🟢 Log server stopped")
		return true
	}
	return false
}
