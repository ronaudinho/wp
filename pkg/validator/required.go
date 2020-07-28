package validator

import (
	"fmt"
	"reflect"
)

// RequiredError
type RequiredError struct {
	Field     string
	FieldJSON string
}

// Error ensures RequiredError implements error interface
func (re RequiredError) Error() string {
	return fmt.Sprintf("required: field missing in JSON: %s", re.FieldJSON)
}

// validateRequired validates struct field with tag `validate` that contains `required`
func validateRequired(f reflect.StructField, v reflect.Value) error {
	var err error
	if v.IsNil() {
		err = RequiredError{
			Field:     f.Name,
			FieldJSON: f.Tag.Get("json"),
		}
	}
	return err
}
