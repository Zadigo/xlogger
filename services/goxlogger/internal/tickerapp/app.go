package tickerapp

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/Zadigo/goxlogger/internal/utils"
	"github.com/go-co-op/gocron"
	"github.com/redis/go-redis/v9"
)

type LogsApp struct {
	ctx       context.Context
	rootDir   string
	scheduler *gocron.Scheduler
	isStarted atomic.Bool
	debugMode bool
}

func (l *LogsApp) Start(serverConfig *utils.ServerConfig, redisClient *redis.Client) {
	l.isStarted.Store(true)
	
	log.Printf("🟢 Starting log server with interval %s\n", serverConfig.LogServer.Interval)

	ch := make(chan error, 1)

	go func() {
		_, err := l.scheduler.Cron(serverConfig.LogServer.Interval).Do(func() {
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

func NewLogsApp(ctx context.Context) *LogsApp {
	rootDir := ctx.Value("rootDir").(string)
	debugMode := ctx.Value("debugMode").(bool)

	return &LogsApp{
		ctx:       ctx,
		rootDir:   rootDir,
		scheduler: gocron.NewScheduler(time.UTC),
		debugMode: debugMode,
		isStarted: atomic.Bool{},
	}
}
