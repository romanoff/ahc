package server

type TemplateJson struct {
	JsonGroups []*JsonGroup
}

type JsonGroup struct {
	Name   string
	Params map[string]interface{}
}
