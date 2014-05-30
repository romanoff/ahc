package htmldiff

import (
	"bytes"
	"code.google.com/p/go-html-transform/h5"
	"code.google.com/p/go.net/html"
	"fmt"
	"github.com/foize/go.sgr"
)

type Error struct {
	Expected    string
	Got         string
	Description string
}

func CompareNode(originalNode, expectedNode *html.Node) *Error {
	err := &Error{
		Got:      h5.RenderNodesToString([]*html.Node{originalNode}),
		Expected: h5.RenderNodesToString([]*html.Node{expectedNode}),
	}
	if originalNode.Type != expectedNode.Type {
		err.Description = "Node type does not match"
		return err
	}
	if originalNode.Data != expectedNode.Data {
		err.Description = "Nodes data does not match"
		return err
	}
	for _, attr := range expectedNode.Attr {
		attrFound := false
		attrValueSame := false
		for _, originalAttr := range originalNode.Attr {
			if originalAttr.Key == attr.Key {
				attrFound = true
				if originalAttr.Val == attr.Val {
					attrValueSame = true
				}
			}
		}
		if !attrFound {
			err.Description = fmt.Sprintf("Attribute %v not found in node", attr.Key)
			return err
		}
		if !attrValueSame {
			err.Description = fmt.Sprintf("Attribute %v value is different", attr.Key)
			return err
		}
	}
	if len(originalNode.Attr) != len(expectedNode.Attr) {
		err.Description = "Different number of node attributes"
		return err
	}
	return CompareNodes(h5.Children(originalNode), h5.Children(expectedNode))
}

func CompareNodes(originalNodes, expectedNodes []*html.Node) *Error {
	if len(originalNodes) != len(expectedNodes) {
		return &Error{
			Description: fmt.Sprintf("Expected node to have %v elements, but got %v", len(expectedNodes), len(originalNodes)),
			Got:         h5.RenderNodesToString(originalNodes),
			Expected:    h5.RenderNodesToString(expectedNodes),
		}
	}
	for i, node := range originalNodes {
		expectedNode := expectedNodes[i]
		err := CompareNode(node, expectedNode)
		if err != nil {
			return err
		}
	}
	return nil
}

func Compare(actual, expected []byte) *Error {
	actualNodes, err := h5.Partial(bytes.NewReader(actual))
	if err != nil {
		return &Error{Description: "Parsing error", Got: string(actual)}
	}
	expectedNodes, err := h5.Partial(bytes.NewReader(expected))
	if err != nil {
		return &Error{Description: "Parsing error", Got: string(expected)}
	}
	return CompareNodes(actualNodes, expectedNodes)
}

func PrettyPrint(err *Error) {
	if err == nil {
		return
	}
	fmt.Println(sgr.MustParse("[fg-11]" + err.Description + "[reset]"))
	fmt.Println(sgr.MustParse("[fg-green]Expected:[reset]"))
	expected := fmt.Sprintf("-->%s<--\n", err.Expected)
	fmt.Printf(sgr.MustParse("[fg-green]" + expected + "[reset]"))
	fmt.Println(sgr.MustParse("[fg-red]Got:[reset]"))
	got := fmt.Sprintf("-->%s<--\n", err.Got)
	fmt.Printf(sgr.MustParse("[fg-red]" + got + "[reset]"))
}
