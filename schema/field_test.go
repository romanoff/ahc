package schema

import (
	"testing"
)

func TestFieldValidate(t *testing.T) {
	field := &Field{Type: StringField}
	err := field.Validate("field value")
	if err != nil {
		t.Errorf("Expected to get no string validation error, but got %v", err)
	}
	err = field.Validate(1)
	if err == nil {
		t.Errorf("Expected to get string validation error (int supplied), but got nil")
	}
}
