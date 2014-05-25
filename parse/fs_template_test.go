package parse

import (
	"testing"
)

func TestParseTemplate(t *testing.T) {
	fs := &Fs{}
	tmpl, err := fs.ParseTemplate("test_templates/index.html", "test_templates")
	if err != nil {
		t.Errorf("Expected to parse index template, but got errpr: '%v'", err)
	}
	if tmpl.Path != "index" {
		t.Errorf("Expected to get 'index' as template path, but got '%v'", tmpl.Path)
	}
	if tmpl.Schema == nil {
		t.Error("Expected to get 'index' template schema, but got nil")
	}
}

func TestParseTemplateJson(t *testing.T) {
	fs := &Fs{}
	templateJson, err := fs.ParseTemplateJson("test_templates/index.json")
	if err != nil {
		t.Errorf("Expected to get no error while parsing index template json, but got: '%v'", err)
	}
	expected := 2
	if len(templateJson.JsonGroups) != expected {
		t.Errorf("Expected to get:\n%v\n, but got:\n%v", expected, len(templateJson.JsonGroups))
	}
}
