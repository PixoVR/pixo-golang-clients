package abstract

import (
	"context"
	"net/http"
	"strings"
)

// Get makes a GET request to the URL
func (a *ServiceClient) Get(ctx context.Context, path string) (*http.Response, error) {
	req, err := a.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return a.client.Do(req)
}

// Post makes a POST request to the URL
func (a *ServiceClient) Post(ctx context.Context, path string, body []byte) (*http.Response, error) {
	req, err := a.NewRequestWithContext(ctx, http.MethodPost, path, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	return a.client.Do(req)
}

// Patch makes a PATCH request to the URL
func (a *ServiceClient) Patch(ctx context.Context, path string, body []byte) (*http.Response, error) {
	req, err := a.NewRequestWithContext(ctx, http.MethodPatch, path, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	return a.client.Do(req)
}

// Put makes a PUT request to the URL
func (a *ServiceClient) Put(ctx context.Context, path string, body []byte) (*http.Response, error) {
	req, err := a.NewRequestWithContext(ctx, http.MethodPut, path, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	return a.client.Do(req)
}

// Delete makes a DELETE request to the URL
func (a *ServiceClient) Delete(ctx context.Context, path string) (*http.Response, error) {
	req, err := a.NewRequestWithContext(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	return a.client.Do(req)
}
