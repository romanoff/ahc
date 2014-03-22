package component

import (
	"testing"
	"text/template"
)

func TestComponentTemplate(t *testing.T) {
	tmpl := template.Must(template.New("button").
		Parse("<div class='button'>{{.name}}</div>"))
	c := &Component{Namespace: "goog.button", Template: tmpl}
	params := make(map[string]interface{})
	params["name"] = "Click me"
	html, err := c.Render(params)
	if err != nil {
		t.Errorf("Expected no error while rendering, but got %v", err)
	}
	expected := "<div class='button'>Click me</div>"
	if string(html) != expected {
		t.Errorf("Expected to get:\n%v\n, but got:\n%v", expected, string(html))
	}
}
