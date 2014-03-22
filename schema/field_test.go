package schema

import (
	"testing"
)

func TestStringFieldValidate(t *testing.T) {
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

func TestIntFieldValidate(t *testing.T) {
	field := &Field{Type: IntField}
	err := field.Validate(1)
	if err != nil {
		t.Errorf("Expected to get no int validation error, but got %v", err)
	}
	err = field.Validate("string field")
	if err == nil {
		t.Errorf("Expected to get int validation error (string supplied), but got nil")
	}
}

func TestBoolFieldValidate(t *testing.T) {
	field := &Field{Type: BoolField}
	err := field.Validate(true)
	if err != nil {
		t.Errorf("Expected to get no bool validation error, but got %v", err)
	}
	err = field.Validate("string field")
	if err == nil {
		t.Errorf("Expected to get bool validation error (string supplied), but got nil")
	}
}
