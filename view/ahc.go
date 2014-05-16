package view

import (
	"errors"
	"fmt"
)

type AhcTag struct {
	Uuid         string
	Name         string
	Params       map[string][]Tag
	DefaultParam []Tag
}

func (self *AhcTag) GetContent(rParams *RenderParams) ([]byte, error) {
	pool := rParams.Pool
	if pool == nil {
		return nil, errors.New("Components pool is missing while rendering ahc tag")
	}
	component := pool.GetComponent(self.Name)
	if component == nil {
		return nil, errors.New(fmt.Sprintf("Component missing: %v", self.Name))
	}
	params := make(map[string]interface{})
	for name, tags := range self.Params {
		tagsContent, err := getTagsContent(tags, rParams)
		if err != nil {
			return nil, err
		}
		params[name] = string(tagsContent)
	}
	if len(self.DefaultParam) > 0 {
		if component.DefaultParam == "" {
			return nil, errors.New(fmt.Sprintf("Default parameter is missing from '%v' component", component.Namespace))
		}
		tagsContent, err := getTagsContent(self.DefaultParam, rParams)
		if err != nil {
			return nil, err
		}
		params[component.DefaultParam] = string(tagsContent)
	}
	// Populate params
	if rParams.Safe {
		return component.RenderSafe(params, pool)
	} else {
		return component.Render(params, pool)
	}
}
