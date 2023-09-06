package abstract_client

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

// GetURL returns the URL of the restClient
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

// GetURLWithPath returns the URL of the restClient with a path appended
func (p *PixoAbstractAPIClient) GetURLWithPath(path string) string {
	return fmt.Sprintf("%s/%s", p.url, path)
}

// IsAuthenticated returns true if the client is authenticated
func (p *PixoAbstractAPIClient) IsAuthenticated() bool {
	return p.token != ""
}

// FormatRequest formats the request headers needed for authentication
func (p *PixoAbstractAPIClient) FormatRequest() *resty.Request {
	return p.restClient.R().
		SetHeader("x-access-token", p.token).
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", p.token))
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
		return nil, errors.New(string(res.Body()))
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
