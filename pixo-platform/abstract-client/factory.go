package abstract_client

import (
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
)

// PixoAbstractAPIClient is a struct that contains the url of the Pixo Service and a restClient to make requests
type PixoAbstractAPIClient struct {
	url            string
	timeoutSeconds int
	token          string
	headers        map[string]string
	restClient     *resty.Client
	conn           *websocket.Conn
}

// NewClient is a function that returns a PixoAbstractAPIClient
func NewClient(token, apiURL string) *PixoAbstractAPIClient {

	return &PixoAbstractAPIClient{
		url:        apiURL,
		restClient: resty.New(),
		token:      token,
		headers:    make(map[string]string),
	}
}
