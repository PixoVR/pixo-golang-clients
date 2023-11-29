package abstract_client

import (
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
)

// PixoAbstractAPIClient is a struct that contains the url of the Pixo Service and a restClient to make requests
type PixoAbstractAPIClient struct {
	url            string
	token          string
	headers        map[string]string
	restClient     *resty.Client
	conn           *websocket.Conn
	timeoutSeconds int
}

// NewClient is a function that returns a PixoAbstractAPIClient
func NewClient(token, apiURL string, timeoutSeconds ...int) *PixoAbstractAPIClient {

	if len(timeoutSeconds) == 0 {
		timeoutSeconds = []int{30}
	}

	return &PixoAbstractAPIClient{
		url:            apiURL,
		restClient:     resty.New(),
		token:          token,
		headers:        make(map[string]string),
		timeoutSeconds: timeoutSeconds[0],
	}
}
