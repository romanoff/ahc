package parse

import (
	"errors"
	"fmt"
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
	return nil
}

var provideRe *regexp.Regexp = regexp.MustCompile("@provide\\s+['\"](.+)['\"]")
var defaultParamRe *regexp.Regexp = regexp.MustCompile("@default_param\\s+['\"](.+)['\"]")

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
	return nil
}
