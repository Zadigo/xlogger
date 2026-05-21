package logic

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

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

func (l *Logs) StartServer(serverConfig *models.ServerConfig, redisClient *redis.Client) {
	l.isStarted.Store(true)
	log.Printf("🟢 Starting log server with interval %s\n", serverConfig.YamlConfig.LogServer.Interval)

	ch := make(chan error, 1)

	go func() {
		_, err := l.scheduler.Cron(serverConfig.YamlConfig.LogServer.Interval).Do(func() {
			fileRedis := NewFileRedis(l.ctx, l.rootDir, redisClient)

			// Get all the log files in the folder
			// path, err := filepath.Abs(serverConfig.Config.LogsFolder.Name)
			logFiles, err := fileRedis.GetLocalLogs()
			if err != nil {
				ch <- fmt.Errorf("🔴 Could not get log files: %w", err)
			}

			fmt.Printf("📁 Found %d log files\n", len(logFiles))

			for _, logFile := range logFiles {
				logs, err := fileRedis.ReadFile(logFile.Path, serverConfig)
				if err != nil {
					log.Printf("🔴 Could not read file %s: %s\n", logFile.Path, err)
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

				logsRedis := LogRedis{ctx: l.ctx, redisClient: redisClient}
				logsRedis.SaveLogs(logLines)
			}

			if l.debugMode {
				l.StopServer()
				log.Print("🟢 Log server stopped after one run in debug mode")
			}
		})

		if err != nil {
			ch <- fmt.Errorf("🔴 Could not schedule log server: %w", err)
		}

		l.scheduler.StartAsync()
	}()

	for {
		select {
		case err := <-ch:
			log.Printf("🔴 Log server error: %s\n", err)
			return
		case <-l.ctx.Done():
			l.StopServer()
			return
		}
	}

	// select {
	// case err := <-ch:
	// 	log.Printf("🔴 Log server error: %s\n", err)
	// case <-l.ctx.Done():
	// 	l.StopServer()
	// 	log.Print("🟢 Log server stopped")
	// }
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

func NewLogsService(ctx context.Context, rootDir string, debugMode bool) *Logs {
	return &Logs{
		ctx:       ctx,
		rootDir:   rootDir,
		scheduler: gocron.NewScheduler(time.UTC),
		debugMode: debugMode,
	}
}
