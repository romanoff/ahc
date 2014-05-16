package view

import (
	"testing"
)

func TestGetTemplate(t *testing.T) {
	template := &Template{}
	templates := make(map[string]*Template)
	templates["/user"] = template
	pool := &Pool{Templates: templates}
	tmpl := pool.GetTemplate("/about")
	if tmpl != nil {
		t.Error("Expected not to get template from the pool, but got one")
	}
	tmpl = pool.GetTemplate("/user")
	if tmpl == nil {
		t.Error("Expected to get template from the pool, but got nil")
	}
}

func TestAddTemplate(t *testing.T) {
	template := &Template{Path: "/users", Content: "Username: {{.name}}"}
	pool := InitPool()
	pool.AddTemplate(template)

	params := make(map[string]interface{})
	params["name"] = "Jimmy"
	content, err := pool.Render("/users", params)
	if err != nil {
		t.Errorf("Not expected to get error while rendering templates pool template, but got %v", err)
	}
	expected := "Username: Jimmy"
	if expected != string(content) {
		t.Errorf("Expected to get:\n%v\n, but got:\n%v", expected, string(content))
	}
}
