package schema

type Schema struct {
	Fields []*Field
}

func (self *Schema) Validate(params map[string]interface{}) []error {
	errors := []error{}
	for _, field := range self.Fields {
		err := field.Validate(params[field.Name])
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

func (self *Schema) GetSchemaParams(params map[string]interface{}) map[string]interface{} {
	schemaParams := make(map[string]interface{})
	for _, field := range self.Fields {
		schemaParams[field.Name] = params[field.Name]
	}
	return schemaParams
}
