package view

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/romanoff/ahc/component"
	"html/template"
)

func InitPool() *Pool {
	return &Pool{Templates: make(map[string]*Template), templates: template.New("")}
}

type Pool struct {
	Templates      map[string]*Template
	templates      *template.Template
	Pools          []*Pool
	componentsPool *component.Pool
}

func (self *Pool) AddTemplate(tmpl *Template) error {
	if tmpl.Path == "" {
		return errors.New("Cannot add template without path to the pool")
	}
	self.Templates[tmpl.Path] = tmpl
	_, err := self.templates.Parse("{{define \"" + tmpl.Path + "\"}}" + tmpl.Content + "{{end}}")
	if err != nil {
		return errors.New(fmt.Sprintf("Error while parsing '%v' template: %v", tmpl.Path, err))
	}
	return nil
}

func (self *Pool) render(path string, params map[string]interface{}, safe bool) ([]byte, error) {
	template := self.GetTemplate(path)
	if template == nil {
		return nil, errors.New(fmt.Sprintf("Cannot find template '%v'", path))
	}
	if template.Schema != nil && safe {
		schemaErrors := template.Schema.Validate(params)
		if len(schemaErrors) != 0 {
			errorString := ""
			for _, err := range schemaErrors {
				errorString += fmt.Sprintf("%v\n", err)
			}
			return nil, errors.New(errorString)
		}
	}
	out := bytes.Buffer{}
	err := self.templates.Lookup(path).Execute(&out, params)
	if err != nil {
		return nil, err
	}
	view, err := InitView(out.Bytes())
	if err != nil {
		return nil, err
	}
	return view.GetContent(&RenderParams{Pool: self.componentsPool, Safe: safe})
}

func (self *Pool) Render(path string, params map[string]interface{}) ([]byte, error) {
	return self.render(path, params, false)
}

func (self *Pool) RenderSafe(path string, params map[string]interface{}) ([]byte, error) {
	return self.render(path, params, true)
}

func (self *Pool) GetTemplate(path string) *Template {
	tmpl := self.Templates[path]
	if tmpl != nil {
		return tmpl
	}
	for _, pool := range self.Pools {
		tmpl = pool.GetTemplate(path)
		if tmpl != nil {
			return tmpl
		}
	}
	return nil
}
