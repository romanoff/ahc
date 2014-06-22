package view

import (
	"github.com/romanoff/ahc/component"
	"testing"
)

func TestInitView(t *testing.T) {
	viewContent := `
<div class="container">
  <a-button name="Fancy button"/>
</div>
`
	view, err := InitView([]byte(viewContent))
	if err != nil {
		t.Errorf("Expected no error while initializing view, but got %v", err)
	}

	tmpl := "<div class='button'>{{.name}}</div>"
	c := &component.Component{Namespace: "goog.a-button", Template: &component.Template{Content: tmpl}}
	pool := &component.Pool{Components: []*component.Component{c}}

	content, err := view.GetContent(&RenderParams{Pool: pool})
	if err != nil {
		t.Errorf("Expected no error while getting view content, but got %v", err)
	}
	expected := `<div class="container"><div class="button">Fancy button</div></div>`
	if expected != string(content) {
		t.Errorf("Expected to get:\n%v\n, but got:\n%v", expected, string(content))
	}
}

func TestGetContent(t *testing.T) {
	viewContent := `
<a-button>Fancy name</a-button>
`
	view, err := InitView([]byte(viewContent))
	if err != nil {
		t.Errorf("Expected no error while initializing view, but got %v", err)
	}
	tmpl := "<div class='button'>{{.name}}</div>"
	c := &component.Component{Namespace: "goog.a-button", Template: &component.Template{Content: tmpl}, DefaultParam: "name"}
	pool := &component.Pool{Components: []*component.Component{c}}

	content, err := view.GetContent(&RenderParams{Pool: pool})
	if err != nil {
		t.Errorf("Expected no error while getting view content, but got %v", err)
	}
	expected := `<div class="button">Fancy name</div>`
	if expected != string(content) {
		t.Errorf("Expected to get:\n%v\n, but got:\n%v", expected, string(content))
	}
}
