package abstract

import (
	"encoding/json"
	"io"
	"net/http"
)

const (
	whatsMyIPURL = "https://httpbin.org/ip"
)

type Body struct {
	Origin string `json:"origin"`
}

// GetIPAddress returns the IP address of the client using a get request to httpbin.org/ip
func (a *ServiceClient) GetIPAddress() (string, error) {
	req, err := http.NewRequest(http.MethodGet, whatsMyIPURL, nil)
	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	body, _ := io.ReadAll(res.Body)

	var resBody Body
	if err = json.Unmarshal(body, &resBody); err != nil {
		return "", err
	}

	return resBody.Origin, nil
}
