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
	// call websocket to inform that a message has been created
	go pushMessage(s.addr, *msg.Message)
	return s.repository.CreateMessage(msg)
}

// GetMessages get all previously sent messages
func (s *Service) GetMessages(ctx context.Context) ([]model.Message, error) {
	msg, err := s.repository.GetMessages()
	return msg, err
}

// pushMessage pushes created message to websocket to display in realtime
// short-lived websocket client that acts as a wrapper
// since message is actually created on a different endpoint
func pushMessage(addr, msg string) {
	u := url.URL{Scheme: "ws", Host: addr, Path: "/websocket"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println("dial:", err)
	}
	defer conn.Close()
	err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Println("error writing:", err)
	}
	// TODO: delete this
	// simply check if message is sent to client
	// since my mental model is a bit murky right now
	_, res, err := conn.ReadMessage()
	if err != nil {
		log.Printf("error reading: %v", err)
	}
	log.Printf("successfully recv from %s%s: %s", u.Host, u.Path, string(res))
}
