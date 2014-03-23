package schema

import (
	"testing"
)

func TestSchemaValidation(t *testing.T) {
	schema := &Schema{Fields: []*Field{
		&Field{Name: "name", Type: StringField},
		&Field{Name: "count", Type: NumField},
	}}
	params := make(map[string]interface{})
	params["name"] = "Hello"
	params["count"] = 20
	errors := schema.Validate(params)
	if len(errors) != 0 {
		t.Errorf("Expected to get 0 schema validation errors,  but got %v", len(errors))
	}
	params["count"] = "Hello"
	errors = schema.Validate(params)
	if len(errors) != 1 {
		t.Errorf("Expected to get 1 schema validation errors,  but got %v", len(errors))
	}
}
