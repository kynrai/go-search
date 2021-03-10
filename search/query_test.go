package search_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/kynrai/go-search/search"
)

func TestQueryJSON(t *testing.T) {
	t.Parallel()
	q := search.NewQuery("firstname", "surname", "postcode")

	// note: the tab spacing here very important, it just does not look nice formatted
	const expected = `{
	"query": {
		"bool": {
			"must": {
				"multi_match": {
					"query": "smith",
					"fields": [
						"firstname",
						"firstname._2gram",
						"firstname._3gram",
						"surname",
						"surname._2gram",
						"surname._3gram",
						"postcode",
						"postcode._2gram",
						"postcode._3gram"
					],
					"type": "bool_prefix"
				}
			},
			"filter": {
				"terms": {
					"gender": [
						"male"
					]
				}
			}
		}
	},
	"aggs": {
		"genders": {
			"terms": {
				"field": "gender"
			}
		}
	}
}
`
	params := search.QueryParams{
		Query: "ja",
		Terms: []search.Term{
			{Name: "genders", Field: "gender"},
			{Name: "locations", Field: "location"},
		},
		Filters: []search.Filter{
			{Field: "gender", Values: []string{"male", "female"}},
			{Field: "location", Values: []string{"London"}},
		},
		Sort: []search.Sort{
			{Field: "gender", Direction: search.SortDesc},
		},
		// Size: search.Int(1),
		// From: search.Int(1),
	}

	buf := &bytes.Buffer{}
	search := q.Search(params)
	if err := search.JSON(buf); err != nil {
		t.Fatal(err)
	}
	fmt.Println(buf.String())
	if expected != buf.String() {
		t.Fatal("unexpected query JSON output")
	}
}
