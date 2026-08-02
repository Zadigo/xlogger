package tests

import (
	"testing"

	"github.com/Zadigo/goxlogger/internal/httpapp"
	"github.com/Zadigo/goxlogger/internal/logic"
	"github.com/redis/go-redis/v9"
)

func TestLogs(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	serverConfig := httpapp.LoadConfig("../config.yaml")
	logs := logic.NewLogsService(t.Context())
	logs.StartServer(serverConfig, redisClient)
}
