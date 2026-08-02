package models

import "github.com/redis/go-redis/v9"

type AppInterface interface {
	Start() error
	GetRedisClient() *redis.Client
}
