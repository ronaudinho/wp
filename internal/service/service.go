package service

import (
	"github.com/ronaudinho/wp/internal/model"
)

// Service layer is supposedly for business logic implementation
// however, since there is not much business logic in this application
// this layer is not used much
type Service struct {
	repository Repository
	// addr tells the address of the websocket
	// it feels weird to put it here though
	addr string
}

// Repository defines interface that needs to be implemented in repository layer
type Repository interface {
	CreateMessage(model.Message) error
	GetMessages() ([]model.Message, error)
}

// New creates new instance of Service
func New(r Repository, addr string) *Service {
	return &Service{
		repository: r,
		addr:       addr,
	}
}
