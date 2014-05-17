package parse

import (
	"testing"
)

func TestSchema(t *testing.T) {
	schemaContent := `
name {string} username
address {object} address object
  country {string} country name
  city {string} city name
people {array} Array of people
  name {string} person name
  last_name {string} person last name
`
	fs := &Fs{}
	schema, err := fs.ParseSchema([]byte(schemaContent), "schemapath - for display purposes only")
	if err != nil {
		t.Errorf("Did not expect to get error while parsing schema: '%v'", err)
	}
	if len(schema.Fields) != 3 {
		t.Errorf("Expected schema to have 3 fields, but got %v", len(schema.Fields))
	}
}
