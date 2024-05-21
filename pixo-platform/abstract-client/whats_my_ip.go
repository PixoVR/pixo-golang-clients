package abstract_client

import (
	"encoding/json"
)

const (
	whatsMyIPURL = "https://httpbin.org/ip"
)

func (a *AbstractServiceClient) GetIPAddress() (string, error) {
	response, err := a.restyClient.R().Get(whatsMyIPURL)
	if err != nil {
		return "", err
	}

	var body struct {
		Origin string `json:"origin"`
	}
	if err = json.Unmarshal(response.Body(), &body); err != nil {
		return "", err
	}

	return body.Origin, nil
}
