package handlers

import (
	"errors"
	"io"

	"github.com/gorilla/websocket"
)

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
