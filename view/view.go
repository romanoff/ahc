package view

import (
	"github.com/romanoff/ahc/component"
)

type View struct {
	Tags []Tag
}

func (self *View) GetContent(rParams *RenderParams) ([]byte, error) {
	return getTagsContent(self.Tags, rParams)
}

func getTagsContent(tags []Tag, rParams *RenderParams) ([]byte, error) {
	tagsContent := []byte{}
	for _, tag := range tags {
		content, err := tag.GetContent(rParams)
		if err != nil {
			return nil, err
		}
		tagsContent = append(tagsContent, content...)
	}
	return tagsContent, nil
}

type Attribute struct {
	Name  string
	Value string
}

type RenderParams struct {
	Pool *component.Pool
	Safe bool
}

type Tag interface {
	GetContent(*RenderParams) ([]byte, error)
}
