package view

import (
	"github.com/jteeuwen/go-pkg-xmlx"
	"strings"
)

func InitView(template []byte) (*View, error) {
	document := xmlx.New()
	document.LoadBytes(template, nil)
	tags, err := getViewTags(document.Root.Children)
	if err != nil {
		return nil, err
	}
	view := &View{Tags: tags}
	return view, nil
}

func getViewTags(nodes []*xmlx.Node) ([]Tag, error) {
	tags := []Tag{}
	for _, node := range nodes {
		if node.Type == xmlx.NT_TEXT {
			if strings.TrimSpace(node.Value) != "" {
				tags = append(tags, &Text{Content: []byte(node.Value)})
			}
		}
		if node.Type != xmlx.NT_ELEMENT {
			continue
		}
		namespace := node.Name.Local
		//If html node, not custom node
		if strings.Index(namespace, "-") == -1 {
			htmlTag := &HtmlTag{Name: namespace}
			childTags, err := getViewTags(node.Children)
			if err != nil {
				return nil, err
			}
			htmlTag.Children = childTags
			for _, attribute := range node.Attributes {
				attr := &Attribute{Name: attribute.Name.Local, Value: attribute.Value}
				htmlTag.Attributes = append(htmlTag.Attributes, attr)
			}
			tags = append(tags, htmlTag)
			continue
		}
		//If ahc node
		ahcTag := &AhcTag{Name: namespace}
		params, defaultParam, err := getAhcNodeParams(node)
		if err != nil {
			return nil, err
		}
		ahcTag.Params = params
		ahcTag.DefaultParam = defaultParam
		tags = append(tags, ahcTag)

	}
	return tags, nil
}

func getAhcNodeParams(node *xmlx.Node) (map[string][]Tag, []Tag, error) {
	params := make(map[string][]Tag)
	defaultParam := make([]Tag, 0, 0)
	for _, attribute := range node.Attributes {
		params[attribute.Name.Local] = []Tag{&Text{Content: []byte(attribute.Value)}}
	}
	if len(node.Children) > 0 && nodesHaveNamespace(node.Children, node.Name.Local) {
		for _, child := range node.Children {
			tags, err := getViewTags(child.Children)
			if err != nil {
				return nil, nil, err
			}
			params[child.Name.Local] = tags
		}
	} else {
		tags, err := getViewTags(node.Children)
		if err != nil {
			return nil, nil, err
		}
		defaultParam = tags
	}
	return params, defaultParam, nil
}

func nodesHaveNamespace(nodes []*xmlx.Node, namespace string) bool {
	for _, node := range nodes {
		if node.Name.Space != namespace {
			return false
		}
	}
	return true
}
