package rest

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ronaudinho/wp/internal/model"
)

// REST
type REST struct {
	service Service
}

// Service defines interface that needs to be implemented in service layer
type Service interface {
	CreateMessage(context.Context, model.Message) error
	GetMessages(context.Context) ([]model.Message, error)
}

// New creates new instance of REST handler
func New(s Service) *REST {
	return &REST{service: s}
}

// WithRoutes sets routes for REST handler
func (rst *REST) WithRoutes(rou *mux.Router) {
	rou.HandleFunc("/healthcheck", healthcheck).Methods("GET")
	rou.HandleFunc("/message", rst.GetMessages).Methods("GET")
	rou.HandleFunc("/message", rst.CreateMessage).Methods("POST")
	rou.HandleFunc("/message/live", rst.LiveMessage).Methods("GET")
}

// healthcheck endpoint for checking if the application is up
func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
