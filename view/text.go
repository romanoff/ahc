package view

import (
	"github.com/romanoff/ahc/component"
)

type Text struct {
	Uuid    string
	Content []byte
}

func (self *Text) GetContent(pool *component.Pool) []byte {
	return self.Content
}
