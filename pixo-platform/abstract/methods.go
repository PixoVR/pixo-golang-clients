package abstract

import (
	"context"
	"net/http"
	"strings"
)

// Get makes a GET request to the URL
func (a *ServiceClient) Get(ctx context.Context, path string) (*http.Response, error) {
	req, err := a.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req.WithContext(ctx))
}

// Post makes a POST request to the URL
func (a *ServiceClient) Post(ctx context.Context, path string, body []byte) (*http.Response, error) {
	req, err := a.NewRequest(http.MethodPost, path, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req.WithContext(ctx))
}

// Patch makes a PATCH request to the URL
func (a *ServiceClient) Patch(ctx context.Context, path string, body []byte) (*http.Response, error) {
	req, err := a.NewRequest(http.MethodPatch, path, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req.WithContext(ctx))
}

// Put makes a PUT request to the URL
func (a *ServiceClient) Put(ctx context.Context, path string, body []byte) (*http.Response, error) {
	req, err := a.NewRequest(http.MethodPut, path, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(req.WithContext(ctx))
}

// Delete makes a DELETE request to the URL
func (a *ServiceClient) Delete(ctx context.Context, path string) (*http.Response, error) {
	req, err := a.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(req.WithContext(ctx))
}
