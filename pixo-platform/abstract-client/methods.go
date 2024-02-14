package abstract_client

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

// Get makes a GET request to the API
func (a *AbstractServiceClient) Get(path string) (*resty.Response, error) {
	url := a.GetURLWithPath(path)

	res, err := a.FormatRequest().Get(url)
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
func (a *AbstractServiceClient) Post(path string, body []byte) (*resty.Response, error) {
	url := a.GetURLWithPath(path)

	req := a.FormatRequest().SetBody(body)
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
func (a *AbstractServiceClient) Patch(path string, body []byte) (*resty.Response, error) {
	url := a.GetURLWithPath(path)

	res, err := a.FormatRequest().SetBody(body).Patch(url)
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
func (a *AbstractServiceClient) Put(path string, body []byte) (*resty.Response, error) {
	url := a.GetURLWithPath(path)

	res, err := a.FormatRequest().SetBody(body).Put(url)
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
func (a *AbstractServiceClient) Delete(path string) (*resty.Response, error) {
	url := a.GetURLWithPath(path)

	res, err := a.FormatRequest().Delete(url)
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
