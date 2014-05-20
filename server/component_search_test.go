package server

import (
	"github.com/romanoff/ahc/component"
	"github.com/romanoff/ahc/view"
	"testing"
)

func TestGetUsedNamespaces(t *testing.T) {
	template := &view.Template{Path: "/users", Content: "Username: {{.name}} <div><a-button name='some' /></div>"}
	pool := view.InitPool()
	pool.AddTemplate(template)

	tmpl := "<div class='button'>{{.name|html}}</div>"
	c := &component.Component{Namespace: "goog.a-button", Template: &component.Template{Content: tmpl}}
	componentsPool := &component.Pool{Components: []*component.Component{c}}
	pool.ComponentsPool = componentsPool

	search := InitComponentSearch(pool)
	search.Search("/users")
	expected := 1
	if len(search.Components) != expected {
		t.Errorf("Expected to find 1 components for users template, but got: %v", len(search.Components))
	}
}
