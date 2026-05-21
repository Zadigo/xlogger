package logic

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/redis/go-redis/v9"
)

// LogRedis contains all the logic to save, retrieve
// and manage logs in Redis
type LogRedis struct {
	ctx         context.Context
	redisClient *redis.Client
	broadcastCh chan LogLine
	mu          sync.Mutex
	Key         string
}

func (l *LogRedis) SaveLogs(logLines []LogLine) error {
	var values []any
	for _, logLine := range logLines {
		data, err := json.Marshal(logLine)
		if err != nil {
			return err
		}
		values = append(values, data)
	}
	cmd := l.redisClient.SAdd(l.ctx, l.Key, values...)
	return cmd.Err()
}

func (l *LogRedis) GetLogs() ([]LogLine, error) {
	cmd := l.redisClient.SMembers(l.ctx, l.Key)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	var logLines []LogLine
	for _, data := range cmd.Val() {
		var logLine LogLine
		if err := json.Unmarshal([]byte(data), &logLine); err != nil {
			return nil, err
		}
		logLines = append(logLines, logLine)
	}
	return logLines, nil
}

func (l *LogRedis) DeleteLogs() error {
	cmd := l.redisClient.Del(l.ctx, l.Key)
	return cmd.Err()
}

func (l *LogRedis) BroadcastLog(logLine LogLine) {
	l.mu.Lock()
	l.broadcastCh <- logLine
	l.mu.Unlock()
}

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

func NewLogsRedis(ctx context.Context, redisClient *redis.Client) *LogRedis {
	return &LogRedis{
		ctx:         ctx,
		redisClient: redisClient,
		broadcastCh: make(chan LogLine, 100),
		Key:         "go-xlogger:all_logs",
	}
}
