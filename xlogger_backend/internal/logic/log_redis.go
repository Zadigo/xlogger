package logic

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// LogRedis contains all the logic to save, retrieve
// and manage logs in Redis
type LogRedis struct {
	ctx         context.Context
	redisClient *redis.Client
}

func (l *LogRedis) SaveLogs(logLines []LogLine) error {
	// Save the log lines in Redis
	return nil
}

func (l *LogRedis) GetLogs() ([]LogLine, error) {
	// Retrieve logs from Redis and return them as a slice of LogLine
	return nil, nil
}
