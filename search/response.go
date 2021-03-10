package search

import "encoding/json"

// Response from a search call
type Response struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore interface{} `json:"max_score"`
		Hits     []struct {
			Index  string      `json:"_index"`
			Type   string      `json:"_type"`
			ID     string      `json:"_id"`
			Score  interface{} `json:"_score"`
			Source interface{} `json:"_source"`
			Sort   []string    `json:"sort"`
		} `json:"hits"`
	} `json:"hits"`
	Aggregations map[string]struct {
		DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
		SumOtherDocCount        int `json:"sum_other_doc_count"`
		Buckets                 []struct {
			Key      string `json:"key"`
			DocCount int    `json:"doc_count"`
		} `json:"buckets"`
	} `json:"aggregations"`
}

// ParseResponse from a search call into a more friendly struct
func ParseResponse(b []byte) (*Response, error) {
	var resp Response
	err := json.Unmarshal(b, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// HitsSource into a given slice, v must be a slice
func (r *Response) HitsSource(v interface{}) error {
	var sources []interface{}
	for _, hit := range r.Hits.Hits {
		sources = append(sources, hit.Source)
	}
	b, _ := json.Marshal(sources)
	return json.Unmarshal(b, v)
}

// TotalHits conveniance method to return total number of results
func (r *Response) TotalHits() int {
	return r.Hits.Total.Value
}
