package component

import (
	"github.com/romanoff/ahc/schema"
	"testing"
	"text/template"
)

func TestComponentTemplate(t *testing.T) {
	tmpl := template.Must(template.New("button").
		Parse("<div class='button'>{{.name}}</div>"))
	c := &Component{Namespace: "goog.a-button", Template: tmpl}
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

func TestRenderSafeNoValidation(t *testing.T) {
	tmpl := template.Must(template.New("a-button").
		Parse("<div class='button'>{{.name}}</div>"))
	c := &Component{Namespace: "goog.a-button", Template: tmpl}
	params := make(map[string]interface{})
	params["name"] = "Click me"
	_, err := c.RenderSafe(params)
	if err == nil {
		t.Errorf("Expected to get missing schema error for RenderSafe, but got nil")
	}
}

func TestRenderSafeValidation(t *testing.T) {
	tmpl := template.Must(template.New("button").
		Parse("<div class='button'>{{.name}}</div>{{or .not_in_schema \"\"}}"))
	c := &Component{
		Namespace: "goog.a-button",
		Template:  tmpl,
		Schema: &schema.Schema{Fields: []*schema.Field{
			&schema.Field{Name: "name", Required: true, Type: schema.StringField},
		}},
	}
	params := make(map[string]interface{})
	_, err := c.RenderSafe(params)
	if err == nil {
		t.Errorf("Expected to get missing schema parameter error for RenderSafe, but got nil")
	}
	params["name"] = "Click me"
	params["not_in_schema"] = "This text should not appear in button"
	html, err := c.RenderSafe(params)
	if err != nil {
		t.Errorf("Expected no error while rendering, but got %v", err)
	}
	expected := "<div class='button'>Click me</div>"
	if string(html) != expected {
		t.Errorf("Expected to get:\n%v\n, but got:\n%v", expected, string(html))
	}
}
