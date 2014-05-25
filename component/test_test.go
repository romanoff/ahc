package component

import (
	"github.com/romanoff/ahc/schema"
	"github.com/romanoff/htmlcompressor"
	"testing"
)

func TestComponentTest(t *testing.T) {
	params := make(map[string]interface{})
	params["name"] = "Button"
	expected := `
<div class="ahc_button">Button</div>
`
	test := &Test{Params: params, Expected: []byte(expected)}
	tmpl := `<div class="ahc_button">{{.name|html}}</div>`
	c := &Component{Namespace: "goog.a-button", Template: &Template{Content: tmpl}, Schema: &schema.Schema{}}
	compressor := htmlcompressor.InitAll()
	err := test.Run(c, nil, compressor)
	if err != nil {
		t.Errorf("Expected to get no test error, but got %v", err)
	}
	expected = `
	<div class="ahc_button"> Button </div>
	`
	test.Expected = []byte(expected)
	err = test.Run(c, nil, compressor)
	if err == nil {
		t.Errorf("Expected to get test error, but got nil")
	}
}
