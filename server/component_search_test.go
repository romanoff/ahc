package server

import (
	"github.com/romanoff/ahc/component"
	"github.com/romanoff/ahc/view"
	"testing"
)

func TestGetUsedNamespaces(t *testing.T) {
	template := &view.Template{Path: "users", Content: `Username: {{.name}} <div>{{ template "_header" }}</div>`}
	pool := view.InitPool()
	pool.AddTemplate(template)
	partial := &view.Template{Path: "_header", Content: `{{template "_header1"}}`}
	pool.AddTemplate(partial)
	partial1 := &view.Template{Path: "_header1", Content: `<div class="header1"><a-button name='some' /></div>`}
	pool.AddTemplate(partial1)

	tmpl := "<div class='button'>{{.name|html}}</div>"
	c := &component.Component{Namespace: "goog.a-button", Template: &component.Template{Content: tmpl}}
	componentsPool := &component.Pool{Components: []*component.Component{c}}
	pool.ComponentsPool = componentsPool

	search := InitComponentSearch(pool)
	search.Search("users")
	expected := 1
	if len(search.Components) != expected {
		t.Errorf("Expected to find 1 components for users template, but got: %v", len(search.Components))
	}
}

func TestGetUsedTemplates(t *testing.T) {
	search := InitComponentSearch(nil)
	templates := search.GetUsedTemplates([]byte(`{{template "some" .}} {{template "_header" }}`))
	if len(templates) != 2 {
		t.Errorf("Expected to get 2 templates, but got: %v", len(templates))
	}
}
