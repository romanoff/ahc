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
	levelFields := make(map[int]*schema.Field)
	for _, line := range lines {
		field, level, err := self.parseSchemaField(line)
		if err != nil && len(bytes.TrimSpace(line)) != 0 {
			return nil, errors.New(fmt.Sprintf("Couldn't parse field: '%s' for schema in following file: '%v'", line, filepath))
		}
		if err == nil {
			if level == 0 {
				fields = append(fields, field)
			} else {
				parentField := levelFields[level-1]
				if parentField == nil {
					return nil, errors.New(fmt.Sprintf("No parent field for the following line: '%s' in following file: '%v'", line, filepath))
				}
				if parentField.Type == schema.ArrayField {
					parentField.ArrayValues = append(parentField.ArrayValues, field)
				} else if parentField.Type == schema.ObjectField {
					parentField.ObjectFields = append(parentField.ObjectFields, field)
				} else {
					return nil, errors.New(fmt.Sprintf("Non array of object parent field for the following line: '%v' in following file: '%v'", line, filepath))
				}

			}
			levelFields[level] = field
		}
	}
	schema := &schema.Schema{Fields: fields}
	return schema, nil
}

var fieldRe *regexp.Regexp = regexp.MustCompile("(\\s*)([\\w|-|_]+)\\s*\\{(\\w+)(=)?\\}({[^}]+})?\\s*(.*)\\s*")

func (self *Fs) parseSchemaField(fieldContent []byte) (*schema.Field, int, error) {
	matches := fieldRe.FindSubmatch(fieldContent)
	level := 0
	if len(matches) == 7 {
		level = len(matches[1]) / 2
		field := &schema.Field{
			Name:        string(matches[2]),
			Description: string(matches[6]),
			Required:    !(string(matches[4]) == "="),
		}
		fieldType := string(matches[3])
		switch {
		case fieldType == "string":
			field.Type = schema.StringField
		case fieldType == "num":
			field.Type = schema.NumField
		case fieldType == "bool":
			field.Type = schema.BoolField
		case fieldType == "object":
			field.Type = schema.ObjectField
		case fieldType == "array":
			field.Type = schema.ArrayField
		default:
			return nil, level, errors.New(fmt.Sprintf("Undefined field type: %s", fieldType))
		}
		if len(matches[5]) > 0 {
			allowedValues := string(matches[5][1 : len(matches[5])-1])
			field.AllowedValues = strings.Split(allowedValues, "|")
		}
		return field, level, nil
	}

	return nil, level, errors.New(fmt.Sprintf("Field '%s' could not be parsed", fieldContent))
}
