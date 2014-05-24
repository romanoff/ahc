package parse

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/romanoff/ahc/server"
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
		Path:    strings.TrimPrefix(strings.TrimPrefix(absFilepath, absBasepath), "/"),
	}
	return template, nil
}

func (self *Fs) ParseTemplateJson(jsonpath string) (*server.TemplateJson, error) {
	if _, err := os.Stat(jsonpath); err != nil {
		return nil, errors.New(fmt.Sprintf("Error whie parsing json params for template: %v file doesn't exist", jsonpath))
	}
	content, err := ioutil.ReadFile(jsonpath)
	if err != nil {
		return nil, err
	}
	lines := bytes.Split(content, []byte("\n"))
	name := []byte{}
	jsonContent := []byte{}
	nameParsing := true
	templateJson := &server.TemplateJson{JsonGroups: make([]*server.JsonGroup, 0, 0)}
	for _, line := range lines {
		delimiterLine := false
		if len(strings.TrimSpace(string(line))) == 0 {
			delimiterLine = true
		}
		if len(strings.TrimRight(strings.TrimSpace(string(line)), "-")) == 0 {
			delimiterLine = true
		}
		if delimiterLine {
			if nameParsing == false {
				jsonGroup, err := self.createJsonGroup(name, jsonContent)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("%v : %v", err, jsonpath))
				}
				templateJson.JsonGroups = append(templateJson.JsonGroups, jsonGroup)
				name = []byte{}
				jsonContent = []byte{}
				nameParsing = true
			} else {
				nameParsing = false
			}
			continue
		}
		if nameParsing {
			if len(name) != 0 {
				name = append(name, []byte("\n")...)
			}
			name = append(name, line...)
		} else {
			jsonContent = append(jsonContent, []byte(" ")...)
			jsonContent = append(jsonContent, line...)
		}
	}
	return templateJson, nil
}

func (self *Fs) createJsonGroup(name, jsonContent []byte) (*server.JsonGroup, error) {
	//TODO: Recover from invalid json parsing
	if len(strings.TrimSpace(string(name))) == 0 {
		return nil, errors.New("json group name is missing")
	}
	if len(strings.TrimSpace(string(jsonContent))) == 0 {
		return nil, errors.New("json group json is missing")
	}
	jsonGroup := &server.JsonGroup{Name: string(name)}
	params := make(map[string]interface{})
	err := json.Unmarshal(jsonContent, &params)
	if err != nil {
		return nil, err
	}
	jsonGroup.Params = params
	return jsonGroup, nil
}

func (self *Fs) GetTemplateCustomCss(path string) ([]byte, error) {
	cssPath := "templates/" + path + ".css"
	if _, err := os.Stat(cssPath); err != nil {
		return []byte{}, nil
	}
	return ioutil.ReadFile(cssPath)
}
