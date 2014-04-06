package component

import (
	"github.com/romanoff/ahc/schema"
	"testing"
	"text/template"
)

func TestGetComponent(t *testing.T) {
	button1 := &Component{Namespace: "mp.a-button"}
	button2 := &Component{Namespace: "goog.a-button"}
	pool := &Pool{Components: []*Component{button1, button2}}
	if pool.GetComponent("a-button") != button1 {
		t.Errorf("Expected to get mp.a-button from pool, but got %v", pool.GetComponent("a-button").Namespace)
	}
	if pool.GetComponent("goog.a-button") != button2 {
		t.Errorf("Expected to get goog.a-button from pool, but got %v", pool.GetComponent("a-button").Namespace)
	}
}

func TestPoolRender(t *testing.T) {
	tmpl := template.Must(template.New("button").
		Parse("<div class='button'>{{.name}}</div>"))
	c := &Component{Namespace: "goog.a-button", Template: tmpl}
	pool := &Pool{Components: []*Component{c}}
	ahcx := `
<a-button name="Click me" />
<a-button name="Click me!" />
`
	html, err := pool.Render([]byte(ahcx))
	if err != nil {
		t.Errorf("Expected to get no error while rendering template, but got %v", err)
	}
	expected := `<div class="button">Click me</div><div class="button">Click me!</div>`
	if string(html) != expected {
		t.Errorf("Expected to get:\n%v\n, but got:\n%v\n", expected, string(html))
	}
}

func TestPoolRenderUsingNamespaceParams(t *testing.T) {
	tmpl := template.Must(template.New("button").
		Parse("<div class='button'>{{.name}}</div>"))
	c := &Component{Namespace: "goog.a-button", Template: tmpl}
	pool := &Pool{Components: []*Component{c}}
	ahcx := `
<a-button><a-button:name>Click me<a-button><a-button:name>Click me</a-button:name></a-button></a-button:name></a-button>
`
	html, err := pool.Render([]byte(ahcx))
	if err != nil {
		t.Errorf("Expected to get no error while rendering template, but got %v", err)
	}
	expected := `<div class="button">Click me<div class="button">Click me</div></div>`
	if string(html) != expected {
		t.Errorf("Expected to get:\n%v\n, but got:\n%v\n", expected, string(html))
	}
}

func TestPoolRenderWithDefaultParams(t *testing.T) {
	tmpl := template.Must(template.New("button").
		Parse("<div class='button'>{{.name}}</div>"))
	c := &Component{Namespace: "goog.a-button", Template: tmpl, DefaultParam: "name"}
	pool := &Pool{Components: []*Component{c}}
	ahcx := `
<a-button>Click me</a-button>
`
	html, err := pool.Render([]byte(ahcx))
	if err != nil {
		t.Errorf("Expected to get no error while rendering template, but got %v", err)
	}
	expected := `<div class="button">Click me</div>`
	if string(html) != expected {
		t.Errorf("Expected to get:\n%v\n, but got:\n%v\n", expected, string(html))
	}
}

func TestPoolRenderSafe(t *testing.T) {
	tmpl := template.Must(template.New("button").
		Parse("<div class='button'>{{.name}}</div>{{or .not_in_schema \"\"}}"))
	c := &Component{
		Namespace: "goog.a-button",
		Template:  tmpl,
		Schema: &schema.Schema{Fields: []*schema.Field{
			&schema.Field{Name: "name", Required: true, Type: schema.StringField},
		}},
	}
	pool := &Pool{Components: []*Component{c}, Safe: true}
	ahcx := `
<a-button></a-button>
`
	_, err := pool.Render([]byte(ahcx))
	if err == nil {
		t.Errorf("Expected error while rendering pool component with not enough params, but got nil")
	}
}

func TestPoolRenderContainerWithDefaultParams(t *testing.T) {
	tmpl := template.Must(template.New("button").
		Parse("<div class='button'>{{.name}}</div>"))
	c := &Component{Namespace: "goog.a-button", Template: tmpl}
	containerTmpl := template.Must(template.New("container").
		Parse("<div class='container'>{{.content}}</div>"))
	c1 := &Component{Namespace: "goog.a-container", Template: containerTmpl, DefaultParam: "content"}
	pool := &Pool{Components: []*Component{c, c1}}
	ahcx := `
<a-container><a-button name="Click me"/></a-container>
`

	html, err := pool.Render([]byte(ahcx))
	if err != nil {
		t.Errorf("Expected to get no error while rendering template, but got %v", err)
	}
	expected := `<div class="container"><div class="button">Click me</div></div>`
	if string(html) != expected {
		t.Errorf("Expected to get:\n%v\n, but got:\n%v\n", expected, string(html))
	}
}
