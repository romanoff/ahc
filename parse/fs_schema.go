package parse

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/romanoff/ahc/schema"
	"regexp"
	"strings"
)

func (self *Fs) ParseSchema(content []byte, filepath string) (*schema.Schema, error) {
	fields := make([]*schema.Field, 0, 0)
	lines := bytes.Split(content, []byte("\n"))
	for _, line := range lines {
		field, err := self.parseSchemaField(line)
		if err != nil && len(bytes.TrimSpace(line)) != 0 {
			return nil, errors.New(fmt.Sprintf("Couldn't parse field: '%s' for schema in following file: '%v'", line, filepath))
		}
		if err == nil {
			fields = append(fields, field)
		}
	}
	schema := &schema.Schema{Fields: fields}
	return schema, nil
}

var fieldRe *regexp.Regexp = regexp.MustCompile("\\s*([\\w|-|_]+)\\s*\\{(\\w+)(=)?\\}({[^}]+})?\\s*(.*)\\s*")

func (self *Fs) parseSchemaField(fieldContent []byte) (*schema.Field, error) {
	matches := fieldRe.FindSubmatch(fieldContent)
	if len(matches) == 6 {
		field := &schema.Field{
			Name:        string(matches[1]),
			Description: string(matches[5]),
			Required:    !(string(matches[3]) == "="),
		}
		fieldType := string(matches[2])
		switch {
		case fieldType == "string":
			field.Type = schema.StringField
		case fieldType == "num":
			field.Type = schema.NumField
		case fieldType == "bool":
			field.Type = schema.BoolField
		default:
			return nil, errors.New(fmt.Sprintf("Undefined field type: %s", fieldType))
		}
		if len(matches[4]) > 0 {
			allowedValues := string(matches[4][1 : len(matches[4])-1])
			field.AllowedValues = strings.Split(allowedValues, "|")
		}
		return field, nil
	}
	return nil, errors.New(fmt.Sprintf("Field '%s' could not be parsed", fieldContent))
}
