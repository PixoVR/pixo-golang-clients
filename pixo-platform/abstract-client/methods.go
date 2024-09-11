package abstract_client

import (
	"github.com/go-resty/resty/v2"
)

// Get makes a GET request to the URL
func (a *AbstractServiceClient) Get(path string) (*resty.Response, error) {
	url := a.GetURLWithPath(path)

	res, err := a.FormatRequest().Get(url)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Post makes a POST request to the URL
func (a *AbstractServiceClient) Post(path string, body []byte) (*resty.Response, error) {
	url := a.GetURLWithPath(path)

	req := a.FormatRequest()
	if body != nil {
		req = req.SetBody(body)
	}

	res, err := req.Post(url)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Patch makes a PATCH request to the URL
func (a *AbstractServiceClient) Patch(path string, body []byte) (*resty.Response, error) {
	url := a.GetURLWithPath(path)

	res, err := a.FormatRequest().SetBody(body).Patch(url)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Put makes a PUT request to the URL
func (a *AbstractServiceClient) Put(path string, body []byte) (*resty.Response, error) {
	url := a.GetURLWithPath(path)

	res, err := a.FormatRequest().SetBody(body).Put(url)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Delete makes a DELETE request to the URL
func (a *AbstractServiceClient) Delete(path string) (*resty.Response, error) {
	url := a.GetURLWithPath(path)

	res, err := a.FormatRequest().Delete(url)
	if err != nil {
		return nil, err
	}

	return res, nil
}
