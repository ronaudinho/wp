package service_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/ronaudinho/wp/internal/model"
	"github.com/ronaudinho/wp/internal/service"
)

// can use mocking library here
type repoMock struct {
	CreateMessageFunc func(model.Message) error
	GetMessagesFunc   func() ([]model.Message, error)
}

func (rm *repoMock) CreateMessage(msg model.Message) error {
	return rm.CreateMessageFunc(msg)
}

func (rm *repoMock) GetMessages() ([]model.Message, error) {
	msg, err := rm.GetMessagesFunc()
	return msg, err
}

var (
	addr = "127.0.0.1:3195"
	ctx  = context.Background()
)

var createTests = []struct {
	name   string
	in     model.Message
	errnil bool
}{
	{
		name:   "ok",
		in:     model.Message{Message: strptr("brainfcuk is fun")},
		errnil: true,
	},
	{
		name:   "mock database error",
		in:     model.Message{Message: strptr("brainfcuk is fun")},
		errnil: false,
	},
}

var getTests = []struct {
	name   string
	in     []model.Message
	out    []model.Message
	errnil bool
}{
	{
		name: "ok",
		in: []model.Message{
			model.Message{Message: strptr("brainfcuk is fun")},
			model.Message{Message: strptr("hi")},
		},
		out: []model.Message{
			model.Message{Message: strptr("brainfcuk is fun")},
			model.Message{Message: strptr("hi")},
		},
		errnil: true,
	},
	{
		name: "output not as expected",
		in: []model.Message{
			model.Message{Message: strptr("brainfcuk is fun")},
			model.Message{Message: strptr("hi")},
			model.Message{Message: strptr("really?")},
		},
		out: []model.Message{
			model.Message{Message: strptr("brainfcuk is fun")},
			model.Message{Message: strptr("hi")},
		},
		errnil: false,
	},
}

func strptr(s string) *string {
	return &s
}

// a bit unnecessary at this point
// since there is not much business logic that goes here
func TestService_CreateMessage(t *testing.T) {
	for _, tt := range createTests {
		t.Run(tt.name, func(t *testing.T) {
			rm := &repoMock{
				CreateMessageFunc: func(model.Message) error {
					if !tt.errnil {
						return errors.New("mock error")
					}
					return nil
				},
			}
			msvc := service.New(rm, addr)
			err := msvc.CreateMessage(ctx, tt.in)
			if tt.errnil != (err == nil) {
				t.Errorf("got errnil %v, want errnil %v", err != nil, tt.errnil)
			}
		})
	}
}

// a bit unnecessary at this point
// since there is not much business logic that goes here
func TestService_GetMessages(t *testing.T) {
	for _, tt := range getTests {
		t.Run(tt.name, func(t *testing.T) {
			rm := &repoMock{
				GetMessagesFunc: func() ([]model.Message, error) {
					var err error
					if !reflect.DeepEqual(tt.in, tt.out) {
						err = errors.New("mock error")
					}
					return tt.out, err
				},
			}
			msvc := service.New(rm, addr)
			msg, err := msvc.GetMessages(ctx)
			if tt.errnil != (err == nil) {
				t.Errorf("got errnil %v, want errnil %v", err != nil, tt.errnil)
			}
			if tt.errnil && !reflect.DeepEqual(msg, tt.in) {
				t.Errorf("got %v, want %v", msg, tt.in)
			}
		})
	}
}
