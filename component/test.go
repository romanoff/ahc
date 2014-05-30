package component

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/foize/go.sgr"
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
	fmt.Println()
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
		fmt.Printf(sgr.MustParse("[fg-red].[reset]"))
		errorContent := fmt.Sprintf("Error while rendering component %v : %v\n", component.Namespace, err)
		if self.Identifier != "" {
			fmt.Printf(sgr.MustParse("[fg-red]" + errorContent + "[reset]"))
		}
		return errors.New(errorContent)
	}
	diffError := htmldiff.Compare(self.FormatHtml(html, compressor),
		self.FormatHtml(self.Expected, compressor))
	if diffError != nil {
		fmt.Printf(sgr.MustParse("[fg-red].[reset]"))
		if self.Identifier != "" {
			errorContent := fmt.Sprintf("%v test failed:\n", self.Identifier)
			fmt.Printf(sgr.MustParse("[fg-red]" + errorContent + "[reset]"))
			htmldiff.PrettyPrint(diffError)
		}
		return errors.New("Test failed")
	}
	fmt.Printf(sgr.MustParse("[fg-green].[reset]"))
	return nil
}

func (self *Test) FormatHtml(html []byte, compressor *htmlcompressor.HtmlCompressor) []byte {
	return bytes.TrimSpace(compressor.Compress(html))
}
