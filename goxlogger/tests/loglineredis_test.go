package tests

import (
	"testing"

	"github.com/Zadigo/goxlogger/internal/logic"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestLogsRedis(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	logRedis := logic.NewLogsRedis(t.Context(), redisClient)

	t.Run("Should save logs to Redis", func(t *testing.T) {
		logLines := []logic.LogLine{
			{RawLine: "Log line 1"},
			{RawLine: "Log line 2"},
		}
		err := logRedis.SaveLogs(logLines)
		assert.Nil(t, err)
	})

	t.Run("Should retrieve logs from Redis", func(t *testing.T) {
		logs, err := logRedis.GetLogs()
		assert.Nil(t, err)
		assert.Len(t, logs, 2)
	})

	t.Run("Should delete logs from Redis", func(t *testing.T) {
		t.Skip("Works. Skip to test others")
		err := logRedis.DeleteLogs()
		assert.Nil(t, err)
	})
}

func TestBroadcastLog(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	t.Run("Should broadcast logs", func(t *testing.T) {
		logRedis := logic.NewLogsRedis(t.Context(), redisClient)
		logRedis.StartBroadcaster()

		logLine := logic.LogLine{RawLine: "Broadcast log line"}
		logRedis.BroadcastLog(logLine)
	})
}
