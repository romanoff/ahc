package parse

import (
	"errors"
	"fmt"
	"github.com/romanoff/ahc/view"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func (self *Fs) ParseTemplate(templatepath string, basepath string) (*view.Template, error) {
	if basepath == "" {
		return nil, errors.New(fmt.Sprintf("Basepath not specified while parsing template: '%v'", templatepath))
	}
	if _, err := os.Stat(templatepath); err != nil {
		return nil, errors.New(fmt.Sprintf("Error whie parsing component: %v file doesn't exist", templatepath))
	}

	basePath := strings.TrimSuffix(templatepath, path.Ext(templatepath))
	schemaPath := basePath + ".schema"
	//schemaContent
	content, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error while reading template '%v' schema: %v", templatepath, err))
	}
	schema, err := self.ParseSchema(content, schemaPath)
	if err != nil {
		return nil, err
	}

	content, err = ioutil.ReadFile(templatepath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error while reading template '%v': %v", templatepath, err))
	}
	absFilepath, err := filepath.Abs(templatepath)
	if err != nil {
		return nil, err
	}
	absFilepath = strings.TrimSuffix(absFilepath, path.Ext(templatepath))
	absBasepath, err := filepath.Abs(basepath)
	if err != nil {
		return nil, err
	}
	template := &view.Template{
		Content: string(content),
		Schema:  schema,
		Path:    strings.TrimPrefix(absFilepath, absBasepath),
	}
	return template, nil
}
