package search

import (
	"encoding/json"
	"fmt"
	"io"
)

// ESError is the format elasticsearch errors are returned
type ESError struct {
	Err struct {
		RootCause []struct {
			Type      string `json:"type"`
			Reason    string `json:"reason"`
			IndexUUID string `json:"index_uuid"`
			Index     string `json:"index"`
		} `json:"root_cause"`
		Type      string `json:"type"`
		Reason    string `json:"reason"`
		IndexUUID string `json:"index_uuid"`
		Index     string `json:"index"`
	} `json:"error"`
	Status int `json:"status"`
}

func parseError(body io.Reader) *ESError {
	var esErr ESError
	err := json.NewDecoder(body).Decode(&esErr)
	if err != nil {
		esErr.Err.Reason = err.Error()
	}
	return &esErr
}

func (e *ESError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Status, e.Err.Reason)
}
