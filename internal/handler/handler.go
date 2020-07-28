package handler

import (
	"github.com/gorilla/mux"
)

// Handler
type Handler struct {
	handler []IHandler
}

// IHandler defines interface that needs to be implemented
// by type of handlers, namely REST and WebSocket in this codebase
// usage of I in front of actual struct name to signify interface
// differs from practice in other files since it is needed
// to define IHandler when creating new instance of Handler
// using New()
type IHandler interface {
	WithRoutes(*mux.Router)
}

// New creates a new instance of Handler
func New(h []IHandler) *Handler {
	return &Handler{handler: h}
}

// WithRoutes aggregates all routes set in each handler
// TODO: option to print down all routes ala Gin?
func (h *Handler) WithRoutes(rou *mux.Router) {
	for _, hnd := range h.handler {
		hnd.WithRoutes(rou)
	}
}
