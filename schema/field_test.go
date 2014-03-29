package schema

import (
	"encoding/json"
	"math"
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

func TestRequiredFieldValidate(t *testing.T) {
	field := &Field{Type: StringField, Required: true}
	err := field.Validate(nil)
	if err == nil {
		t.Errorf("Expected to get string validation error (nil supplied), but got nil")
	}
	field.Required = false
	err = field.Validate(nil)
	if err != nil {
		t.Errorf("Expected to get no string validation error, but got %v", err)
	}
}

func TestAllowedValuesValidation(t *testing.T) {
	field := &Field{Type: StringField, AllowedValues: []string{"a", "b"}}
	err := field.Validate("c")
	if err == nil {
		t.Errorf("Expected to get string validation error (not allowed value), but got nil")
	}
	err = field.Validate("a")
	if err != nil {
		t.Errorf("Expected to get no string validation error, but got %v", err)
	}

}

func TestNumFieldValidate(t *testing.T) {
	field := &Field{Type: NumField}
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

func TestArrayFieldValidate(t *testing.T) {
	field := &Field{Type: ArrayField}
	err := field.Validate([]int{1, 2, 3})
	if err != nil {
		t.Errorf("Expected to get no array validation error, but got %v", err)
	}
	err = field.Validate("string field")
	if err == nil {
		t.Errorf("Expected to get array validation error (string supplied), but got nil")
	}
	field.ArrayValues = &Field{Type: NumField}
	err = field.Validate([]int{1, 2, 3})
	if err != nil {
		t.Errorf("Expected to get no array validation error, but got %v", err)
	}
	err = field.Validate([]string{"a", "b", "c"})
	if err == nil {
		t.Errorf("Expected to get array validation error (string, but expected integer values), but got nil")
	}
}

func TestObjectFieldValidate(t *testing.T) {
	jsonString := `
{"number": 5, "object": {"key": "value"}}
`
	var params map[string]interface{}
	json.Unmarshal([]byte(jsonString), &params)
	field := &Field{Type: ObjectField}
	err := field.Validate(params["object"])
	if err != nil {
		t.Errorf("Expected to get no object validation error, but got %v", err)
	}
	err = field.Validate(params["number"])
	if err == nil {
		t.Errorf("Expected to get object validation error (int supplied), but got nil")
	}
	field.ObjectFields = []*Field{&Field{Name: "numbers", Type: NumField}}
	err = field.Validate(params)
	if err == nil {
		t.Errorf("Expected to get object validation error (wrong key supplied), but got nil")
	}
	field.ObjectFields = []*Field{
		&Field{Name: "number", Type: NumField},
		&Field{Name: "object", Type: ObjectField,
			ObjectFields: []*Field{
				&Field{Name: "key", Type: StringField},
			},
		},
	}
	err = field.Validate(params)
	if err != nil {
		t.Errorf("Expected to get no object validation error, but got %v", err)
	}
}

var boolCastTests = []struct {
	input  interface{}
	result interface{}
}{
	{"true", true},
	{"false", false},
	{1, true},
	{0, false},
	{"1", true},
	{"0", false},
	{true, true},
	{false, false},
}

func TestCastBoolean(t *testing.T) {
	field := &Field{Type: BoolField}
	for i, tt := range boolCastTests {
		value, ok := field.Cast(tt.input).(bool)
		if ok != true {
			t.Errorf("%d. boolean cast from '%v' was unsuccessful", i, tt.input)
		}
		if value != tt.result {
			t.Errorf("%d. boolean cast from '%v' to boolean is %v, but expected %v", i, tt.input, value, tt.result)
		}
	}
}

var numCastTests = []struct {
	input  interface{}
	result float64
}{
	{"2", 2},
	{"2.2", 2.2},
	{2, 2},
	{float32(2), 2},
}

func TestCastNumber(t *testing.T) {
	field := &Field{Type: NumField}
	eps := 1e-6
	for i, tt := range numCastTests {
		value, ok := field.Cast(tt.input).(float64)
		if ok != true {
			t.Errorf("%d. number cast from '%v' was unsuccessful", i, tt.input)
		}
		if !(math.Abs(value-tt.result) < eps) {
			t.Errorf("%d. number cast from '%v' to number is %v, but expected %v", i, tt.input, value, tt.result)
		}
	}
}
