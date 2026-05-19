package handlers

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// CORS middleware to handle cross-origin requests
func Cors(next http.Handler) http.Handler {
	handle := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(handle)
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
