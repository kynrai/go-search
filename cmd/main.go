package main

import (
	"context"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/kynrai/go-search/search"
)

// Doc is a document for the index
type Doc struct {
	ID        string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Postcode  string `json:"postcode"`
	Gender    string `json:"gender"`
}

var docs = []*Doc{
	{
		ID:        "1",
		Firstname: "John",
		Lastname:  "Smith",
		Postcode:  "AB1 2CD",
		Gender:    "male",
	},
	{
		ID:        "2",
		Firstname: "Jane",
		Lastname:  "Smith",
		Postcode:  "EF1 2GH",
		Gender:    "female",
	},
	{
		ID:        "3",
		Firstname: "Brian",
		Lastname:  "Jones",
		Postcode:  "IJ1 2KL",
		Gender:    "male",
	},
	{
		ID:        "4",
		Firstname: "Tom",
		Lastname:  "Evans",
		Postcode:  "MN1 2OP",
		Gender:    "male",
	},
	{
		ID:        "5",
		Firstname: "Sally",
		Lastname:  "Brown",
		Postcode:  "QR1 2ST",
		Gender:    "female",
	},
	{
		ID:        "6",
		Firstname: "Leslie",
		Lastname:  "Ann",
		Postcode:  "UV1 2WX",
		Gender:    "female",
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
			},
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
	params := search.QueryParams{
		Query: "smith",
		Terms: []search.Term{
			{Name: "genders", Field: "gender"},
		},
		Filters: []search.Filter{
			{Field: "gender", Values: []string{"male", "female"}},
		},
		Sort: []search.Sort{
			{Field: "gender", Direction: search.SortDesc},
		},
		Size: search.Int(1),
		From: search.Int(1),
	}
	b, err := index.Search(ctx, params)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
