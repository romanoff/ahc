package view

import (
	"html/template"
)

type Template struct {
	Content []byte
}

func (self *Template) Render(params map[string]interface{}, pool *component.Pool) ([]byte, error) {
	return nil, nil
}
