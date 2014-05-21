package server

import (
	"github.com/romanoff/ahc/component"
)

func GetComponentsCss(pool *component.Pool, components []*component.Component) ([]byte, error) {
	css := []byte{}
	for i, component := range components {
		componentCss, err := component.GetCss(pool)
		if err != nil {
			return nil, err
		}
		if i != 0 {
			css = append(css, []byte("\n")...)
		}
		css = append(css, componentCss...)
	}
	return css, nil
}
