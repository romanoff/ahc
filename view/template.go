package view

import (
	"bytes"
	"github.com/romanoff/ahc/component"
	"html/template"
)

type Template struct {
	Content  string
	compiled bool
	Template *template.Template
}

func (self *Template) compileTemplate() error {
	compiledTemplate, err := template.New("template").Parse(self.Content)
	if err != nil {
		return err
	}
	self.Template = compiledTemplate
	return nil
}

func (self *Template) render(params map[string]interface{}, pool *component.Pool, safe bool) ([]byte, error) {
	if !self.compiled {
		err := self.compileTemplate()
		if err != nil {
			return nil, err
		}
	}
	out := bytes.Buffer{}
	err := self.Template.Execute(&out, params)
	if err != nil {
		return nil, err
	}
	view, err := InitView(out.Bytes())
	if err != nil {
		return nil, err
	}
	return view.GetContent(&RenderParams{Pool: pool, Safe: safe})
}

func (self *Template) Render(params map[string]interface{}, pool *component.Pool) ([]byte, error) {
	return self.render(params, pool, false)
}

func (self *Template) RenderSafe(params map[string]interface{}, pool *component.Pool) ([]byte, error) {
	return self.render(params, pool, true)
}
