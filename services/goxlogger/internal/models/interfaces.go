package models

import (
	"context"
	"net/http"

	"github.com/Zadigo/goxlogger/internal/utils"
	"github.com/redis/go-redis/v9"
)

type ErrorInterface interface {
	LogErrorMessage(err ...error)
	SendErrorMessage(w http.ResponseWriter, err ...error)
}

type BaseServerInterface interface {
	Start(ctx context.Context) error
}

type MainServerInterface interface {
	BaseServerInterface
	GetContext() context.Context
	GetConfig() *utils.ServerOptions
	GetRedisDb() *redis.Client
}

type AppInterface interface {
	BaseServerInterface
	GetRedisClient() *redis.Client
}
