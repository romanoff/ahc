package schema

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
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
	AllowedValues []string //For StringField only
	ObjectFields  []*Field //Fields that are supposed to be in object (optional)
	ArrayValues   []*Field //Fields that ArrayField should consist of (optional)
}

func (self *Field) Validate(value interface{}) error {
	if value == nil && self.Required == false {
		return nil
	}
	switch self.Type {
	case StringField:
		strVal, success := value.(string)
		if !success {
			return errors.New(fmt.Sprintf("Expected string value, but got %v", value))
		}
		if self.AllowedValues != nil && !StringInSlice(strVal, self.AllowedValues) {
			return errors.New(fmt.Sprintf("Expected string value to be in (%v), but was %v",
				strings.Join(self.AllowedValues, ", "), strVal))
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
			if len(self.ArrayValues) == 0 {
				return nil
			}
			valueMap, success := arrayValue.(map[string]interface{})
			if !success {
				//TODO: Fix this. cannot convert map. Success is false
				return nil
				return errors.New(fmt.Sprintf("Expected object value, but got %v", arrayValue))
			}
			for _, arraySchemaValue := range self.ArrayValues {
				err := arraySchemaValue.Validate(valueMap[arraySchemaValue.Name])
				if err != nil {
					return err
				}
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

func (self *Field) Cast(value interface{}) interface{} {
	switch self.Type {
	case BoolField:
		stringValue, success := value.(string)
		if success {
			return stringValue == "true" || stringValue == "1"
		}
		intValue, success := value.(int)
		if success {
			return intValue == 1
		}
	case NumField:
		stringValue, success := value.(string)
		if success {
			floatValue, err := strconv.ParseFloat(stringValue, 32)
			if err == nil {
				return floatValue
			}
		}
		intValue, success := value.(int)
		if success {
			return float64(intValue)
		}
		floatValue, success := value.(float32)
		if success {
			return float64(floatValue)
		}
	}
	return value
}
