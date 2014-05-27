package parse

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/romanoff/ahc/component"
	"io/ioutil"
	"os"
)

const (
	INPUT = iota
	EXPECTED
)

func (self *Fs) ParseComponentTest(filepath string, pool *component.Pool) (*component.TestSuite, error) {
	if pool == nil {
		return nil, errors.New("Components pool not supplied while parsing test")
	}
	c, err := self.ParseComponent(filepath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error while parsing component: %v", err))
	}
	namespace := c.Namespace
	c = pool.GetComponent(namespace)
	if c == nil {
		return nil, errors.New(fmt.Sprintf("Component %v not found in components pool", namespace))
	}
	testSuite := &component.TestSuite{Pool: pool, Component: c, Tests: make([]*component.Test, 0, 0)}
	if _, err := os.Stat(filepath); err != nil {
		return nil, errors.New(fmt.Sprintf("Error whie parsing component test: %v file doesn't exist", filepath))
	}
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error while reading test file: %v", filepath))
	}
	lines := bytes.Split(content, []byte("\n"))
	input := []byte{}
	expected := []byte{}
	n := 0
	mode := INPUT
	for _, line := range lines {
		line := bytes.TrimSpace(line)
		if string(line) == "Input:" {
			if len(input) != 0 && len(expected) != 0 {
				n++
				test, err := getTest(input, expected)
				if err != nil {
					return nil, err
				}
				test.Identifier = fmt.Sprintf("%s test %s", filepath, n)
				testSuite.Tests = append(testSuite.Tests, test)
			}
			input = []byte{}
			expected = []byte{}
			mode = INPUT
			continue
		}
		if string(line) == "Expected:" {
			mode = EXPECTED
			continue
		}
		if mode == INPUT {
			input = append(input, line...)
		}
		if mode == EXPECTED {
			expected = append(expected, line...)
		}
	}
	if len(input) != 0 && len(expected) != 0 {
		n++
		test, err := getTest(input, expected)
		if err != nil {
			return nil, err
		}
		test.Identifier = fmt.Sprintf("%s test %s", filepath, n)
		testSuite.Tests = append(testSuite.Tests, test)
	}
	return testSuite, nil
}

func getTest(input, expected []byte) (*component.Test, error) {
	var params map[string]interface{}
	err := json.Unmarshal(input, &params)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error while parsing test json: %v", err))
	}
	return &component.Test{Params: params, Expected: expected}, nil
}
