package component

import (
	"bytes"
	"text/template"
)

// Ahc component
type Component struct {
	Namespace string
	Requires  []string
	Css       string
	Template  *template.Template
}

// Renders component with data
func (self *Component) Render(params map[string]interface{}) ([]byte, error) {
	out := bytes.Buffer{}
	err := self.Template.Execute(&out, params)
	return out.Bytes(), err
}
