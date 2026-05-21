package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Zadigo/goxlogger/internal/logic"
	"github.com/Zadigo/goxlogger/internal/models"
	"github.com/redis/go-redis/v9"
)

type BaseRouteHandlers struct {
	ctx          context.Context
	rooDir       string
	redisClient  *redis.Client
	serverConfig *models.ServerConfig
}

func (h *BaseRouteHandlers) LiveWsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := CustomRequestUpgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	middleware := WebsocketMiddleware{}
	middleware.Handle(conn)

	logic.NewFileRedis(h.ctx, h.rooDir, h.redisClient)

	for {
		var message any
		err := conn.ReadJSON(&message)

		if err != nil {
			break
		}
	}
}

func (h *BaseRouteHandlers) GetLogs(w http.ResponseWriter, r *http.Request) {
	logsRedis := logic.NewLogsRedis(h.ctx, h.redisClient)
	logs, err := logsRedis.GetLogs()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// data, err := json.Marshal(logs)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// w.Write(data)

	// Chunked transfer encoding
	// w.Header().Set("Transfer-Encoding", "chunked")
	// w.Header().Set("Content-Type", "text/plain")

	// flusher, ok := w.(http.Flusher)
	// if !ok {
	//     http.Error(w, "streaming not supported", http.StatusInternalServerError)
	//     return
	// }

	// w.Header().Set("Content-Type", "application/json")
	// // No need to set Transfer-Encoding manually, Go sets it automatically

	// encoder := json.NewEncoder(w)
	// for _, log := range logs {
	//     if err := encoder.Encode(log); err != nil {
	//         return // client likely disconnected
	//     }
	//     flusher.Flush() // sends the chunk immediately
	// }

	// SSE
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for i, log := range logs {
		data, err := json.Marshal(log)
		if err != nil {
			continue
		}

		// SSE format: "data: <payload>\n\n"
		fmt.Fprintf(w, "id: %d\nevent: log\ndata: %s\n\n", i, data)
		flusher.Flush()
	}

	// Optional: signal the client that the stream is done
	fmt.Fprintf(w, "event: done\ndata: {}\n\n")
	flusher.Flush()
}

func NewBaseRouteHandlers(ctx context.Context, rooDir string, redisClient *redis.Client, serverConfig *models.ServerConfig) *BaseRouteHandlers {
	return &BaseRouteHandlers{
		ctx:          ctx,
		rooDir:       rooDir,
		redisClient:  redisClient,
		serverConfig: serverConfig,
	}
}
