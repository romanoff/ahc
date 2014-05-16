package view

type Text struct {
	Uuid    string
	Content []byte
}

func (self *Text) GetContent(rParams *RenderParams) ([]byte, error) {
	return self.Content, nil
}
