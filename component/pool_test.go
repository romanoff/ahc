package component

import (
	"testing"
	"text/template"
)

func TestGetComponent(t *testing.T) {
	button1 := &Component{Namespace: "mp.button"}
	button2 := &Component{Namespace: "goog.button"}
	pool := &Pool{Components: []*Component{button1, button2}}
	if pool.GetComponent("button") != button1 {
		t.Errorf("Expected to get mp.button from pool, but got %v", pool.GetComponent("button").Namespace)
	}
	if pool.GetComponent("goog.button") != button2 {
		t.Errorf("Expected to get goog.button from pool, but got %v", pool.GetComponent("button").Namespace)
	}
}

func TestPoolRender(t *testing.T) {
	tmpl := template.Must(template.New("button").
		Parse("<div class='button'>{{.name}}</div>"))
	c := &Component{Namespace: "goog.button", Template: tmpl}
	pool := &Pool{Components: []*Component{c}}
	ahcx := `
<button name="Click me" />
<button name="Click me!" />
`
	html, err := pool.Render([]byte(ahcx))
	if err != nil {
		t.Errorf("Expected to get no error while rendering template, but got %v", err)
	}
	expected := "<div class='button'>Click me</div><div class='button'>Click me!</div>"
	if string(html) != expected {
		t.Errorf("Expected to get:\n%v\n, but got:\n%v\n", expected, string(html))
	}
}
