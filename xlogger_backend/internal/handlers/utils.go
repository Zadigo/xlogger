package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
)

var allowedOrigins = map[string]bool{
	"http://localhost:3000": true,
	"https://example.com":   true,
}

var CustomRequestUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(request *http.Request) bool {
		origin := request.Header.Get("Origin")

		_, ok := allowedOrigins[origin]
		if !ok {
			return false
		}

		return allowedOrigins[origin]
	},
}

func IsWebsocketClose(err error) bool {
	if websocket.IsCloseError(err,
		websocket.CloseNormalClosure,   // 1000
		websocket.CloseGoingAway,       // 1001
		websocket.CloseAbnormalClosure, // 1006
	) {
		return true
	}

	// Also catches abrupt disconnects
	// (io.EOF, reset by peer, etc.)
	if errors.Is(err, io.EOF) {
		return true
	}

	return false
}
