package search_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/kynrai/go-search/search"
)

func TestSchemaJSON(t *testing.T) {
	t.Parallel()
	s := search.NewSchema(search.SchemaParams{
		SearchFields:  []string{"firstname", "surname", "postcode"},
		KeywordFields: []string{"gender"},
	})
	// note: the tab spacing here very important, it just does not look nice formatted
	const expected = `{
	"mappings": {
		"properties": {
			"firstname": {
				"type": "search_as_you_type"
			},
			"gender": {
				"type": "keyword"
			},
			"postcode": {
				"type": "search_as_you_type"
			},
			"surname": {
				"type": "search_as_you_type"
			}
		}
	}
}
`
	buf := &bytes.Buffer{}
	if err := s.JSON(buf); err != nil {
		t.Fatal(err)
	}
	fmt.Println(expected)
	fmt.Println(buf.String())
	if expected != buf.String() {
		t.Fatal("unexpected schema JSON output")
	}
}
