package handlers

import (
	"time"

	"github.com/gorilla/websocket"
)

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
