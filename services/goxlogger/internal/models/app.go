package models

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type AppInterface interface {
	Start() error
	GetRedisClient() *redis.Client
	GetAppContext() context.Context
}
