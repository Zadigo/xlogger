package handlers

import (
	"net/http"

	"github.com/redis/go-redis/v9"
)

func LiveWsHandler(w http.ResponseWriter, r *http.Request, redisClient *redis.Client) {
	// Do something
}
