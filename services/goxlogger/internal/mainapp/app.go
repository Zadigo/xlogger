package mainapp

import (
	"context"
	"log"

	"github.com/Zadigo/goxlogger/internal/httpapp"
	"github.com/Zadigo/goxlogger/internal/tickerapp"
	"github.com/Zadigo/goxlogger/internal/utils"
	"github.com/redis/go-redis/v9"
)

type MainServerApp struct {
	ctx context.Context
	errorMessages chan error
	config *utils.ServerConfig
	redisDb *redis.Client
}

func (m *MainServerApp) Start() {
	if m.ctx == nil {
		log.Print("🔴 Context is nil, cannot start HTTP app")
		return
	}

	log.Print("⚡️ Starting main server...")
	httpServer := httpapp.NewApp(m.ctx)

	go func() {
		m.errorMessages <- httpServer.Start()
	}()

	tickerServer := tickerapp.NewLogsApp(m.ctx)

	go func() {
		tickerServer.Start(m.config, m.redisDb)
	}()

	<- m.ctx.Done()

	log.Print("🔴 Stopping X-Logger...")
	m.redisDb.Close()
}

func NewMainServerApp(ctx context.Context, config *utils.ServerConfig) *MainServerApp {
	redisDb := redis.NewClient(&redis.Options{
		Addr: config.Redis.Addr,
	})
	
	return &MainServerApp{
		ctx: ctx,
		errorMessages: make(chan error),
		config: config,
		redisDb: redisDb,
	}
}
