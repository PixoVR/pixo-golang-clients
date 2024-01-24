package abstract_client

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

// Client returns the resty client
func (p *PixoAbstractAPIClient) Client() *resty.Client {
	return p.restyClient
}

// GetURL returns the url of the restClient
func (p *PixoAbstractAPIClient) GetURL() string {
	return p.url
}

// GetToken returns the token of the restClient
func (p *PixoAbstractAPIClient) GetToken() string {
	return p.token
}

// SetToken sets the token of the restClient
func (p *PixoAbstractAPIClient) SetToken(token string) {
	p.token = token
}

// GetAPIKey returns the token of the restClient
func (p *PixoAbstractAPIClient) GetAPIKey() string {
	return p.key
}

// SetAPIKey sets the token of the restClient
func (p *PixoAbstractAPIClient) SetAPIKey(key string) {
	p.key = key
}

// GetURLWithPath returns the url of the restClient with a path appended
func (p *PixoAbstractAPIClient) GetURLWithPath(path string) string {
	return fmt.Sprintf("%s/%s", p.url, path)
}

// IsAuthenticated returns true if the client is authenticated
func (p *PixoAbstractAPIClient) IsAuthenticated() bool {
	return p.token != ""
}

// FormatRequest formats the request headers needed for authentication
func (p *PixoAbstractAPIClient) FormatRequest() *resty.Request {
	req := p.restyClient.R().
		SetHeader("Content-Type", "application/json")

	if p.token != "" {
		req.SetHeader("Authorization", fmt.Sprintf("Bearer %s", p.token))
	}

	if p.key != "" {
		req.SetHeader("x-api-key", p.key)
	}

	for key, value := range p.headers {
		req.SetHeader(key, value)
	}

	return req
}

func (p *PixoAbstractAPIClient) AddHeader(key string, value string) {
	p.headers[key] = value
}
