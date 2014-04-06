package component

import (
	"io"
	"text/template"
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

func (self *Template) Execute(wr io.Writer, data interface{}) error {
	if !self.compiled {
		err := self.compileTemplate()
		if err != nil {
			return err
		}
	}
	return self.Template.Execute(wr, data)
}
