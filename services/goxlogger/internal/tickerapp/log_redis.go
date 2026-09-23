package tickerapp

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/redis/go-redis/v9"
)

// LogRedis contains all the logic to save,
// retrieve and manage logs in Redis
type LogRedis struct {
	ctx         context.Context
	redisClient *redis.Client
	broadcastCh chan LogLine
	mu          sync.Mutex
	Key         string
}

// Transform transforms the given string logs into LogLine structs
func (l *LogRedis) Transform(strLogs []string) []LogLine {
	logLines := make([]LogLine, len(strLogs))

	for i, strLog := range strLogs {
		instance := LogLine{RawLine: strLog}
		_, err := instance.ParseLine()

		if err == nil {
			logLines[i] = instance
		}
	}

	return logLines
}

// SaveTransform transforms the given string logs into LogLine structs and saves them in Redis
func (l *LogRedis) SaveTransform(strLines []string) (lines []LogLine, err error) {
	values := make([]any, len(strLines))

	logLines := l.Transform(strLines)

	for i, logLine := range logLines {
		if data, err := json.Marshal(logLine); err == nil {
			values[i] = data
		} else {
			return nil, err
		}
	}

	cmd := l.redisClient.SAdd(l.ctx, l.Key, values...)
	return logLines, cmd.Err()
}

// NewLogsRedis creates a new instance of LogRedis that is used to manage logs in Redis
func NewLogsRedis(ctx context.Context, redisClient *redis.Client) *LogRedis {
	return &LogRedis{
		ctx:         ctx,
		redisClient: redisClient,
		broadcastCh: make(chan LogLine, 100),
		Key:         "go-xlogger:all_logs",
	}
}
