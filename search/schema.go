package search

import (
	"encoding/json"
	"io"
)

const (
	typeSearchAsYouType = `search_as_you_type`
	typeKeyword         = `keyword`
)

// Schema sets the mappings for the elasticsearch index
type Schema struct {
	Mappings struct {
		Properties map[string]schemaProperties `json:"properties"`
	} `json:"mappings"`
}

type schemaProperties struct {
	Type string `json:"type"`
}

// SchemaParams collects all the schema data needed to build mappings
type SchemaParams struct {
	// SearchFields defines fields used in search as you type
	SearchFields []string
	// Fields allow custom types to be set for fields e.g. date
	Fields []Field
	// Allow searching thorough arrays
	Arrays []string
}

// Field defined a field and its mapping type
type Field struct {
	Name string
	Type string
}

// NewSchema returns a Schema with the mappings for the input fields set as search_as_you_type
func NewSchema(params SchemaParams) *Schema {
	s := &Schema{}
	s.Mappings.Properties = make(map[string]schemaProperties)
	for _, field := range params.SearchFields {
		s.Mappings.Properties[field] = schemaProperties{Type: typeSearchAsYouType}
	}
	for _, field := range params.Fields {
		s.Mappings.Properties[field.Name] = schemaProperties{Type: field.Type}
	}
	return s
}

// JSON writes the json representation of the schema to w with indent of tabs
func (s *Schema) JSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	return enc.Encode(s)
}
