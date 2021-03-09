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
		"multi_match": {
			"query": "smith",
			"type": "bool_prefix",
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
			]
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
	search := q.Search("smith", []search.Term{
		{Name: "genders", Field: "gender"},
	})
	if err := search.JSON(buf); err != nil {
		t.Fatal(err)
	}
	fmt.Println(expected)
	fmt.Println(buf.String())
	if expected != buf.String() {
		t.Fatal("unexpected query JSON output")
	}
}
