package validator

import (
	"fmt"
	"reflect"
	"strings"
)

// NotOneOfError
type NotOneOfError struct {
	FieldJSON string
	Word      string
}

// Error ensures NotOneOfError implements error interface
func (noo NotOneOfError) Error() string {
	return fmt.Sprintf("not one of: field %s contains banned word: %s", noo.FieldJSON, noo.Word)
}

// validateNotOneOf validates struct field with tag `validate` that contains `notoneof`
func validateNotOneOf(f reflect.StructField, v reflect.Value, noo string) error {
	var err error
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for _, w := range strings.Split(noo, ",") {
		if strings.Contains(v.String(), w) {
			err = NotOneOfError{
				FieldJSON: f.Tag.Get("json"),
				Word:      noo,
			}
		}
	}
	return err
}
