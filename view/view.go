package view

import (
	"github.com/romanoff/ahc/component"
)

type View struct {
	Tags []Tag
}

type Attribute struct {
	Name  string
	Value string
}

type HtmlTag struct {
	Name       string
	Attributes []*Attribute
	Uuid       string
	Children   []Tag
}

func (self *HtmlTag) GetContent(pool *component.Pool) []byte {
	return nil
}

type AhcTag struct {
	Uuid         string
	Name         string
	Params       map[string][]Tag
	DefaultParam []Tag
}

func (self *AhcTag) GetContent(pool *component.Pool) []byte {
	return nil
}

type Tag interface {
	GetContent(*component.Pool) []byte
}
