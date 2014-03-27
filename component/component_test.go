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
	html, err := c.Render(params, nil)
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
	_, err := c.RenderSafe(params, nil)
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
	_, err := c.RenderSafe(params, nil)
	if err == nil {
		t.Errorf("Expected to get missing schema parameter error for RenderSafe, but got nil")
	}
	params["name"] = "Click me"
	params["not_in_schema"] = "This text should not appear in button"
	html, err := c.RenderSafe(params, nil)
	if err != nil {
		t.Errorf("Expected no error while rendering, but got %v", err)
	}
	expected := "<div class='button'>Click me</div>"
	if string(html) != expected {
		t.Errorf("Expected to get:\n%v\n, but got:\n%v", expected, string(html))
	}
}

func TestRender(t *testing.T) {
	tmpl := template.Must(template.New("button").
		Parse("<div class='button'>{{.name}}</div>"))
	tmpl1 := template.Must(template.New("multibutton").
		Parse("<div class='multibutton'><a-button name='one'/><a-button name='two'/></div>"))
	c := &Component{Namespace: "goog.a-button", Template: tmpl}
	multibutton := &Component{Namespace: "goog.a-multibutton", Template: tmpl1}
	pool := &Pool{Components: []*Component{c, multibutton}}
	params := make(map[string]interface{})
	html, err := multibutton.Render(params, pool)
	if err != nil {
		t.Errorf("Expected not to get error while rendering complex component, but got %v", err)
	}
	expected := `<div class="multibutton"><div class="button">one</div><div class="button">two</div></div>`
	if expected != string(html) {
		t.Errorf("Expected to get:\n%v\n, but got:\n%v", expected, string(html))
	}

	tmpl2 := template.Must(template.New("multibutton").
		Parse("<div class='multibutton'><a-button name='one'/><a-button name='two'/></div><img src='image.png' />"))
	multibutton.Template = tmpl2
	html, err = multibutton.Render(params, pool)
	if err != nil {
		t.Errorf("Expected not to get error while rendering complex component, but got %v", err)
	}
	expected = `<div class="multibutton"><div class="button">one</div><div class="button">two</div></div><img src="image.png" />`
	if expected != string(html) {
		t.Errorf("Expected to get:\n%v\n, but got:\n%v", expected, string(html))
	}

}
