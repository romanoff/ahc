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
		if node.Type == xmlx.NT_TEXT {
			if strings.TrimSpace(node.Value) != "" {
				html = append(html, []byte(node.Value)...)
			}
		}
		if node.Type != xmlx.NT_ELEMENT {
			continue
		}
		namespace := node.Name.Local
		//If html node, not custom node
		if strings.Index(namespace, "-") == -1 {
			childNodesHtml, err := self.getNodesHtml(node.Children)
			if err != nil {
				return nil, err
			}
			nodeAttributes := ""
			for _, attribute := range node.Attributes {
				nodeAttributes += " " + attribute.Name.Local + "=\"" + attribute.Value + "\""
			}
			result := []byte("<" + namespace + nodeAttributes + ">")
			result = append(result, childNodesHtml...)
			result = append(result, []byte("</"+namespace+">")...)
			return result, nil
		}
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
	if len(node.Children) > 0 && nodesHaveNamespace(node.Children, node.Name.Local) {
		for _, child := range node.Children {
			content, err := self.getNodesHtml(child.Children)
			if err != nil {
				return nil, err
			}
			params[child.Name.Local] = string(content)
		}
	} else if component.DefaultParam != "" {
		content, err := self.getNodesHtml(node.Children)
		if err != nil {
			return nil, err
		}
		params[component.DefaultParam] = string(content)
	}
	return params, nil
}

func nodesHaveNamespace(nodes []*xmlx.Node, namespace string) bool {
	for _, node := range nodes {
		if node.Name.Space != namespace {
			return false
		}
	}
	return true
}
