package abstract_client

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strings"
)

// Client returns the resty client
func (a *AbstractServiceClient) Client() *http.Client {
	client := &http.Client{Transport: a}
	return client
}

// Path returns the path of this service client
func (a *AbstractServiceClient) Path() string {
	return a.serviceConfig.Service
}

// GetURL returns the url of the service client for the given protocol
func (a *AbstractServiceClient) GetURL(protocolInput ...string) string {
	return a.serviceConfig.FormatURL(protocolInput...)
}

// GetToken returns the token of the service client
func (a *AbstractServiceClient) GetToken() string {
	a.lock.Lock()
	defer a.lock.Unlock()

	return a.token
}

// SetToken sets the token of the service client
func (a *AbstractServiceClient) SetToken(token string) {
	a.lock.Lock()
	a.token = token
	a.lock.Unlock()

	a.headers.Store(AuthorizationHeader, a.BearerAuthHeaderValue())
}

// GetAPIKey returns the token of the service client
func (a *AbstractServiceClient) GetAPIKey() string {
	a.lock.Lock()
	defer a.lock.Unlock()

	return a.key
}

// SetAPIKey sets the token of the service client
func (a *AbstractServiceClient) SetAPIKey(key string) {
	a.headers.Store(APIKeyHeader, key)

	a.lock.Lock()
	defer a.lock.Unlock()

	a.key = key
}

// GetURLWithPath returns the url of the service client with the given path appended
func (a *AbstractServiceClient) GetURLWithPath(path string, protocolInput ...string) string {
	if strings.HasPrefix(path, "/") {
		return fmt.Sprintf("%s%s", a.serviceConfig.FormatURL(protocolInput...), path)
	}

	return fmt.Sprintf("%s/%s", a.serviceConfig.FormatURL(protocolInput...), path)
}

// IsAuthenticated returns true if the client is authenticated
func (a *AbstractServiceClient) IsAuthenticated() bool {
	return a.GetAPIKey() != "" || a.GetToken() != ""
}

// NewRequest returns a new resty request
func (a *AbstractServiceClient) NewRequest() *resty.Request {
	req := a.client.R()
	a.addHeaders(req)
	return req
}

func (a *AbstractServiceClient) addHeaders(req *resty.Request) {
	a.headers.Range(func(key, value interface{}) bool {
		req.SetHeader(key.(string), value.(string))
		return true
	})
}

func (a *AbstractServiceClient) BearerAuthHeaderValue() string {
	return fmt.Sprintf("Bearer %s", a.GetToken())
}

// SetHeader sets a header of the service client
func (a *AbstractServiceClient) SetHeader(key string, value string) {
	a.headers.Store(key, value)
}

// GetAuthHeader gets the current auth header of the service client
func (a *AbstractServiceClient) GetAuthHeader() http.Header {
	if a.GetAPIKey() != "" {
		return http.Header{APIKeyHeader: []string{a.GetAPIKey()}}
	} else if a.GetToken() != "" {
		return http.Header{"Authorization": []string{fmt.Sprintf("Bearer %s", a.GetToken())}}
	}

	return http.Header{}
}

func (a *AbstractServiceClient) RoundTrip(req *http.Request) (*http.Response, error) {
	a.headers.Range(func(key, value interface{}) bool {
		req.Header.Add(key.(string), value.(string))
		return true
	})
	return a.client.GetClient().Transport.RoundTrip(req)
}
