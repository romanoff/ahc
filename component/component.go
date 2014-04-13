package component

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jteeuwen/go-pkg-xmlx"
	"github.com/romanoff/ahc/schema"
)

// Ahc component
type Component struct {
	Namespace    string
	Requires     []string
	DefaultParam string
	Css          string
	Template     *Template
	Schema       *schema.Schema
}

// Renders component with params
func (self *Component) Render(params map[string]interface{}, pool *Pool) ([]byte, error) {
	params = self.CastParams(params)
	out := bytes.Buffer{}
	err := self.Template.Execute(&out, params)
	if err != nil {
		return nil, err
	}
	if pool == nil {
		return out.Bytes(), err
	}
	document := xmlx.New()
	document.LoadBytes(out.Bytes(), nil)
	return pool.getNodesHtml(document.Root.Children)
}

// Renders component with params and verifies params schema
func (self *Component) RenderSafe(params map[string]interface{}, pool *Pool) ([]byte, error) {
	if self.Schema == nil {
		return nil, errors.New(fmt.Sprintf("No schema provided for %v", self.Namespace))
	}
	schemaErrors := self.Schema.Validate(params)
	if len(schemaErrors) != 0 {
		errorString := ""
		for _, err := range schemaErrors {
			errorString += fmt.Sprintf("%v\n", err)
		}
		return nil, errors.New(errorString)
	}
	filteredParams := self.Schema.GetSchemaParams(params)
	return self.Render(filteredParams, nil)
}

func (self *Component) CastParams(params map[string]interface{}) map[string]interface{} {
	if self.Schema == nil {
		return params
	}
	for _, field := range self.Schema.Fields {
		if params[field.Name] != nil {
			params[field.Name] = field.Cast(params[field.Name])
		}
	}
	return params
}

// Returns component css after using preprocessor
func (self *Component) GetCss(pool *Pool) (string, error) {
	return self.Css, nil
}
