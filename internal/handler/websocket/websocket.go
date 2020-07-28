package websocket

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

// WebSocket
type WebSocket struct {
	client    map[*websocket.Conn]bool
	broadcast chan []byte
}

// New creates new instance of WebSocket
func New() *WebSocket {
	return &WebSocket{
		client:    make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte),
	}
}

// WithRoutes sets routes for WebSocket handler
func (ws *WebSocket) WithRoutes(r *mux.Router) {
	r.HandleFunc("/websocket", ws.ListenMessage).Methods("GET")
}
