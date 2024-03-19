package abstract_client

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

// Client returns the resty client
func (a *AbstractServiceClient) Client() *resty.Client {
	return a.restyClient
}

// GetURL returns the url of the restClient
func (a *AbstractServiceClient) GetURL() string {
	return a.url
}

// GetToken returns the token of the restClient
func (a *AbstractServiceClient) GetToken() string {
	return a.token
}

// SetToken sets the token of the restClient
func (a *AbstractServiceClient) SetToken(token string) {
	a.token = token
}

// GetAPIKey returns the token of the restClient
func (a *AbstractServiceClient) GetAPIKey() string {
	return a.key
}

// SetAPIKey sets the token of the restClient
func (a *AbstractServiceClient) SetAPIKey(key string) {
	a.key = key
}

// GetURLWithPath returns the url of the restClient with a path appended
func (a *AbstractServiceClient) GetURLWithPath(path string) string {
	return fmt.Sprintf("%s/%s", a.url, path)
}

// IsAuthenticated returns true if the client is authenticated
func (a *AbstractServiceClient) IsAuthenticated() bool {
	return a.token != "" || a.key != ""
}

// FormatRequest formats the request headers needed for authentication
func (a *AbstractServiceClient) FormatRequest() *resty.Request {
	req := a.restyClient.R().
		SetHeader("Content-Type", "application/json")

	if a.token != "" {
		req.SetHeader("Authorization", fmt.Sprintf("Bearer %s", a.token))
	}

	if a.key != "" {
		req.SetHeader("x-api-key", a.key)
	}

	for key, value := range a.headers {
		req.SetHeader(key, value)
	}

	return req
}

func (a *AbstractServiceClient) SetHeader(key string, value string) {
	a.headers[key] = value
}
