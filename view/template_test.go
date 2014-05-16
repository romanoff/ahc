package view

import (
	"github.com/romanoff/ahc/component"
	"testing"
)

func TestTemplateRender(t *testing.T) {
	tmpl := "<div class='button'>{{.name|html}}</div>"
	c := &component.Component{Namespace: "goog.a-button", Template: &component.Template{Content: tmpl}}
	pool := &component.Pool{Components: []*component.Component{c}}

	tmplContent := `
<div class="container">
  {{ range .buttons }}
  <a-button name="{{.}}"/>
  {{ end }}
</div>
`
	template := &Template{Content: tmplContent}
	params := make(map[string]interface{})
	params["buttons"] = []string{"A", "B", "C"}
	content, err := template.Render(params, pool)
	if err != nil {
		t.Errorf("Not expected to get an error while rendering template, but got %v", err)
	}
	expected := `<div class="container"><div class="button">A</div><div class="button">B</div><div class="button">C</div></div>`
	if expected != string(content) {
		t.Errorf("Expected to get:\n%v\n, but got:\n%v", expected, string(content))
	}
}
