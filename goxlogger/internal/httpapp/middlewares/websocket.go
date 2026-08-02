package middlewares

import (
	"net/http"

	"github.com/Zadigo/goxlogger/internal/utils"
	"github.com/gorilla/websocket"
)

var CustomRequestUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(request *http.Request) bool {
		origin := request.Header.Get("Origin")

		_, ok := utils.AllowedOrigins[origin]
		if !ok {
			return false
		}

		return utils.AllowedOrigins[origin]
	},
}
