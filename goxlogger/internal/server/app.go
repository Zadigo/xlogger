package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Zadigo/goxlogger/internal/logic"
	"github.com/Zadigo/goxlogger/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

type App struct {
	ctx         context.Context
	redisClient *redis.Client
	config      *models.ServerConfig
	router      *chi.Mux
}

func (a *App) Start() error {
	// Redis client
	err := a.redisClient.Ping(a.ctx).Err()
	if err != nil {
		return fmt.Errorf("🔴 Failed to load Redis: %w", err)
	}

	defer func() {
		a.redisClient.Close()
		log.Print("🔴 Redis client closed")
	}()

	port := os.Getenv("XLOGGER_PORT")
	if port == "" {
		port = "9000"
	}

	// HTTP server
	server := http.Server{
		Addr:    ":" + port,
		Handler: a.router,
	}

	ch := make(chan error, 1)

	go func() {
		log.Printf("🟢 Server is listening on port %s", port)
		ch <- server.ListenAndServe()
	}()

	// Log server
	logServer := logic.NewLogsService(a.ctx, a.config.RootDir, false)
	logServer.StartServer(a.config, a.redisClient)

	select {
	case err := <-ch:
		return err

	case <-a.ctx.Done():
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return server.Shutdown(timeoutCtx)
	}
}

func NewApp(ctx context.Context, config *models.ServerConfig) *App {
	redisAddr := os.Getenv("REDIS_ADDR")

	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	app := &App{
		ctx:    ctx,
		config: config,
		redisClient: redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Username: "",
			Password: "",
			DB:       0,
		}),
	}
	app.loadRoutes()
	return app
}
