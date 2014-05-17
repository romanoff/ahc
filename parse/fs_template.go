package parse

import (
	"errors"
	"fmt"
	"github.com/romanoff/ahc/view"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func (self *Fs) ParseTemplate(filepath string, basepath string) (*view.Template, error) {
	if basepath == "" {
		return nil, errors.New(fmt.Sprintf("Basepath not specified while parsing template: '%v'", filepath))
	}
	if _, err := os.Stat(filepath); err != nil {
		return nil, errors.New(fmt.Sprintf("Error whie parsing component: %v file doesn't exist", filepath))
	}

	filename := path.Base(filepath)
	filename = strings.TrimSuffix(filename, path.Ext(filepath))
	basePath := path.Dir(filepath) + "/" + filename
	schemaPath := basePath + ".schema"
	//schemaContent
	_, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error while reading template '%v' schema: %v", filepath, err))
	}
	// Parse schema content
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error while reading template '%v': %v", filepath, err))
	}
	template := &view.Template{Content: string(content)}
	return template, nil
}
