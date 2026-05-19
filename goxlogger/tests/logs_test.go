package tests

import (
	"testing"

	"github.com/Zadigo/goxlogger/internal/logic"
	"github.com/Zadigo/goxlogger/internal/server"
	"github.com/redis/go-redis/v9"
)

func TestLogs(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	serverConfig := server.LoadConfig("../config.yaml")
	logs := logic.NewLogsService(t.Context(), "../data", true)
	logs.StartServer(serverConfig, redisClient)
}
