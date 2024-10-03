package abstract

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Client returns the resty client
func (a *ServiceClient) Client() *http.Client {
	return &http.Client{Transport: a}
}

// Path returns the path of this service client
func (a *ServiceClient) Path() string {
	return a.serviceConfig.Service
}

// GetURL returns the url of the service client for the given protocol
func (a *ServiceClient) GetURL(protocolInput ...string) string {
	return a.serviceConfig.FormatURL(protocolInput...)
}

// GetToken returns the token of the service client
func (a *ServiceClient) GetToken() string {
	a.lock.Lock()
	defer a.lock.Unlock()

	return a.token
}

// SetToken sets the token of the service client
func (a *ServiceClient) SetToken(token string) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.token = token
	a.headers[AuthorizationHeaderKey] = bearerAuthHeaderValue(token)
}

// GetAPIKey returns the token of the service client
func (a *ServiceClient) GetAPIKey() string {
	a.lock.Lock()
	defer a.lock.Unlock()

	return a.key
}

// SetAPIKey sets the token of the service client
func (a *ServiceClient) SetAPIKey(key string) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.key = key
	a.headers[APIKeyHeaderKey] = key
}

// GetURLWithPath returns the url of the service client with the given path appended
func (a *ServiceClient) GetURLWithPath(path string, protocolInput ...string) string {
	if strings.HasPrefix(path, "/") {
		return fmt.Sprintf("%s%s", a.serviceConfig.FormatURL(protocolInput...), path)
	}

	return fmt.Sprintf("%s/%s", a.serviceConfig.FormatURL(protocolInput...), path)
}

// IsAuthenticated returns true if the client is authenticated
func (a *ServiceClient) IsAuthenticated() bool {
	return a.GetAPIKey() != "" || a.GetToken() != ""
}

// NewRequest returns a new resty request
func (a *ServiceClient) NewRequest(method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, a.GetURLWithPath(path), body)
	if err != nil {
		return nil, err
	}
	return a.addHeaders(req), nil
}

// AddQueryParams adds query parameters to the request
func (a *ServiceClient) AddQueryParams(req *http.Request, params map[string]string) *http.Request {
	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
	return req
}

func (a *ServiceClient) addHeaders(req *http.Request) *http.Request {
	a.lock.Lock()
	defer a.lock.Unlock()

	for key, value := range a.headers {
		req.Header.Set(key, value)
	}
	return req
}

func bearerAuthHeaderValue(token string) string {
	return fmt.Sprintf("Bearer %s", token)
}

// SetHeader sets a header of the service client
func (a *ServiceClient) SetHeader(key string, value string) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.headers[key] = value
}

// GetAuthHeader gets the current auth header of the service client
func (a *ServiceClient) GetAuthHeader() http.Header {
	if a.GetAPIKey() != "" {
		return http.Header{APIKeyHeaderKey: []string{a.GetAPIKey()}}
	} else if a.GetToken() != "" {
		return http.Header{AuthorizationHeaderKey: []string{fmt.Sprintf("Bearer %s", a.GetToken())}}
	}

	return http.Header{}
}

func (a *ServiceClient) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header = a.addHeaders(req).Header
	return http.DefaultClient.Do(req)
}
