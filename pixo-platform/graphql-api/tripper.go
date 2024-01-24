package graphql_api

import (
	"fmt"
	"net/http"
)

type transport struct {
	token               string
	key                 string
	underlyingTransport http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.token))
	}

	if t.key != "" {
		req.Header.Add("X-Api-Key", t.key)
	}

	return t.underlyingTransport.RoundTrip(req)
}
