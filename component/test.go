package component

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/romanoff/ahc/htmldiff"
	"github.com/romanoff/htmlcompressor"
)

type TestPool struct {
	TestSuites []*TestSuite
}

func (self *TestPool) Run(stopOnFailure bool) error {
	for _, testSuite := range self.TestSuites {
		err := testSuite.Run(stopOnFailure)
		if err != nil && stopOnFailure {
			return err
		}
	}
	return nil
}

type TestSuite struct {
	Compressor *htmlcompressor.HtmlCompressor
	Pool       *Pool
	Component  *Component
	Tests      []*Test
}

func (self *TestSuite) Run(stopOnFailure bool) error {
	for _, test := range self.Tests {
		err := test.Run(self.Component, self.Pool, self.Compressor)
		if err != nil && stopOnFailure {
			return err
		}
	}
	return nil
}

type Test struct {
	Identifier string
	Params     map[string]interface{}
	Expected   []byte
}

func (self *Test) Run(component *Component, pool *Pool, compressor *htmlcompressor.HtmlCompressor) error {
	html, err := component.Render(self.Params, pool)
	if err != nil {
		errorContent := fmt.Sprintf("Error while rendering component %v : %v", component.Namespace, err)
		if self.Identifier != ""{
			fmt.Println(errorContent)
		}
		return errors.New(errorContent)
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
