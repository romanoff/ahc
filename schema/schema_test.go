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

func TestGetSchemaParams(t *testing.T) {
	params := make(map[string]interface{})
	params["name"] = "John"
	params["count"] = 5
	schema := &Schema{Fields: []*Field{
		&Field{Name: "name", Type: StringField},
	}}
	schemaParams := schema.GetSchemaParams(params)
	if schemaParams["count"] != nil {
		t.Errorf("Expected not to get count parameter as it's not specified in schema, but got %v", schemaParams["count"])
	}
	nameValue, _ := schemaParams["name"].(string)
	if nameValue != "John" {
		t.Errorf("Expected to get John as name parameter as it's specified in schema, but got %v", nameValue)
	}
}
