package graphql_api

import (
	"fmt"
	"net/http"
)

type transport struct {
	token               string
	underlyingTransport http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.token))
		req.Header.Add("X-Access-Token", t.token)
	}

	return t.underlyingTransport.RoundTrip(req)
}
