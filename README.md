# Go Search

Out of the box configuration for a common elasticsearch use case.

Planned to support:

- Search as you type
- Term aggregations (categories e.g. gender,genre,category etc...)
- Filters

This should cover most search scenarios with optional facets like a person database or product database

## Overview

This is a conveniance wrapper around managing an elastic search index which is used for multi field search as you type scenarios.

For example if you had documents which you wanted to run search as you type across `firstname` `lastname` and `postcode`, you can use this library to easily manage the indicies and search query boilerplate.

Difficult to manage query syntax and response formats becomes:

Full example [here](cmd/main.go)

```go
index := search.NewIndex(
    es,
    "testing",
    search.SchemaParams{
        SearchFields:  []string{"firstname", "lastname", "postcode"},
        KeywordFields: []string{"gender"},
    },
)
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
}
err = index.Search(ctx, params)
if err != nil {
    log.Fatal(err)
}
```
