package search

import (
	"context"
	"fmt"
	"io/ioutil"

	opensearch "github.com/opensearch-project/opensearch-go"
	opensearchapi "github.com/opensearch-project/opensearch-go/opensearchapi"
	opensearchutil "github.com/opensearch-project/opensearch-go/opensearchutil"
)

// Index manages a search index
type Index struct {
	Client *opensearch.Client
	Name   string
	Params SchemaParams
	Schema *Schema
	Query  *Query
}

// NewIndex returns a new managed search index with fields set to be the search_as_you_type fields
func NewIndex(client *opensearch.Client, name string, params SchemaParams) *Index {
	return &Index{
		Client: client,
		Name:   name,
		Params: params,
		Schema: NewSchema(params),
		Query:  NewQuery(params.SearchFields, params.Arrays),
	}
}

// Create sets the mappings for an index on a server
func (i *Index) Create(ctx context.Context) error {
	req := opensearchapi.IndicesCreateRequest{
		Index: i.Name,
		Body:  opensearchutil.NewJSONReader(i.Schema),
	}
	resp, err := req.Do(ctx, i.Client)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.IsError() {
		return parseError(resp.Body)
	}
	return nil
}

// Delete an index
func (i *Index) Delete(ctx context.Context) error {
	req := opensearchapi.IndicesDeleteRequest{
		Index: []string{i.Name},
	}
	resp, err := req.Do(ctx, i.Client)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.IsError() {
		return parseError(resp.Body)
	}
	return nil
}

// Update sets the mappings for an index on a server
func (i *Index) Update(ctx context.Context) error {
	req := opensearchapi.IndicesCreateRequest{
		Index: i.Name,
		Body:  opensearchutil.NewJSONReader(i.Schema),
	}
	resp, err := req.Do(ctx, i.Client)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.IsError() {
		return parseError(resp.Body)
	}
	return nil
}

// InsertDocument adds a document into the index
func (i *Index) InsertDocument(ctx context.Context, id string, doc interface{}) error {
	req := opensearchapi.IndexRequest{
		Index:      i.Name,
		DocumentID: id,
		Body:       opensearchutil.NewJSONReader(doc),
		Refresh:    "true",
	}
	resp, err := req.Do(context.Background(), i.Client)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.IsError() {
		return parseError(resp.Body)
	}
	return nil
}

// Search returns results for a given search terms
func (i *Index) Search(ctx context.Context, params QueryParams) ([]byte, error) {
	req := opensearchapi.SearchRequest{
		Index:  []string{i.Name},
		Body:   opensearchutil.NewJSONReader(i.Query.Search(params)),
		Pretty: true,
		Size:   params.Size,
		From:   params.From,
	}
	resp, err := req.Do(context.Background(), i.Client)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.IsError() {
		fmt.Println(resp)
		return nil, parseError(resp.Body)
	}
	return ioutil.ReadAll(resp.Body)
}
