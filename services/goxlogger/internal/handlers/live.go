package handlers

import (
	"encoding/base64"
	"log"
	"net/http"
	"path"

	"github.com/Zadigo/goxlogger/internal/models"
	"github.com/Zadigo/goxlogger/internal/tickerapp"
	"github.com/Zadigo/goxlogger/internal/utils"
)

type BaseRouteHandlers struct {
	models.BaseHandler
	serverConfig *utils.ServerConfig
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

	tickerapp.NewFileRedis(h.GetApp().GetAppContext(), h.GetApp().GetRedisClient())

	for {
		var message any
		err := conn.ReadJSON(&message)

		if err != nil {
			break
		}
	}
}

// GetFilesHandler handles the HTTP request to retrieve the list of log files. It first checks 
// the Redis cache for the files, and if not found, collects them from the "/data" folder.
func (h *BaseRouteHandlers) GetFilesHandler(w http.ResponseWriter, r *http.Request) {
	httpErrors := HttpErrors{}

	filesRedis := tickerapp.NewFileRedis(h.GetApp().GetAppContext(), h.GetApp().GetRedisClient())
	files, err := filesRedis.GetCachedFiles()
	if err != nil {
		httpErrors.FailedToCollectFiles(w, err)
		return
	}

	
	if len(files) == 0 {
		filesCollector := &tickerapp.FileCollector{}
		files, err = filesCollector.CollectFilesInFolder(h.GetApp().GetRootDir(), "/data")
		if err != nil {
			httpErrors.FailedToCollectFiles(w, err)
			return
		}

		filesRedis.SaveFiles(files)
	}

	if len(files) == 0 {
		// If no files are found at all, return an empty
		// array instead of null
		files = []models.File{}
	}

	utils.JsonResponse(w, files, http.StatusOK)
}

func (h *BaseRouteHandlers) GetLogsHandler(w http.ResponseWriter, r *http.Request) {
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

	fileRedis := tickerapp.NewFileRedis(h.GetApp().GetAppContext(), h.GetApp().GetRedisClient())

	var logs []tickerapp.LogLine

	// Check if the cached data for the file exists in Redis
	result := fileRedis.HasCachedData(string(decodedFileName))
	if !result {
		fullPath := path.Join(h.GetApp().GetRootDir(), "data", string(decodedFileName))
		strLogs, err := fileRedis.ReadFile(fullPath, h.serverConfig)
		if err != nil {
			log.Println("🔴 Failed to read file:", err)
			httpErrors.FailedToReadFile(w, err)
			return
		}

		logRedis := tickerapp.NewLogsRedis(h.GetApp().GetAppContext(), h.GetApp().GetRedisClient())
		logs, err = logRedis.SaveTransform(strLogs)
		if err != nil {
			httpErrors.FailedToReadFile(w, err)
			return
		}
	} else {
		if logs, err = fileRedis.GetLogs(string(decodedFileName)); err != nil {
			httpErrors.FailedToGetLogs(w, err)
			return
		}
	}

	logs, err = resolveQuery(r, logs)
	if err != nil {
		httpErrors.FailedToGetLogs(w, err)
		return
	}

	// Pagination
	logs, err = PaginateData(r, logs)
	if err != nil {
		httpErrors.InvalidLimitOffset(w, err)
		return
	}

	utils.JsonResponse(w, logs, http.StatusOK)
}
