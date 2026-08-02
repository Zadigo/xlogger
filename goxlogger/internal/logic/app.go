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
			logFileAnalyzer(l.ctx, ch, serverConfig, redisClient)
		})

		if err != nil {
			ch <- fmt.Errorf("🔴 Could not schedule log server: %w", err)
		}

		l.scheduler.StartBlocking()
	}()

	for {
		select {
		case err := <-ch:
			log.Printf("🔴 Log server error: %s\n", err)
		case <-l.ctx.Done():
			l.scheduler.Stop()
			l.isStarted.Store(false)

			close(ch)

			log.Print("🟢 Log server stopped")
		}
	}
}

func NewLogsService(ctx context.Context) *Logs {
	rootDir := ctx.Value("rootDir").(string)
	debugMode := ctx.Value("debugMode").(bool)

	return &Logs{
		ctx:       ctx,
		rootDir:   rootDir,
		scheduler: gocron.NewScheduler(time.UTC),
		debugMode: debugMode,
		isStarted: atomic.Bool{},
	}
}
