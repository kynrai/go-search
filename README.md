# Go Search

Out of the box configuration for a common elasticsearch use case

## Overview

This is a conveniance wrapper around managing an elastic search index which is used for multi field search as you type scenarios.

For example if you had documents which you wanted to run search as you type across `firstname` `lastname` and `postcode`, you can use this library to easily manage the indicies and search query boilerplate.

Difficult to manage query syntax and response formats becomes:

```go
index := search.NewIndex(
    es,
    "testing",
    []string{"firstname", "lastname", "postcode"},
)

err = index.Create(ctx)
if err != nil {
    log.Print(err)
}
for _, doc := range docs {
    err = index.Insert(ctx, doc.ID, doc)
    if err != nil {
        log.Fatal(err)
    }
}
index.Search(ctx, "bro s ef")
```
