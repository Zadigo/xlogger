package handlers

import (
	"encoding/base64"
	"net/http"

	"github.com/Zadigo/goxlogger/internal/models"
	"github.com/Zadigo/goxlogger/internal/tickerapp"
	"github.com/Zadigo/goxlogger/internal/utils"
	"github.com/redis/go-redis/v9"
)

type BaseRouteHandlers struct {
	app          models.AppInterface
	rootDir      string
	redisClient  *redis.Client
	serverConfig *utils.ServerConfig
}

func (h *BaseRouteHandlers) SetApp(app models.AppInterface) {
	h.app = app
	h.redisClient = app.GetRedisClient()
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

	tickerapp.NewFileRedis(h.app.GetAppContext(), h.redisClient)

	for {
		var message any
		err := conn.ReadJSON(&message)

		if err != nil {
			break
		}
	}
}

func (h *BaseRouteHandlers) GetFiles(w http.ResponseWriter, r *http.Request) {
	filesRedis := tickerapp.NewFileRedis(h.app.GetAppContext(), h.redisClient)
	files, err := filesRedis.GetFiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, files, http.StatusOK)
}

func (h *BaseRouteHandlers) GetLogs(w http.ResponseWriter, r *http.Request) {
	httpErrors := HttpErrors{}

	fileId := r.Context().Value("fileId").(string)
	if fileId == "" {
		httpErrors.InvalidFileId(w)
		return
	}

	decodedFileName, err := base64.StdEncoding.DecodeString(fileId)
	if err != nil {
		httpErrors.InvalidFileId(w)
		return
	}

	fileRedis := tickerapp.NewFileRedis(h.app.GetAppContext(), h.redisClient)
	logs, err := fileRedis.GetLogs(string(decodedFileName))

	if err != nil {
		httpErrors.FailedToGetLogs(w)
		return
	}

	utils.JsonResponse(w, logs, http.StatusOK)
}
