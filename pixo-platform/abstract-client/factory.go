package abstract_client

import (
	"github.com/go-resty/resty/v2"
)

// PixoAbstractAPIClient is a struct that contains the URL of the Pixo Service and a restClient to make requests
type PixoAbstractAPIClient struct {
	url        string
	token      string
	restClient *resty.Client
}

// NewClient is a function that returns a PixoAbstractAPIClient
func NewClient(token, apiURL string) *PixoAbstractAPIClient {
	if apiURL == "" {
		apiURL = getAPIURL()
	}

	return &PixoAbstractAPIClient{
		url:        apiURL,
		restClient: resty.New(),
		token:      token,
	}
}
