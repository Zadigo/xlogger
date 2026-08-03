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
func (l *LogRedis) Transform(strLogs []string) []*LogLine {
	var logLines []*LogLine

	for _, strLog := range strLogs {
		instance := &LogLine{RawLine: strLog}
		_, err := instance.ParseLine()

		if err != nil {
			logLines = append(logLines, instance)
		}
	}

	return logLines
}

// SaveTransform transforms the given string logs into LogLine structs and saves them in Redis
func (l *LogRedis) SaveTransform(strLines []string) (lines []*LogLine, err error) {
	var values []any

	logLines := l.Transform(strLines)

	for _, logLine := range logLines {
		if data, err := json.Marshal(logLine); err == nil {
			values = append(values, data)
		} else {
			return nil, err
		}
	}

	cmd := l.redisClient.SAdd(l.ctx, l.Key, values...)
	return logLines, cmd.Err()
}

// Deprecated: use files_redis.go
func (l *LogRedis) GetLogs() error {
	cmd := l.redisClient.SMembers(l.ctx, l.Key)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	var logLines []LogLine
	for _, data := range cmd.Val() {
		var logLine LogLine
		if err := json.Unmarshal([]byte(data), &logLine); err != nil {
			return err
		}
		logLines = append(logLines, logLine)
	}
	return nil
}

// Deprecated: use files_redis.go
func (l *LogRedis) DeleteLogs() error {
	cmd := l.redisClient.Del(l.ctx, l.Key)
	return cmd.Err()
}

// Deprecated: use files_redis.go
func (l *LogRedis) BroadcastLog(logLine LogLine) {
	l.mu.Lock()
	l.broadcastCh <- logLine
	l.mu.Unlock()
}

// Deprecated: use files_redis.go
func (l *LogRedis) StartBroadcaster() <-chan error {
	ch := make(chan error, 1)

	go func() {
		pubSub := l.redisClient.Subscribe(l.ctx, "logs_channel")

		for {
			select {
			case logLine := <-l.broadcastCh:
				data, err := json.Marshal(logLine)
				if err != nil {
					ch <- err
					continue
				}
				l.redisClient.Publish(l.ctx, "logs_channel", data)
			case <-pubSub.Channel():
			// case msg := <-redisCh:
			// 	var logLine LogLine
			// 	if err := json.Unmarshal([]byte(msg.Payload), &logLine); err != nil {
			// 		continue
			// 	}
			// 	l.mu.Lock()
			// 	l.broadcastCh <- logLine
			// 	l.mu.Unlock()
			// }

			case <-l.ctx.Done():
				ch <- pubSub.Close()
				return
			}
		}
	}()

	return ch
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
