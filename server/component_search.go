package server

import (
	"errors"
	"fmt"
	"github.com/romanoff/ahc/component"
	"github.com/romanoff/ahc/view"
	"strings"
)

func InitComponentSearch(templatesPool *view.Pool) *ComponentSearch {
	return &ComponentSearch{
		Components:     make([]*component.Component, 0, 0),
		UsedNamespaces: make(map[string]bool),
		TemplatesPool:  templatesPool,
	}
}

type ComponentSearch struct {
	Components     []*component.Component
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
	return usedNamespaces, nil
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
