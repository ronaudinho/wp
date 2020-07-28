package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// ListenMessage listens for message using websocket protocol
func (ws *WebSocket) ListenMessage(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	ws.client[conn] = true
	// would probably be better to use hub like chat examples in gorilla/websocket repo
	// reading message from websocket
	go func() {
		defer conn.Close()
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("error reading: %v", err)
				delete(ws.client, conn)
				break
			}
			ws.broadcast <- msg
		}
	}()
	// broadcasting message to client
	go func() {
		defer conn.Close()
		for {
			msg := <-ws.broadcast
			for c := range ws.client {
				err := c.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					log.Printf("error writing: %v", err)
					delete(ws.client, c)
				}
			}
		}
	}()
}
