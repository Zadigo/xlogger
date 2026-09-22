package mainapp

import (
	"context"
	"log"
	"os"

	"github.com/Zadigo/goxlogger/internal/httpapp"
	"github.com/Zadigo/goxlogger/internal/models"
	"github.com/Zadigo/goxlogger/internal/tickerapp"
	"github.com/Zadigo/goxlogger/internal/utils"
	"github.com/redis/go-redis/v9"
)

type MainServerApp struct {
	models.MainServerInterface
	ctx context.Context
	errorMessages chan error
	config *utils.ServerConfig
	redisDb *redis.Client
}

func (m *MainServerApp) GetContext() context.Context {
	return m.ctx
}

func (m *MainServerApp) GetConfig() *utils.ServerConfig {
	return m.config
}

func (m *MainServerApp) GetRedisDb() *redis.Client {
	return m.redisDb
}

func (m *MainServerApp) Start() {
	if m.ctx == nil {
		log.Print("🔴 Context is nil, cannot start HTTP app")
		return
	}

	os.Setenv("XLOGGER_STATE", "true")

	log.Print("⚡️ Starting main server...")
	httpServer := httpapp.NewApp(m.ctx)

	go func() {
		m.errorMessages <- httpServer.Start()
	}()

	tickerServer := tickerapp.NewLogsApp(m.ctx)

	go func() {
		tickerServer.Start(m)
	}()

	<- m.ctx.Done()

	log.Print("🔴 Stopping X-Logger...")

	os.Unsetenv("XLOGGER_STATE")
	m.redisDb.Close()
}

func NewMainServerApp(ctx context.Context, config *utils.ServerConfig) models.MainServerInterface {
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
