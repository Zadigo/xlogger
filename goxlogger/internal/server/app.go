package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
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

func (a *App) Start(ctx context.Context) error {
	a.ctx = ctx

	// Redis client
	err := a.redisClient.Ping(a.ctx).Err()
	if err != nil {
		return fmt.Errorf("🔴 Failed to load Redis: %w", err)
	}

	// Log server
	logServer := &logic.Logs{}
	logServer.StartServer(ctx, a.config, a.redisClient)

	defer func() {
		a.redisClient.Close()
		log.Print("🔴 Redis client closed")
	}()

	// HTTP server
	server := http.Server{
		Addr:    ":8080",
		Handler: a.router,
	}

	ch := make(chan error, 1)

	go func() {
		log.Print("🟢 Server is listening on port 8080")
		ch <- server.ListenAndServe()
	}()

	select {
	case err := <-ch:
		return err

	case <-ctx.Done():
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return server.Shutdown(timeoutCtx)
	}
}

func NewApp(config *models.ServerConfig) *App {
	app := &App{
		config: config,
		redisClient: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
			DB:   0,
		}),
	}
	app.loadRoutes()
	return app
}
