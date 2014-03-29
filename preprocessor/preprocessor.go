package preprocessor

type Css struct {
	Content []byte
}

func (self *Css) Get() ([]byte, error) {
	return self.Content, nil
}
