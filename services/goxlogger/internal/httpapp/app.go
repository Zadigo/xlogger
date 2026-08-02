package httpapp

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Zadigo/goxlogger/internal/backend"
	"github.com/Zadigo/goxlogger/internal/models"
	"github.com/go-chi/chi"
	"github.com/redis/go-redis/v9"
)

type App struct {
	ctx         context.Context
	parentCtx   context.Context
	router      *chi.Mux
	chanErr     chan error
	redisClient *redis.Client
}

func (a *App) Start() error {
	var cancel context.CancelFunc
	a.ctx, cancel = context.WithCancel(a.parentCtx)
	defer cancel()

	if a.router == nil {
		panic("Router is not initialized. Please call SetupRouter() before starting the server.")
	}

	if a.redisClient == nil {
		panic("Redis client is not initialized. Please call NewApp() to initialize the Redis client.")
	}

	if a.ctx == nil {
		panic("Context is not initialized. Please call NewApp() to initialize the context.")
	}

	port, err := strconv.ParseUint(os.Getenv("GO_PORT"), 10, 16)
	if err != nil {
		port = 9000
	}

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: a.router,
	}

	go func() {
		log.Printf("⚡️ Starting server on port %d...", port)
		a.chanErr <- server.ListenAndServe()
	}()

	select {
	case err := <-a.chanErr:
		return fmt.Errorf("🔴 %s HTTP server error: %v", os.Getenv("SERVICE_NAME"), err)
	case <-a.ctx.Done():
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		a.redisClient.Close()
		shutdownErr := server.Shutdown(timeoutCtx)

		// close(a.chanErr)

		return errors.Join(fmt.Errorf("🔴 Shutting down %s HTTP server...", os.Getenv("SERVICE_NAME")), shutdownErr, a.ctx.Err())
	}
}

func (a *App) GetRedisClient() *redis.Client {
	return a.redisClient
}

func (a *App) GetAppContext() context.Context {
	return a.ctx
}

func NewApp(ctx context.Context) models.AppInterface {
	redisClient := backend.NewRedisBackend(ctx)

	app := &App{
		ctx:         nil,
		parentCtx:   ctx,
		router:      nil,
		chanErr:     make(chan error),
		redisClient: redisClient,
	}

	app.loadRoutes()

	return app
}
