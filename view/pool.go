package view

type Pool struct {
	Templates map[string]*Template
	Pools     []*Pool
}

func (self *Pool) GetTemplate(path string) *Template {
	template := self.Templates[path]
	if template != nil {
		return template
	}
	for _, pool := range self.Pools {
		template = pool.GetTemplate(path)
		if template != nil {
			return template
		}
	}
	return nil
}
