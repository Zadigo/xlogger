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
}

func (l *LogRedis) SaveLogs(logLines []LogLine) error {
	data, err := json.Marshal(logLines)
	if err != nil {
		return err
	}
	cmd := l.redisClient.SAdd(l.ctx, "all_logs", data)
	return cmd.Err()
}

func (l *LogRedis) GetLogs() ([]LogLine, error) {
	cmd := l.redisClient.SMembers(l.ctx, "all_logs")
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	var logLines []LogLine
	for _, data := range cmd.Val() {
		var logs []LogLine
		if err := json.Unmarshal([]byte(data), &logs); err != nil {
			return nil, err
		}
		logLines = append(logLines, logs...)
	}
	return logLines, nil
}

func (l *LogRedis) DeleteLogs() error {
	cmd := l.redisClient.Del(l.ctx, "all_logs")
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
	}
}
