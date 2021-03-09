package search_test

import (
	"bytes"
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
	buf := &bytes.Buffer{}
	params := search.QueryParams{
		Query: "smith",
		Terms: []search.Term{
			{Name: "genders", Field: "gender"},
		},
		Filters: []search.Filter{
			{Field: "gender", Values: []string{"male"}},
		},
	}
	search := q.Search(params)
	if err := search.JSON(buf); err != nil {
		t.Fatal(err)
	}
	if expected != buf.String() {
		t.Fatal("unexpected query JSON output")
	}
}
