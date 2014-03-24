package component

import (
	"errors"
	"fmt"
	"github.com/jteeuwen/go-pkg-xmlx"
	"strings"
)

type Pool struct {
	Components []*Component
	Pools      []*Pool
}

func (self *Pool) GetComponent(namespace string) *Component {
	for _, component := range self.Components {
		if namespace == component.Namespace || strings.HasSuffix(component.Namespace, "."+namespace) {
			return component
		}

	}
	for _, pool := range self.Pools {
		component := pool.GetComponent(namespace)
		if component != nil {
			return component
		}
	}
	return nil
}

func (self *Pool) Render(template []byte) ([]byte, error) {
	document := xmlx.New()
	document.LoadBytes(template, nil)
	html, err := self.getNodesHtml(document.Root.Children)
	if err != nil {
		return nil, err
	}
	return html, nil
}

func (self *Pool) getNodesHtml(nodes []*xmlx.Node) ([]byte, error) {
	html := []byte{}
	for _, node := range nodes {
		if node.Type != xmlx.NT_ELEMENT {
			continue
		}
		namespace := node.Name.Local
		component := self.GetComponent(namespace)
		if component == nil {
			return nil, errors.New(fmt.Sprintf("Component missing: %v", namespace))
		}
		params, err := self.getComponentParams(component, node)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error while parsing %v params: %v", namespace, err))
		}
		componentHtml, err := component.Render(params)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error while rendering %v: %v", namespace, err))
		}
		html = append(html, componentHtml...)
	}
	return html, nil
}

func (self *Pool) getComponentParams(component *Component, node *xmlx.Node) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	for _, attribute := range node.Attributes {
		params[attribute.Name.Local] = attribute.Value
	}
	return params, nil
}
