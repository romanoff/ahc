package parse

import (
	"github.com/kr/pretty"
	"testing"
)

func TestNonExistingFileParse(t *testing.T) {
	fs := &Fs{}
	_, err := fs.ParseComponent("test_files/non_existing_component/non_existing_component.css")
	if err == nil {
		t.Error("Expected to get file doesn't exist error, but got nil")
	}
}

func TestParseComponent(t *testing.T) {
	fs := &Fs{}
	component, err := fs.ParseComponent("test_files/button/button.css")
	if err != nil {
		t.Errorf("Expected not to get error while parsing button component, but got %v", err)
	}
	if component.Css == "" {
		t.Errorf("Expected component css to not be empty after parsing button css")
	}
	if component.Namespace != "ahc.a-button" {
		t.Errorf("Expected button component namespace to be ahc.button, but got '%v'", component.Namespace)
	}
	if component.DefaultParam != "name" {
		t.Errorf("Expected button component default param to be name, but got '%v'", component.DefaultParam)
	}
	if len(component.Requires) != 1 || component.Requires[0] != "ahc.reset" {
		t.Errorf("Expected to get ahc.reset as require for button component, but got %v", component.Requires)
	}
	if component.Template == nil {
		t.Error("Expected to get button template, but got nil")
	}
	if component.Schema == nil || len(component.Schema.Fields) != 5 {
		t.Errorf("Expected to get schema with 5 fields, but got %# v", pretty.Formatter(component.Schema))
	}
}

func TestParseHtmlComponent(t *testing.T) {
	fs := &Fs{}
	component, err := fs.ParseComponent("test_files/menuitem/menuitem.html")
	if err != nil {
		t.Errorf("Expected not to get error while parsing menuitem component, but got %v", err)
	}
	if component.Css == "" {
		t.Errorf("Expected component css to not be empty after parsing menuitem css")
	}
	if component.Namespace != "ahc.menu-item" {
		t.Errorf("Expected menuitem component namespace to be ahc.menu-item, but got '%v'", component.Namespace)
	}
	if component.DefaultParam != "content" {
		t.Errorf("Expected menuitem component default param to content name, but got '%v'", component.DefaultParam)
	}
	if component.Template == nil {
		t.Error("Expected to get menuitem template, but got nil")
	}
	if component.Schema == nil || len(component.Schema.Fields) != 1 {
		t.Errorf("Expected to get schema with 1 fields, but got %# v", pretty.Formatter(component.Schema))
	}
}
