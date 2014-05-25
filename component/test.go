package component

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/romanoff/ahc/htmldiff"
	"github.com/romanoff/htmlcompressor"
)

type TestSuite struct {
	Compressor *htmlcompressor.HtmlCompressor
	Pool       *Pool
	Component  *Component
	Tests      []*Test
}

type Test struct {
	Identifier string
	Params     map[string]interface{}
	Expected   []byte
}

func (self *Test) Run(component *Component, pool *Pool, compressor *htmlcompressor.HtmlCompressor) error {
	html, err := component.Render(self.Params, pool)
	if err != nil {
		return errors.New(fmt.Sprintf("Error while rendering component %v : %v", component.Namespace, err))
	}
	diffError := htmldiff.Compare(self.FormatHtml(html, compressor),
		self.FormatHtml(self.Expected, compressor))
	if diffError != nil {
		if self.Identifier != "" {
			fmt.Printf("%v test failed:\n", self.Identifier)
			htmldiff.PrettyPrint(diffError)
		}
		return errors.New("Test failed")
	}
	return nil
}

func (self *Test) FormatHtml(html []byte, compressor *htmlcompressor.HtmlCompressor) []byte {
	return bytes.TrimSpace(compressor.Compress(html))
}
