package tickerapp

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/Zadigo/goxlogger/internal/models"
	"github.com/go-co-op/gocron"
)

type LogsApp struct {
	ctx       context.Context
	rootDir   string
	scheduler *gocron.Scheduler
	isStarted atomic.Bool
	debugMode bool
}

func (l *LogsApp) Start(server models.MainServerInterface) {
	l.isStarted.Store(true)
	
	serverConfig := server.GetConfig()
	log.Printf("🟢 Starting log server with interval %s\n", serverConfig.LogServer.Interval)

	ch := make(chan error, 1)

	go func() {
		_, err := l.scheduler.Cron(serverConfig.LogServer.Interval).Do(func() {
			logFileAnalyzer(l.ctx, ch, serverConfig, server.GetRedisDb())
		})

		if err != nil {
			ch <- fmt.Errorf("🔴 Could not schedule log server: %w", err)
		}

		l.scheduler.StartBlocking()
	}()

	for {
		select {
		case err := <-ch:
			log.Printf("🔴 Log server error: %s", err)
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
