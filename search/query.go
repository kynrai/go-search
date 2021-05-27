package search

import (
	"encoding/json"
	"fmt"
	"io"
)

// Query constructs a multimatch query for the fields set as search_as_you_type
type Query struct {
	Query struct {
		Bool struct {
			Must struct {
				MultiMatch struct {
					Query  string   `json:"query"`
					Fields []string `json:"fields"`
					Type   string   `json:"type"`
				} `json:"multi_match"`
			} `json:"must"`
			Filter []*filterTerm `json:"filter,omitempty"`
		} `json:"bool"`
	} `json:"query"`
	Sort map[string]string      `json:"sort,omitempty"`
	Aggs map[string]interface{} `json:"aggs,omitempty"`
}

type filterTerm struct {
	Terms map[string][]string `json:"terms,omitempty"`
}

type aggTerms struct {
	Terms aggTermsField `json:"terms"`
}

type aggTermsField struct {
	Field string `json:"field"`
}

// QueryParams collects all the data needed to build a search
type QueryParams struct {
	Query   string
	Terms   []Term
	Filters []Filter
	Size    *int
	From    *int
	Sort    []Sort
}

// Int returns a pointer to an int
func Int(i int) *int {
	return &i
}

// Filter defines a field should only have given values
type Filter struct {
	// Field we want to apply filters to
	Field string
	// Values that field can be
	Values []string
}

// Term defineds a field to return counts for all occurances
// https://www.elastic.co/guide/en/elasticsearch/reference/current/search-aggregations-bucket-terms-aggregation.html
type Term struct {
	// Name of the search bucket, usually the plural of the field name, e.g. genres, types, locations
	Name string
	// Field name of the document to bucket for each value
	Field string
}

type sortDirection string

// sort directions
const (
	SortAsc  sortDirection = "asc"
	SortDesc sortDirection = "desc"
)

// Sort defineds a field and direction to sort searches by
type Sort struct {
	Field     string
	Direction sortDirection
}

// NewQuery returns a query with the search fields set to search_as_you_type ngrams
func NewQuery(fields, arrays []string) *Query {
	searchFields := []string{}
	for _, field := range fields {
		searchFields = append(searchFields,
			field,
			fmt.Sprintf("%s._2gram", field),
			fmt.Sprintf("%s._3gram", field),
		)
	}
	for _, array := range arrays {
		searchFields = append(searchFields, array)
	}
	q := &Query{}
	q.Query.Bool.Must.MultiMatch.Type = "bool_prefix"
	q.Query.Bool.Must.MultiMatch.Fields = searchFields
	return q
}

// Search returns a clone of the query with search params set
func (q *Query) Search(params QueryParams) *Query {
	clone := q
	clone.Query.Bool.Must.MultiMatch.Query = params.Query
	clone.Aggs = make(map[string]interface{})

	if len(params.Terms) > 0 {
		for _, term := range params.Terms {
			clone.Aggs[term.Name] = &aggTerms{Terms: aggTermsField{Field: term.Field}}
		}
	}
	if len(params.Filters) > 0 {
		filters := []*filterTerm{}
		for _, filter := range params.Filters {
			ft := &filterTerm{}
			ft.Terms = make(map[string][]string)
			ft.Terms[filter.Field] = filter.Values
			filters = append(filters, ft)
		}
		clone.Query.Bool.Filter = filters
	}
	if params.Sort != nil {
		sorts := make(map[string]string)
		for _, sort := range params.Sort {
			sorts[sort.Field] = string(sort.Direction)
		}
		clone.Sort = sorts
	}
	return clone
}

// JSON returns the json body
func (q *Query) JSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	return enc.Encode(q)
}
