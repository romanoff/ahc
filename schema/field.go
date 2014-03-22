package schema

import (
	"errors"
	"fmt"
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
	_, noError := value.(string)
	if !noError {
		return errors.New(fmt.Sprintf("Expected string value, but got %v", value))
	}
	return nil
}
