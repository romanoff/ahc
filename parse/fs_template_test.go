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
