package rest

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error"`
	Message string      `json:"message,omitempty"`
	Status  int         `json:"-"` // omitted since it is available on response header
}

func (r response) respond(w http.ResponseWriter) {
	b, _ := json.Marshal(r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Status)
	w.Write(b)
}
