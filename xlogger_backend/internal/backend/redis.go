package backend

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(serverConfig *ServerConfig) *redis.Client {
	options, err := redis.ParseURL(serverConfig.Config.Redis.Url)
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(options)

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	return client
}
