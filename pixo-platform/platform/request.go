package platform

import "encoding/json"

type GraphQLRequestPayload struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
	OperationName string                 `json:"operationName,omitempty"`
}

type GraphQLResponse struct {
	Messages []string        `json:"messages"`
	Errors   []Error         `json:"errors"`
	Data     json.RawMessage `json:"data"`
}

type Error struct {
	// Path elements are field names or, inside a list, the index of the element
	// the error came from, so they are not all strings.
	Path      []any  `json:"path"`
	Message   string `json:"message"`
	Locations []struct {
		Line int `json:"line"`
	}
}
