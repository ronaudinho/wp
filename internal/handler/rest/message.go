package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ronaudinho/wp/internal/model"
	svc "github.com/ronaudinho/wp/internal/service"
	"github.com/ronaudinho/wp/pkg/validator"
)

// CreateMessage creates message from request
func (rst *REST) CreateMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response{
			Error:  err.Error(),
			Status: http.StatusBadRequest,
		}.respond(w)
		return
	}
	defer r.Body.Close()
	var msg model.Message
	json.Unmarshal(b, &msg)
	err = validator.Validate(msg)
	if err != nil {
		response{
			Error:  err.Error(),
			Status: http.StatusBadRequest,
		}.respond(w)
		return
	}
	err = rst.service.CreateMessage(ctx, msg)
	if err != nil {
		response{
			Error:  err.Error(),
			Status: http.StatusInternalServerError,
		}.respond(w)
		return
	}
	//  call websocket before responding
	go svc.PushMessage(*msg.Message)
	response{
		Data:   msg,
		Error:  nil,
		Status: http.StatusOK,
	}.respond(w)
	return
}

// GetMessages get all previously sent messages
// TODO: optional QueryString for filtering messages?
func (rst *REST) GetMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	msg, err := rst.service.GetMessages(ctx)
	if err != nil {
		response{
			Error:  err.Error(),
			Status: http.StatusInternalServerError,
		}.respond(w)
		return
	}
	response{
		Data:   msg,
		Error:  nil,
		Status: http.StatusOK,
	}.respond(w)
	return
}

// ListenMessage displays messages as they come in live
// rendered in live.html file
// it is not actually an API endpoint to be consumed but
// simply a visual help for seeing websocket in action
func (rst *REST) LiveMessage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "live.html")
}
