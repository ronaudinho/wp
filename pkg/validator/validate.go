package validator

import (
	"encoding/json"
	"reflect"
	"strings"
)

// ValidationError
type ValidationError struct {
	Err []error
}

// Validate validates struct having fields with tags `validate`
func Validate(i interface{}) error {
	err := &ValidationError{Err: []error{}}
	v := reflect.ValueOf(i)
	t := reflect.TypeOf(i)
	for j := 0; j < v.NumField(); j++ {
		if t.Field(j).Tag.Get("validate") == "" {
			continue
		}
		err.validate(t.Field(j), v.Field(j))
	}
	// a hacky way to say return err as nil
	if len(err.Err) == 0 {
		return nil
	}
	return err
}

// Error ensures ValidationError implements error interface
func (ve *ValidationError) Error() string {
	err := []string{}
	for _, e := range ve.Err {
		err = append(err, e.Error())
	}
	j := struct {
		E []string `json:"validation"`
	}{
		E: err,
	}
	b, _ := json.Marshal(j)
	return string(b)
}

// Validate validates a given struct field
// idk might not be the best way to structure this
func (ve *ValidationError) validate(f reflect.StructField, v reflect.Value) {
	rul := strings.Split(f.Tag.Get("validate"), " ")
	for _, r := range rul {
		switch strings.Split(r, "=")[0] {
		case "required":
			e := validateRequired(f, v)
			if e != nil {
				ve.Err = append(ve.Err, e)
			}
		case "notoneof":
			e := validateNotOneOf(f, v, strings.Split(r, "=")[1])
			if e != nil {
				ve.Err = append(ve.Err, e)
			}
		default:
			// TODO maybe log to tell that rule is not defined yet?
			continue
		}
	}
}
