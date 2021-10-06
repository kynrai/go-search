package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/kynrai/go-search/search"
)

// Doc is a document for the index
type Doc struct {
	ID        string   `json:"id"`
	Firstname string   `json:"firstname"`
	Lastname  string   `json:"lastname"`
	Postcode  string   `json:"postcode"`
	Gender    string   `json:"gender"`
	Location  string   `json:"location"`
	Tags      []string `json:"tags"`
}

var docs = []*Doc{
	{
		ID:        "1",
		Firstname: "James",
		Lastname:  "Smith",
		Postcode:  "AB1 2CD",
		Gender:    "male",
		Location:  "London",
		Tags:      []string{"Brit"},
	},
	{
		ID:        "2",
		Firstname: "Jane",
		Lastname:  "Smith",
		Postcode:  "EF1 2GH",
		Gender:    "female",
		Location:  "London",
		Tags:      []string{"Brit"},
	},
	{
		ID:        "3",
		Firstname: "Brian",
		Lastname:  "Jones",
		Postcode:  "IJ1 2KL",
		Gender:    "male",
		Location:  "London",
		Tags:      []string{"Brit"},
	},
	{
		ID:        "4",
		Firstname: "Tom",
		Lastname:  "Evans",
		Postcode:  "MN1 2OP",
		Gender:    "male",
		Location:  "Tokyo",
		Tags:      []string{"Japanese"},
	},
	{
		ID:        "5",
		Firstname: "Sally",
		Lastname:  "Brown",
		Postcode:  "QR1 2ST",
		Gender:    "female",
		Location:  "New York",
		Tags:      []string{"Yank"},
	},
	{
		ID:        "6",
		Firstname: "Leslie",
		Lastname:  "Ann",
		Postcode:  "UV1 2WX",
		Gender:    "female",
		Location:  "Paris",
		Tags:      []string{"French"},
	},
	{
		ID:        "6",
		Firstname: "Janus",
		Lastname:  "Alan",
		Postcode:  "UV1 2WX",
		Gender:    "male",
		Location:  "Berlin",
		Tags:      []string{"German"},
	},
}

func main() {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	index := search.NewIndex(
		es,
		"testing",
		search.SchemaParams{
			SearchFields: []string{"firstname", "lastname", "postcode"},
			Fields: []search.Field{
				{Name: "gender", Type: "keyword"},
				{Name: "location", Type: "keyword"},
			},
			Arrays: []string{"tags"},
		},
	)
	err = index.Delete(ctx)
	if err != nil {
		log.Print(err)
	}

	err = index.Create(ctx)
	if err != nil {
		log.Print(err)
	}
	for _, doc := range docs {
		err = index.InsertDocument(ctx, doc.ID, doc)
		if err != nil {
			log.Fatal(err)
		}
	}

	doSearch(ctx, search.QueryParams{
		Query: "brit",
		// Terms: []search.Term{
		// 	{Name: "genders", Field: "gender"},
		// 	{Name: "locations", Field: "location"},
		// },
		// Filters: []search.Filter{
		// 	{Field: "gender", Values: []string{"male", "female"}},
		// 	{Field: "location", Values: []string{"London", "Berlin"}},
		// },
		// Sort: []search.Sort{
		// 	{Field: "gender", Direction: search.SortDesc},
		// },
		// Size: search.Int(1),
		// From: search.Int(1),
	}, index)

	doSearch(ctx, search.QueryParams{
		Query: "brit",
		Filters: []search.Filter{
			{Field: "gender", Values: []string{"female"}},
		},
	}, index)

	doSearch(ctx, search.QueryParams{
		Query: "brit",
	}, index)

}

func doSearch(ctx context.Context, params search.QueryParams, index *search.Index) {
	fmt.Println("REQUEST START")
	b, err := index.Search(ctx, params)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(string(b))
	resp, err := search.ParseResponse(b)
	if err != nil {
		log.Fatal(err)
	}
	var docs []Doc
	err = resp.HitsSource(&docs)
	if err != nil {
		log.Fatal(err)
	}
	b, _ = json.MarshalIndent(docs, "", "  ")
	fmt.Println(string(b))
	// fmt.Println(resp.TotalHits())
}
