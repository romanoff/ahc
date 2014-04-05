package preprocessor

import (
	"regexp"
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
var selectorRe = regexp.MustCompile("(?s)([^{]*)\\s*{(.*?)}")
var classRe = regexp.MustCompile("\\.(\\w+)")

func (self *Css) Classes() ([]string, error) {
	matches := selectorRe.FindAllSubmatch(self.Content, -1)
	classes := []string{}
	for _, match := range matches {
		if len(match) == 3 {
			classesMatch := classRe.FindAllSubmatch(match[1], -1)
			for _, classMatch := range classesMatch {
				if len(classMatch) == 2 {
					classes = append(classes, string(classMatch[1]))
				}
			}
		}
	}
	return classes, nil
}
