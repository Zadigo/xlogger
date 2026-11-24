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
	"golang.org/x/sync/errgroup"
)

type MainServerApp struct {
	models.MainServerInterface
	errorMessages chan error
	options       *utils.ServerOptions
}

func (m *MainServerApp) GetConfig() *utils.ServerOptions {
	return m.options
}

func (m *MainServerApp) GetRedisDb() *redis.Client {
	return m.options.Redis
}

func (m *MainServerApp) Start(ctx context.Context) error {
	if ctx == nil {
		log.Print("🔴 Context is nil, cannot start HTTP app")
		return nil
	}

	defer m.options.Redis.Close()

	os.Setenv("XLOGGER_STATE", "true")

	log.Print("⚡️ Starting main server...")

	httpServer := httpapp.NewApp(ctx)
	tickerServer := tickerapp.NewLogsApp(ctx)

	g := errgroup.Group{}

	g.Go(func() error {
		return httpServer.Start()
	})

	g.Go(func() error {
		return tickerServer.Start(m)
	})

	return g.Wait()
}

func NewMainServerApp(ctx context.Context, options *utils.ServerOptions) models.MainServerInterface {
	return &MainServerApp{
		errorMessages: make(chan error),
		options:       options,
	}
}
