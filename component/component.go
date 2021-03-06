package component

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jteeuwen/go-pkg-xmlx"
	"github.com/romanoff/ahc/schema"
	"strings"
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

func (self *Component) Validate() error {
	if strings.TrimSpace(self.Namespace) == "" {
		return errors.New("component is missing namespace")
	}
	return nil
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

func (self *Component) GetRawCss(pool *Pool) (string, error) {
	cssContent := ""
	for _, namespace := range self.Requires {
		component := pool.GetComponent(namespace)
		if component == nil {
			return "", errors.New(fmt.Sprintf("Missing require: %v", namespace))
		}
		componentCss, err := component.GetRawCss(pool)
		if err != nil {
			return "", err
		}
		cssContent += componentCss + "\n"
	}
	cssContent += self.Css
	return cssContent, nil
}

// Returns component css after using preprocessor
func (self *Component) GetCss(pool *Pool) (string, error) {
	cssContent, err := self.GetRawCss(pool)
	if err != nil {
		return "", err
	}
	if pool != nil && pool.Preprocessor != nil {
		preprocessor := pool.Preprocessor
		cssContent, err := preprocessor.GetCss([]byte(cssContent))
		if err != nil {
			return "", err
		}
		return string(cssContent), nil
	}
	return cssContent, nil
}
