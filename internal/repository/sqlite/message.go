package sqlite

import (
	"log"
	"time"

	"github.com/ronaudinho/wp/internal/model"
)

// CreateMessaage insert message into SQLite database
// could have returned created message with the ID but i lazy
func (s *Repository) CreateMessage(msg model.Message) error {
	ut := time.Now().Unix()
	_, err := s.db.Exec("INSERT INTO message(message, create_at) VALUES(?, ?)", msg.Message, ut)
	return err
}

// GetMessages get all messages previously sent from SQLite database
func (s *Repository) GetMessages() ([]model.Message, error) {
	var msg []model.Message
	rows, err := s.db.Query("SELECT message, create_at FROM message")
	if err != nil {
		return msg, err
	}
	defer rows.Close()
	for rows.Next() {
		var m model.Message
		err = rows.Scan(&m.Message, &m.CreateAt)
		if err != nil {
			log.Println(err)
			break
		}
		msg = append(msg, m)
	}
	err = rows.Err()
	return msg, err
}
