package component

import (
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
