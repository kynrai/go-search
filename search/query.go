package search

import (
	"encoding/json"
	"fmt"
	"io"
)

// Query constructs a multimatch query for the fields set as search_as_you_type
type Query struct {
	Query struct {
		MultiMatch struct {
			Query  string   `json:"query"`
			Type   string   `json:"type"`
			Fields []string `json:"fields"`
		} `json:"multi_match"`
	} `json:"query"`
	Aggs map[string]interface{} `json:"aggs"`
}
type aggTerms struct {
	Terms aggTermsField `json:"terms"`
}

type aggTermsField struct {
	Field string `json:"field"`
}

// Filter defines a field should only have given values
type Filter struct {
	Field  string
	Values []string
}

// Term defineds a field to return counts for all occurances
// https://www.elastic.co/guide/en/elasticsearch/reference/current/search-aggregations-bucket-terms-aggregation.html
type Term struct {
	Name  string
	Field string
}

// NewQuery returns a query with the search fields set to search_as_you_type ngrams
func NewQuery(fields ...string) *Query {
	searchFields := []string{}
	for _, field := range fields {
		searchFields = append(searchFields,
			field,
			fmt.Sprintf("%s._2gram", field),
			fmt.Sprintf("%s._3gram", field),
		)
	}
	q := &Query{}
	q.Query.MultiMatch.Type = "bool_prefix"
	q.Query.MultiMatch.Fields = searchFields
	return q
}

// Search returns a clone of the query with search params set
func (q *Query) Search(query string, terms []Term) *Query {
	clone := q
	clone.Query.MultiMatch.Query = query
	clone.Aggs = make(map[string]interface{})

	if terms != nil {
		for _, term := range terms {
			t := &aggTerms{Terms: aggTermsField{Field: term.Field}}
			clone.Aggs[term.Name] = t
		}
	}
	return clone
}

// JSON returns the json body
func (q *Query) JSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	return enc.Encode(q)
}
