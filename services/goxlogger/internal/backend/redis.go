package backend

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewRedisBackend(ctx context.Context) *redis.Client {
	instance := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	err := instance.Ping(ctx).Err()
	if err != nil {
		panic(err)
	}

	return instance
}
