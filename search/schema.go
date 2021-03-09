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
		Properties map[string]struct {
			Type string `json:"type"`
		} `json:"properties"`
	} `json:"mappings"`
}

// SchemaParams collects all the schema data needed to build mappings
type SchemaParams struct {
	// SearchFields defines fields used in search as you type
	SearchFields []string
	// KeywordFields defines fields that can be used for aggregations such as terms
	KeywordFields []string
}

// NewSchema returns a Schema with the mappings for the input fields set as search_as_you_type
func NewSchema(params SchemaParams) *Schema {
	s := &Schema{}
	s.Mappings.Properties = make(map[string]struct {
		Type string "json:\"type\""
	})
	for _, field := range params.SearchFields {
		s.Mappings.Properties[field] = struct {
			Type string "json:\"type\""
		}{Type: typeSearchAsYouType}
	}
	for _, field := range params.KeywordFields {
		s.Mappings.Properties[field] = struct {
			Type string "json:\"type\""
		}{Type: typeKeyword}
	}
	return s
}

// JSON writes the json representation of the schema to w with indent of tabs
func (s *Schema) JSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	return enc.Encode(s)
}
