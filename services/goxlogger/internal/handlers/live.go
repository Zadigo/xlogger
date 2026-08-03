package handlers

import (
	"encoding/base64"
	"net/http"
	"path"

	"github.com/Zadigo/goxlogger/internal/models"
	"github.com/Zadigo/goxlogger/internal/tickerapp"
	"github.com/Zadigo/goxlogger/internal/utils"
)

type BaseRouteHandlers struct {
	app          models.AppInterface
	serverConfig *utils.ServerConfig
}

func (h *BaseRouteHandlers) SetApp(app models.AppInterface) {
	h.app = app
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

	tickerapp.NewFileRedis(h.app.GetAppContext(), h.app.GetRedisClient())

	for {
		var message any
		err := conn.ReadJSON(&message)

		if err != nil {
			break
		}
	}
}

func (h *BaseRouteHandlers) GetFiles(w http.ResponseWriter, r *http.Request) {
	filesRedis := tickerapp.NewFileRedis(h.app.GetAppContext(), h.app.GetRedisClient())
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

	fileRedis := tickerapp.NewFileRedis(h.app.GetAppContext(), h.app.GetRedisClient())

	var logs []*tickerapp.LogLine

	// Check if the cached data for the file exists in Redis
	result := fileRedis.HasCachedData(string(decodedFileName))
	if !result {
		fullPath := path.Join(h.app.GetRootDir(), "data", string(decodedFileName))
		strLogs, err := fileRedis.ReadFile(fullPath, h.serverConfig)
		if err != nil {
			httpErrors.FailedToReadFile(w)
			return
		}

		logRedis := tickerapp.NewLogsRedis(h.app.GetAppContext(), h.app.GetRedisClient())
		logs, err = logRedis.SaveTransform(strLogs)
		if err != nil {
			httpErrors.FailedToReadFile(w, err)
			return
		}
	} else {
		if logs, err = fileRedis.GetLogs(string(decodedFileName)); err != nil {
			httpErrors.FailedToGetLogs(w)
			return
		}
	}

	utils.JsonResponse(w, logs, http.StatusOK)
}
