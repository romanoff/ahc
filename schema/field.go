package schema

import (
	"errors"
	"fmt"
	"reflect"
)

const (
	StringField = iota
	NumField    = iota
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
	ObjectFields  []*Field //Fields that are supposed to be in object (optional)
	ArrayValues   *Field   //Field that ArrayField should consist of (optional)
}

func (self *Field) Validate(value interface{}) error {
	switch self.Type {
	case StringField:
		_, success := value.(string)
		if !success {
			return errors.New(fmt.Sprintf("Expected string value, but got %v", value))
		}
	case NumField:
		kind := reflect.TypeOf(value).Kind()
		if kind != reflect.Int && kind != reflect.Int64 && kind != reflect.Float64 {
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
		if self.ArrayValues == nil {
			return nil
		}
		arrayValues := reflect.ValueOf(value)
		for i := 0; i < arrayValues.Len(); i++ {
			arrayValue := arrayValues.Index(i).Interface()
			err := self.ArrayValues.Validate(arrayValue)
			if err != nil {
				return err
			}
		}
	case ObjectField:
		kind := reflect.TypeOf(value).Kind()
		if kind != reflect.Map {
			return errors.New(fmt.Sprintf("Expected object value, but got %v", value))
		}
		if self.ObjectFields == nil {
			return nil
		}
		valueMap, success := value.(map[string]interface{})
		if !success {
			return errors.New(fmt.Sprintf("Expected object value, but got %v", value))
		}
		for _, objectField := range self.ObjectFields {
			objectFieldValue := valueMap[objectField.Name]
			if objectFieldValue == nil {
				return errors.New(objectField.Name + " key is missing")
			}
			err := objectField.Validate(objectFieldValue)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
