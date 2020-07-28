package service

import (
	"context"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/ronaudinho/wp/internal/model"
)

// CreateMessage creates message from request
func (s *Service) CreateMessage(ctx context.Context, msg model.Message) error {
	return s.repository.CreateMessage(msg)
}

// GetMessages get all previously sent messages
func (s *Service) GetMessages(ctx context.Context) ([]model.Message, error) {
	msg, err := s.repository.GetMessages()
	return msg, err
}

// PushMessage pushes created message to websocket to display in realtime
func PushMessage(msg string) {
	addr := "127.0.0.1:3195"
	u := url.URL{Scheme: "ws", Host: addr, Path: "/websocket"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println("dial:", err)
	}
	defer c.Close()
	err = c.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Println("write:", err)
	}
}
