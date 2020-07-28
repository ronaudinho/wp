package service

import (
	"github.com/ronaudinho/wp/internal/model"
)

// Service
// Service layer is supposedly for business logic implementation
// however, since there is not much business logic in this application
// this layer is not used much
type Service struct {
	repository Repository
}

// Repository defines interface that needs to be implemented in repository layer
type Repository interface {
	CreateMessage(model.Message) error
	GetMessages() ([]model.Message, error)
}

// New creates new instance of Service
func New(r Repository) *Service {
	return &Service{repository: r}
}
