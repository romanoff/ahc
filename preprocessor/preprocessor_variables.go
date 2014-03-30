package preprocessor

import (
	"bytes"
	"regexp"
)

func (self *Css) ReplaceVariables() error {
	lines := bytes.Split(self.Content, []byte("\n"))
	varSymbol := []byte("$")
	for i, line := range lines {
		if bytes.Contains(line, varSymbol) {
			line, err := self.processLineWithVariable(line)
			lines[i] = line
			if err != nil {
				return err
			}
		}
	}
	self.Content = bytes.Join(lines, []byte("\n"))
	return nil
}

var setVarRe *regexp.Regexp = regexp.MustCompile("\\s*\\$([\\w-]+)\\s*:\\s*(.+);")
var varRe *regexp.Regexp = regexp.MustCompile("\\$([\\w-]+)")

func (self *Css) processLineWithVariable(line []byte) ([]byte, error) {
	match := setVarRe.FindSubmatch(line)
	if len(match) == 3 {
		self.Variables[string(match[1])] = match[2]
		line = bytes.Replace(line, match[0], []byte{}, -1)
	} else {
		matches := varRe.FindAllSubmatch(line, -1)
		for _, match := range matches {
			varName := string(match[1])
			if self.Variables[varName] != nil {
				line = bytes.Replace(line, match[0], self.Variables[varName], -1)
			}
		}
	}
	return line, nil
}
