package abstract_client

import (
	"encoding/json"
)

const (
	whatsMyIPURL = "https://httpbin.org/ip"
)

type Body struct {
	Origin string `json:"origin"`
}

// GetIPAddress returns the IP address of the client using a get request to httpbin.org/ip
func (a *AbstractServiceClient) GetIPAddress() (string, error) {
	response, err := a.client.R().Get(whatsMyIPURL)
	if err != nil {
		return "", err
	}

	var body Body
	if err = json.Unmarshal(response.Body(), &body); err != nil {
		return "", err
	}

	return body.Origin, nil
}
