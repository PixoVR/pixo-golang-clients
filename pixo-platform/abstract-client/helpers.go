package abstract_client

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

// Client returns the resty client
func (a *AbstractServiceClient) Client() *resty.Client {
	return a.restyClient
}

// Path returns the name of the service client
func (a *AbstractServiceClient) Path() string {
	return a.serviceConfig.Service
}

// GetURL returns the url of the service client for the given protocol
func (a *AbstractServiceClient) GetURL(protocolInput ...string) string {
	return a.serviceConfig.FormatURL(protocolInput...)
}

// GetToken returns the token of the service client
func (a *AbstractServiceClient) GetToken() string {
	return a.token
}

// SetToken sets the token of the service client
func (a *AbstractServiceClient) SetToken(token string) {
	a.token = token
}

// GetAPIKey returns the token of the service client
func (a *AbstractServiceClient) GetAPIKey() string {
	return a.key
}

// SetAPIKey sets the token of the service client
func (a *AbstractServiceClient) SetAPIKey(key string) {
	a.key = key
}

// GetURLWithPath returns the url of the service client with the given path appended
func (a *AbstractServiceClient) GetURLWithPath(path string, protocolInput ...string) string {
	return fmt.Sprintf("%s/%s", a.serviceConfig.FormatURL(protocolInput...), path)
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
		req.SetHeader("x-access-token", a.token)
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
