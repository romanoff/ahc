package parse

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/romanoff/ahc/component"
	"github.com/romanoff/ahc/schema"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
)

type Fs struct {
}

// Parses component when path to css or html file is provided
func (self *Fs) ParseComponent(filepath string) (*component.Component, error) {
	if _, err := os.Stat(filepath); err != nil {
		return nil, errors.New(fmt.Sprintf("Error whie parsing component: %v file doesn't exist", filepath))
	}
	// Get base path (/a/b.css -> /a/b)
	filename := path.Base(filepath)
	filename = strings.TrimSuffix(filename, path.Ext(filepath))
	basePath := path.Dir(filepath) + "/" + filename
	component := &component.Component{}
	err := self.readAll(component, basePath)
	if err != nil {
		return nil, err
	}
	return component, nil
}

func (self *Fs) readAll(component *component.Component, basePath string) error {
	err := self.readCss(component, basePath)
	if err != nil {
		return err
	}
	err = self.readTemplate(component, basePath)
	if err != nil {
		return err
	}
	err = self.readSchema(component, basePath)
	if err != nil {
		return err
	}
	return nil
}

var provideRe *regexp.Regexp = regexp.MustCompile("@provide\\s+['\"](.+)['\"]")
var defaultParamRe *regexp.Regexp = regexp.MustCompile("@default_param\\s+['\"](.+)['\"]")
var requireRe *regexp.Regexp = regexp.MustCompile("@require\\s+['\"](.+)['\"]")

func (self *Fs) readCss(component *component.Component, basePath string) error {
	filepath := basePath + ".css"
	if _, err := os.Stat(filepath); err != nil {
		return nil
	}
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return errors.New(fmt.Sprintf("Error while reading css file: %v", filepath))
	}
	component.Css = string(content)
	matches := provideRe.FindSubmatch(content)
	if len(matches) == 2 {
		component.Namespace = string(matches[1])
	}
	matches = defaultParamRe.FindSubmatch(content)
	if len(matches) == 2 {
		component.DefaultParam = string(matches[1])
	}
	allMatches := requireRe.FindAllSubmatch(content, -1)
	for _, matches := range allMatches {
		if len(matches) == 2 {
			component.Requires = append(component.Requires, string(matches[1]))
		}
	}
	return nil
}

func (self *Fs) readTemplate(c *component.Component, basePath string) error {
	filepath := basePath + ".tmpl"
	if _, err := os.Stat(filepath); err != nil {
		return nil
	}
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return errors.New(fmt.Sprintf("Error while reading tmpl file: %v", filepath))
	}
	template := &component.Template{Content: string(content)}
	c.Template = template
	return nil
}

var fieldRe *regexp.Regexp = regexp.MustCompile("\\s*([\\w|-|_]+)\\s*\\{(\\w+)(=)?\\}({[^}]+})?\\s*(.*)\\s*")

func (self *Fs) readSchema(c *component.Component, basePath string) error {
	filepath := basePath + ".schema"
	if _, err := os.Stat(filepath); err != nil {
		return nil
	}
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return errors.New(fmt.Sprintf("Error while reading schema file: %v", filepath))
	}
	fields := make([]*schema.Field, 0, 0)
	lines := bytes.Split(content, []byte("\n"))
	for _, line := range lines {
		field, err := self.parseSchemaField(line)
		if err != nil && len(bytes.TrimSpace(line)) != 0 {
			return errors.New(fmt.Sprintf("Couldn't parse field: '%s' for schema in following file: '%v'", line, filepath))
		}
		if err == nil {
			fields = append(fields, field)
		}
	}
	schema := &schema.Schema{Fields: fields}
	c.Schema = schema
	return nil
}

func (self *Fs) parseSchemaField(fieldContent []byte) (*schema.Field, error) {
	matches := fieldRe.FindSubmatch(fieldContent)
	if len(matches) == 6 {
		field := &schema.Field{
			Name:        string(matches[1]),
			Description: string(matches[5]),
			Required:    !(string(matches[3]) == "="),
		}
		fieldType := string(matches[2])
		switch {
		case fieldType == "string":
			field.Type = schema.StringField
		case fieldType == "num":
			field.Type = schema.NumField
		case fieldType == "bool":
			field.Type = schema.BoolField
		default:
			return nil, errors.New(fmt.Sprintf("Undefined field type: %s", fieldType))
		}
		if len(matches[4]) > 0 {
			allowedValues := string(matches[4][1 : len(matches[4])-1])
			field.AllowedValues = strings.Split(allowedValues, "|")
		}
		return field, nil
	}
	return nil, errors.New(fmt.Sprintf("Field '%s' could not be parsed", fieldContent))
}
