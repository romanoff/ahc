package preprocessor

import (
	"regexp"
	"fmt"
)

func Init() *Css {
	return &Css{Variables: make(map[string][]byte)}
}

type Css struct {
	Content   []byte
	Variables map[string][]byte
}

func (self *Css) Get() ([]byte, error) {
	err := self.ReplaceVariables()
	if err != nil {
		return nil, err
	}
	return self.Content, nil
}

// Return css classes
var selectorRe = regexp.MustCompile("(?s)(.*)\\s*{(.*?)}")
func (self *Css) Classes() ([]string, error) {
	matches := selectorRe.FindAllSubmatch(self.Content, -1)
	fmt.Println(len(matches))
	for _, match := range matches {
		fmt.Println("--------")
		fmt.Println(string(match[1]))
		fmt.Println("--------*")
	}
	return []string{}, nil
}
