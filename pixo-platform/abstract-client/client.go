package abstract_client

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"net/http"
)

// GetClient returns the resty client
func (p *PixoAbstractAPIClient) GetClient() *http.Client {
	return p.restClient.GetClient()
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
	req := p.restClient.R().
		SetHeader("Content-Type", "application/json")

	if p.token != "" {
		req.SetHeader("x-access-token", p.token).
			SetHeader("Authorization", fmt.Sprintf("Bearer %s", p.token))
	}

	for key, value := range p.headers {
		req.SetHeader(key, value)
	}

	return req
}

func (p *PixoAbstractAPIClient) AddHeader(key string, value string) {
	p.headers[key] = value
}

// Get makes a GET request to the API
func (p *PixoAbstractAPIClient) Get(path string) (*resty.Response, error) {
	url := p.GetURLWithPath(path)

	res, err := p.FormatRequest().Get(url)
	if err != nil {
		log.Error().Err(err).Msg("Failed to perform get request")
		return nil, err
	}

	if res.IsError() {
		log.Error().Err(err).Msg("Failed to get data from API")
		return nil, errors.New(string(res.Body()))
	}

	return res, nil
}

// Post makes a POST request to the API
func (p *PixoAbstractAPIClient) Post(path string, body []byte) (*resty.Response, error) {
	url := p.GetURLWithPath(path)

	req := p.FormatRequest().SetBody(body)
	res, err := req.Post(url)
	if err != nil {
		log.Error().Err(err).Msg("Failed to perform post request")
		return nil, err
	}

	if res.IsError() {
		log.Error().Err(err).Msg("Failed to post data to API")
		return res, errors.New("invalid HTTP response received")
	}

	return res, nil
}

// Patch makes a PATCH request to the API
func (p *PixoAbstractAPIClient) Patch(path string, body []byte) (*resty.Response, error) {
	url := p.GetURLWithPath(path)

	res, err := p.FormatRequest().SetBody(body).Patch(url)
	if err != nil {
		log.Error().Err(err).Msg("Failed to perform patch request")
		return nil, err
	}

	if res.IsError() {
		log.Error().Err(err).Msg("Failed to patch data to API")
		return nil, errors.New(string(res.Body()))
	}

	return res, nil
}

// Put makes a PUT request to the API
func (p *PixoAbstractAPIClient) Put(path string, body []byte) (*resty.Response, error) {
	url := p.GetURLWithPath(path)

	res, err := p.FormatRequest().SetBody(body).Put(url)
	if err != nil {
		log.Error().Err(err).Msg("Failed to perform put request")
		return nil, err
	}

	if res.IsError() {
		log.Error().Err(err).Msg("Failed to put data to API")
		return nil, errors.New(string(res.Body()))
	}

	return res, nil
}

// Delete makes a DELETE request to the API
func (p *PixoAbstractAPIClient) Delete(path string) (*resty.Response, error) {
	url := p.GetURLWithPath(path)

	res, err := p.FormatRequest().Delete(url)
	if err != nil {
		log.Error().Err(err).Msg("Failed to perform delete request")
		return nil, err
	}

	if res.IsError() {
		log.Error().Err(err).Msg("Failed to delete data from API")
		return nil, errors.New(string(res.Body()))
	}

	return res, nil
}
