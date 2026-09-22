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
	Start()
}

type MainServerInterface interface {
	BaseServerInterface
	GetContext() context.Context
	GetConfig() *utils.ServerConfig
	GetRedisDb() *redis.Client
}
