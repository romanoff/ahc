package parse

import (
	"github.com/romanoff/ahc/component"
	"os"
	"errors"
	"fmt"
)

type Fs struct {
}

// Parses component when path to css or html file is provided
func (self *Fs) ParseComponent(filepath string) (*component.Component, error) {
	if _, err := os.Stat(filepath); err != nil {
		return nil, errors.New(fmt.Sprintf("Error whie parsing component: %v file doesn't exist", filepath))
	}
	return nil, nil
}
