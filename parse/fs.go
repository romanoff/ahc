package parse

import (
	"errors"
	"fmt"
	"github.com/jteeuwen/go-pkg-xmlx"
	"github.com/romanoff/ahc/component"
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
	basePath := strings.TrimSuffix(filepath, path.Ext(filepath))
	component := &component.Component{}
	err := self.readAll(component, basePath)
	if err != nil {
		return nil, err
	}
	err = component.Validate()
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
	err = self.readHtml(component, basePath)
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

func (self *Fs) readSchema(c *component.Component, basePath string) error {
	filepath := basePath + ".schema"
	if _, err := os.Stat(filepath); err != nil {
		return nil
	}
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return errors.New(fmt.Sprintf("Error while reading schema file: %v", filepath))
	}
	schema, err := self.ParseSchema(content, filepath)
	if err != nil {
		return err
	}
	c.Schema = schema
	return nil
}

func (self *Fs) readHtml(c *component.Component, basePath string) error {
	filepath := basePath + ".html"
	if _, err := os.Stat(filepath); err != nil {
		return nil
	}
	content, err := ioutil.ReadFile(filepath)
	document := xmlx.New()
	document.LoadBytes(content, nil)
	html := document.Root.Children[0]
	for _, attr := range html.Attributes {
		if attr.Name.Local == "namespace" {
			c.Namespace = attr.Value
		}
		if attr.Name.Local == "default_param" {
			c.DefaultParam = attr.Value
		}
		if attr.Name.Local == "require" {
			c.Requires = []string{attr.Value}
		}
	}
	for _, node := range html.Children {
		if node.Type != xmlx.NT_ELEMENT {
			continue
		}
		if node.Name.Local == "style" {
			c.Css = getXmlNodesContent(node.Children)
		}
		if node.Name.Local == "schema" {
			content := getXmlNodesContent(node.Children)
			lines := strings.Split(content, "\n")
			schemaContent := []byte{}
			for _, line := range lines {
				schemaContent = append(schemaContent, []byte(strings.Replace(line, "    ", "", 1))...)
			}
			schema, err := self.ParseSchema(schemaContent, filepath)
			if err != nil {
				return err
			}
			c.Schema = schema
		}
		if node.Name.Local == "template" {
			c.Template = &component.Template{Content: getXmlNodesContent(node.Children)}
		}
	}
	if err != nil {
		return errors.New(fmt.Sprintf("Error while reading html file: %v", filepath))
	}
	return nil
}

func getXmlNodesContent(nodes []*xmlx.Node) string {
	content := ""
	for _, node := range nodes {
		content += strings.TrimSpace(node.String())
	}
	return content
}
