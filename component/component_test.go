package component

import (
	"github.com/romanoff/ahc/schema"
	"testing"
)

func TestComponentTemplate(t *testing.T) {
	tmpl := `<div class='button'>{{.name}}</div>`
	c := &Component{Namespace: "goog.a-button", Template: &Template{Content: tmpl}}
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
	tmpl := `<div class='button'>{{.name}}</div>`
	c := &Component{Namespace: "goog.a-button", Template: &Template{Content: tmpl}}
	params := make(map[string]interface{})
	params["name"] = "Click me"
	_, err := c.RenderSafe(params, nil)
	if err == nil {
		t.Errorf("Expected to get missing schema error for RenderSafe, but got nil")
	}
}

func TestRenderSafeValidation(t *testing.T) {
	tmpl := "<div class='button'>{{.name}}</div>{{or .not_in_schema \"\"}}"
	c := &Component{
		Namespace: "goog.a-button",
		Template:  &Template{Content: tmpl},
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
	tmpl := "<div class='button'>{{.name}}</div>"
	tmpl1 := "<div class='multibutton'><a-button name='one'/><a-button name='two'/></div>"
	c := &Component{Namespace: "goog.a-button", Template: &Template{Content: tmpl}}
	multibutton := &Component{Namespace: "goog.a-multibutton", Template: &Template{Content: tmpl1}}
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

	tmpl2 := "<div class='multibutton'><a-button name='one'/><a-button name='two'/></div><img src='image.png' />"
	multibutton.Template = &Template{Content: tmpl2}
	html, err = multibutton.Render(params, pool)
	if err != nil {
		t.Errorf("Expected not to get error while rendering complex component, but got %v", err)
	}
	expected = `<div class="multibutton"><div class="button">one</div><div class="button">two</div></div><img src="image.png" />`
	if expected != string(html) {
		t.Errorf("Expected to get:\n%v\n, but got:\n%v", expected, string(html))
	}

}

func TestCastParams(t *testing.T) {
	tmpl := "<div class='button' {{ if .hidden}}style='display: none;'{{end}}>{{.name}}</div>"
	c := &Component{
		Namespace: "goog.a-button",
		Template:  &Template{Content: tmpl},
		Schema: &schema.Schema{Fields: []*schema.Field{
			&schema.Field{Name: "name", Required: true, Type: schema.StringField},
			&schema.Field{Name: "hidden", Type: schema.BoolField},
		}},
	}
	params := make(map[string]interface{})
	params["name"] = "Here goes name"
	params["hidden"] = "false"
	html, err := c.Render(params, nil)
	if err != nil {
		t.Errorf("Expected not to get error while rendering button, but got %v", err)
	}
	expected := `<div class='button' >Here goes name</div>`
	if expected != string(html) {
		t.Errorf("Expected to get:\n%v\n, but got:\n%v", expected, string(html))
	}
}

func TestGetCss(t *testing.T) {
	c := &Component{
		Namespace: "goog.a-button",
		Css:       ".goog_button { color: red; }",
	}
	css, err := c.GetCss(nil)
	if err != nil {
		t.Errorf("Expected to not get error while extracting css, but got %v", err)
	}
	expected := ".goog_button { color: red; }"
	if expected != css {
		t.Errorf("Expected to get:\n%v\n, but got:\n%v", expected, css)
	}
}

func TestGetClassWithRequires(t *testing.T) {
	c1 := &Component{Namespace: "goog.a-button", Css: ".button {}"}
	c2 := &Component{Namespace: "goog.a-multibutton", Css: ".multibutton {}", Requires: []string{"goog.a-button"}}
	pool := &Pool{Components: []*Component{c1, c2}}
	multibuttonCss, err := c2.GetCss(pool)
	if err != nil {
		t.Errorf("Expected to not get error while extracting css, but got %v", err)
	}
	expected := `.button {}
.multibutton {}`
	if expected != multibuttonCss {
		t.Errorf("Expected to get:\n%v\n, but got:\n%v", expected, multibuttonCss)
	}
}

type TestPreprocessor struct {
}

func (self *TestPreprocessor) GetCss(input []byte) ([]byte, error) {
	returnValue := []byte("/*Added by preprocessor*/\n")
	returnValue = append(returnValue, input...)
	return returnValue, nil
}

func TestGetClassWithPreprocessor(t *testing.T) {
	c1 := &Component{Namespace: "goog.a-button", Css: ".button {}"}
	c2 := &Component{Namespace: "goog.a-multibutton", Css: ".multibutton {}", Requires: []string{"goog.a-button"}}
	preprocessor := &TestPreprocessor{}
	pool := &Pool{Components: []*Component{c1, c2}, Preprocessor: preprocessor}
	multibuttonCss, err := c2.GetCss(pool)
	if err != nil {
		t.Errorf("Expected to not get error while extracting css, but got %v", err)
	}
	expected := `/*Added by preprocessor*/
.button {}
.multibutton {}`
	if expected != multibuttonCss {
		t.Errorf("Expected to get:\n%v\n, but got:\n%v", expected, multibuttonCss)
	}
}
