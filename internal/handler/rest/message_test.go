package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/ronaudinho/wp/internal/handler/rest"
	"github.com/ronaudinho/wp/internal/model"
)

// can use mocking library here
type serviceMock struct {
	CreateMessageFunc func(context.Context, model.Message) error
	GetMessagesFunc   func(context.Context) ([]model.Message, error)
}

func (sm *serviceMock) CreateMessage(ctx context.Context, msg model.Message) error {
	return sm.CreateMessageFunc(ctx, msg)
}

func (sm *serviceMock) GetMessages(ctx context.Context) ([]model.Message, error) {
	msg, err := sm.GetMessagesFunc(ctx)
	return msg, err
}

var (
	addr = "127.0.0.1:3195"
	ctx  = context.Background()
)

var createTests = []struct {
	name   string
	in     json.RawMessage
	errnil bool
}{
	{
		name:   "ok",
		in:     json.RawMessage(`{"message": "brainfcuk is fun"}`),
		errnil: true,
	},
	{
		name:   "validate required",
		in:     json.RawMessage(`{"msg": "brainfcuk is fun"}`),
		errnil: false,
	},
	{
		name:   "validate notoneof",
		in:     json.RawMessage(`{"message": "brainfuck is fun"}`),
		errnil: false,
	},
}

var getTests = []struct {
	name   string
	errnil bool
}{
	{
		name:   "ok",
		errnil: true,
	},
	{
		name:   "some error",
		errnil: false,
	},
}

func strptr(s string) *string {
	return &s
}

func TestREST_CreateMessage(t *testing.T) {
	path := "/message"
	target := "/message"
	for _, tt := range createTests {
		t.Run(tt.name, func(t *testing.T) {
			sm := &serviceMock{
				CreateMessageFunc: func(context.Context, model.Message) error {
					return nil
				},
			}
			mrst := rest.New(sm)
			rou := mux.NewRouter()
			rou.HandleFunc(path, mrst.CreateMessage).Methods("POST")
			b, _ := json.Marshal(tt.in)
			req, _ := http.NewRequest("POST", target, bytes.NewBuffer(b))
			rec := httptest.NewRecorder()
			rou.ServeHTTP(rec, req)
			if tt.errnil && (rec.Code != http.StatusOK) {
				t.Errorf("got %d, want %d", rec.Code, http.StatusOK)
			}
			// TODO check response body?
		})
	}
}

// a bit unnecessary at this point
func TestREST_GetMessages(t *testing.T) {
	path := "/message"
	target := "/message"
	for _, tt := range createTests {
		t.Run(tt.name, func(t *testing.T) {
			sm := &serviceMock{
				GetMessagesFunc: func(context.Context) ([]model.Message, error) {
					if !tt.errnil {
						return []model.Message{}, errors.New("mock error")
					}
					return []model.Message{}, nil
				},
			}
			mrst := rest.New(sm)
			rou := mux.NewRouter()
			rou.HandleFunc(path, mrst.GetMessages).Methods("GET")
			req, _ := http.NewRequest("GET", target, nil)
			rec := httptest.NewRecorder()
			rou.ServeHTTP(rec, req)
			if tt.errnil && (rec.Code != http.StatusOK) {
				t.Errorf("got %d, want %d", rec.Code, http.StatusOK)
			}
			// TODO check response body?
		})
	}
}
