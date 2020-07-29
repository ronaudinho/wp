package websocket

import (
	"bytes"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub
// readPump is run in a per-connection goroutine to ensures that
// there is at most one reader on a connection
// by executing all reads from this goroutine
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("websocket error: %v", err)
			}
			break
		}
		msg = bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))
		c.hub.broadcast <- msg
	}
}

// writePump pumps messages from the hub to the websocket connection
// writePump is started for each connection to ensure
// that there is at most one writer to a connection
// by executing all writes from this goroutine
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			// TODO: delete this
			// simply check if message is sent to client
			// since my mental model is a bit murky right now
			// will print twice if the websocket in browser is active
			// since message is sent to both
			// short-lived client that wraps CreateMessage
			// and long-lived client in browser (for display purpose)
			log.Printf("successfully sent to %s: %s", c.conn.LocalAddr(), string(msg))
			// send msg to the current websocket
			w.Write(msg)
			// msg add queued chat messages to the current websocket message.
			// perhaps for long queue of messages
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}
			err = w.Close()
			if err != nil {
				return
			}
		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
