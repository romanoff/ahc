package server

import (
	"errors"
	"fmt"
	"github.com/romanoff/ahc/component"
	"github.com/romanoff/ahc/view"
	"regexp"
	"strings"
)

func InitComponentSearch(templatesPool *view.Pool) *ComponentSearch {
	return &ComponentSearch{
		Components:     make([]*component.Component, 0, 0),
		UsedNamespaces: make(map[string]bool),
		UsedTemplates:  make(map[string]bool),
		TemplatesPool:  templatesPool,
	}
}

type ComponentSearch struct {
	Components     []*component.Component
	UsedTemplates  map[string]bool
	UsedNamespaces map[string]bool
	TemplatesPool  *view.Pool
}

func (self *ComponentSearch) Search(path string) error {
	if self.TemplatesPool == nil {
		return errors.New("Templates pool missing")
	}
	template := self.TemplatesPool.GetTemplate(path)
	if template == nil {
		return errors.New(fmt.Sprintf("Template %v was not found", path))
	}
	namespaces, err := self.GetUsedNamespaces([]byte(template.Content))
	if err != nil {
		return err
	}
	err = self.AddComponents(namespaces)
	if err != nil {
		return err
	}
	//TODO: check partial components as well as internal components of components
	return nil
}

const (
	IGNORE = iota
	TAG    = iota
	READ   = iota
	WAIT   = iota
)

func (self *ComponentSearch) GetUsedNamespaces(content []byte) ([]string, error) {
	usedNamespaces := []string{}
	currentNamespace := ""
	state := IGNORE
	usedMap := make(map[string]bool)
	for _, b := range content {
		if state == IGNORE && b != '<' {
			continue
		}
		if state == IGNORE && b == '<' {
			state = TAG
			continue
		}
		if state != IGNORE && b == '>' {
			//TODO Remove everything after :
			currentNamespace = strings.Split(currentNamespace, ":")[0]
			if currentNamespace != "" && usedMap[currentNamespace] == false && strings.Index(currentNamespace, "-") != -1 {
				usedMap[currentNamespace] = true
				usedNamespaces = append(usedNamespaces, currentNamespace)
			}
			currentNamespace = ""
			state = IGNORE
		}
		if state == TAG && b == '!' { // It's xml comment
			state = IGNORE
			continue
		}
		if state == TAG && b != ' ' && b != '\t' && b != '\n' && b != '/' {
			currentNamespace += string(b)
			state = READ
			continue
		}
		if state == READ {
			if b != ' ' && b != '\t' && b != '\n' && b != '/' {
				currentNamespace += string(b)
			} else {
				state = WAIT
			}
		}
	}
	templates := self.GetUsedTemplates(content)
	for _, path := range templates {
		template := self.TemplatesPool.GetTemplate(path)
		if template == nil {
			return nil, errors.New(fmt.Sprintf("Template %v was not found", path))
		}
		namespaces, err := self.GetUsedNamespaces([]byte(template.Content))
		if err != nil {
			return nil, err
		}
		for _, namespace := range namespaces {
			if !usedMap[namespace] {
				usedNamespaces = append(usedNamespaces, namespace)
				usedMap[namespace] = true
			}
		}
	}
	return usedNamespaces, nil
}

var partialRe *regexp.Regexp = regexp.MustCompile("{{\\s*template\\s+\"([\\w|/]+)\"[\\s|\\.]*}}")

func (self *ComponentSearch) GetUsedTemplates(content []byte) []string {
	allMatches := partialRe.FindAllSubmatch(content, -1)
	usedTemplates := []string{}
	for _, matches := range allMatches {
		if len(matches) == 2 {
			templateName := string(matches[1])
			if self.UsedTemplates[templateName] {
				continue
			}
			usedTemplates = append(usedTemplates, templateName)
			self.UsedTemplates[templateName] = true
		}
	}
	return usedTemplates
}

func (self *ComponentSearch) AddComponents(namespaces []string) error {
	if self.TemplatesPool == nil {
		return errors.New("Templates pool is missing")
	}
	for _, namespace := range namespaces {
		component := self.TemplatesPool.ComponentsPool.GetComponent(namespace)
		if component == nil {
			return errors.New(fmt.Sprintf("Component not found: %v", namespace))
		}
		if self.UsedNamespaces[component.Namespace] == false {
			self.UsedNamespaces[component.Namespace] = true
			self.Components = append(self.Components, component)
		}
	}
	return nil
}
