package preprocessor

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
func (self *Css) Classes() ([]string, error) {
	return []string{}, nil
}
