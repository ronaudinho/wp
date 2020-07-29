package websocket

import (
	"net/http"

	"github.com/gorilla/mux"
)

// WebSocket
type WebSocket struct {
	hub *Hub
}

// New creates new instance of WebSocket
// as well as runs Hub on a separate goroutine
func New() *WebSocket {
	hub := NewHub()
	go hub.run()
	return &WebSocket{
		hub: hub,
	}
}

// WithRoutes sets routes for WebSocket handler
func (ws *WebSocket) WithRoutes(r *mux.Router) {
	r.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		ws.listenMessage(ws.hub, w, r)
	}).Methods("GET")
}
