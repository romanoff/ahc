package component

import (
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
	return nil, nil
}
