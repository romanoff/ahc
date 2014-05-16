package view

import (
	"strings"
)

type HtmlTag struct {
	Name       string
	Attributes []*Attribute
	Uuid       string
	Children   []Tag
}

var selfClosingTags string = "area,base,br,col,embed,hr,img,input,keygen,link,menuitem,meta param,source,track,wbr"

func (self *HtmlTag) GetContent(rParams *RenderParams) ([]byte, error) {
	nodeAttributes := ""
	for _, attribute := range self.Attributes {
		nodeAttributes += " " + attribute.Name + "=\"" + attribute.Value + "\""
	}
	if len(self.Children) == 0 && strings.Index(selfClosingTags, self.Name) != -1 {
		return []byte("<" + self.Name + nodeAttributes + " />"), nil
	}
	html := []byte{}
	html = append(html, []byte("<"+self.Name+nodeAttributes+">")...)
	for _, child := range self.Children {
		content, err := child.GetContent(rParams)
		if err != nil {
			return nil, err
		}
		html = append(html, content...)
	}
	html = append(html, []byte("</"+self.Name+">")...)
	return html, nil
}
