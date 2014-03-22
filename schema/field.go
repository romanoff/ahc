package schema

import (
	"errors"
	"fmt"
	"reflect"
)

const (
	StringField = iota
	IntField    = iota
	BoolField   = iota
	ArrayField  = iota
	ObjectField = iota
)

type Field struct {
	Name          string
	Description   string
	Type          int
	Required      bool
	AllowedValues []string
}

func (self *Field) Validate(value interface{}) error {
	switch self.Type {
	case StringField:
		_, success := value.(string)
		if !success {
			return errors.New(fmt.Sprintf("Expected string value, but got %v", value))
		}
	case IntField:
		_, success := value.(int)
		if !success {
			return errors.New(fmt.Sprintf("Expected int value, but got %v", value))
		}
	case BoolField:
		_, success := value.(bool)
		if !success {
			return errors.New(fmt.Sprintf("Expected bool value, but got %v", value))
		}
	case ArrayField:
		kind := reflect.TypeOf(value).Kind()
		if kind != reflect.Array && kind != reflect.Slice {
			return errors.New(fmt.Sprintf("Expected array value, but got %v", value))
		}
	case ObjectField:
		kind := reflect.TypeOf(value).Kind()
		if kind != reflect.Map {
			return errors.New(fmt.Sprintf("Expected object value, but got %v", value))
		}
	}
	return nil
}
