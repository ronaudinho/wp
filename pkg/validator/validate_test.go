package validator_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ronaudinho/wp/pkg/validator"
)

var vldtTests = []struct {
	name string
	in   interface{} // input to validate
	want error       // expected output
}{
	{
		name: "ok",
		in: struct {
			Message *string `json:"message" validate:"required"`
		}{
			Message: strptr("brainfuck is fun"),
		},
		want: nil,
	},
	{
		name: "required error",
		in: struct {
			Message *string `validate:"required"`
		}{},
		want: &validator.ValidationError{
			Err: []error{
				validator.RequiredError{
					Field:     "Message",
					FieldJSON: "",
				},
			},
		},
	},
	{
		name: "empty ok",
		in: struct {
			Message *string
		}{},
		want: nil,
	},
	{
		name: "notoneof empty ok",
		in: struct {
			Message *string `validate:"notoneof=fuck"`
		}{},
		want: nil,
	},
	{
		name: "notoneof error",
		in: struct {
			Message *string `validate:"notoneof=fuck,fck"`
		}{
			Message: strptr("brainfuck is fun"),
		},
		want: &validator.ValidationError{
			Err: []error{
				validator.NotOneOfError{
					Field:     "Message",
					FieldJSON: "",
					Word:      "fuck,fck",
				},
			},
		},
	},
	{
		name: "notoneof ok",
		in: struct {
			Message *string `validate:"notoneof=fuck,fck"`
		}{
			Message: strptr("brainfcuk is fun"),
		},
		want: nil,
	},
}

var vldtTestsJSON = []struct {
	name string
	json []byte // json input to validate
	want error  // expected output
}{
	{
		name: "ok",
		json: json.RawMessage(`{"message": "brainfcuk is fun"}`),
		want: nil,
	},
	{
		name: "required error",
		json: json.RawMessage(`{"msg": "brainfuck is fun"}`),
		want: &validator.ValidationError{
			Err: []error{
				validator.RequiredError{
					Field:     "Message",
					FieldJSON: "message",
				},
			},
		},
	},
	{
		name: "notoneof error",
		json: json.RawMessage(`{"msg": "brainfuck is fun","message": "brainfuck is fun"}`),
		want: &validator.ValidationError{
			Err: []error{
				validator.NotOneOfError{
					Field:     "Message",
					FieldJSON: "message",
					Word:      "fuck,fck",
				},
			},
		},
	},
}

func strptr(s string) *string {
	return &s
}

func TestValidate(t *testing.T) {
	for _, tt := range vldtTests {
		t.Run(tt.name, func(t *testing.T) {
			got := validator.Validate(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestValidateFromJSON(t *testing.T) {
	type msg struct {
		Message *string `json:"message" validate:"required notoneof=fuck,fck"`
	}
	for _, tt := range vldtTestsJSON {
		t.Run(tt.name, func(t *testing.T) {
			var m msg
			json.Unmarshal(tt.json, &m)
			got := validator.Validate(m)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}
