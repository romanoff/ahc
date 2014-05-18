package server

import (
	"github.com/romanoff/ahc/component"
	"github.com/romanoff/ahc/view"
)

type ComponentSearch struct {
	Components    *[]component.Component
	TemplatesPool *view.Pool
}

func (self *ComponentSearch) Search(path string) error {
	return nil
}
