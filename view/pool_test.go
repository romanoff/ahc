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
