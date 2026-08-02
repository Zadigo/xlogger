package handlers

import (
	"net/http"
	"time"

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

// WebsocketMiddleware sets up the necessary parameters for a
// websocket connection, such as read limits and deadlines.
type WebsocketMiddleware struct{}

func (m *WebsocketMiddleware) Handle(conn *websocket.Conn) {
	conn.SetReadLimit(1024)
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
}
