package abstract_client

import (
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
)

// AbstractServiceClient is a struct that contains the url of the Pixo Service and a restyClient to make requests
type AbstractServiceClient struct {
	url            string
	token          string
	key            string
	headers        map[string]string
	restyClient    *resty.Client
	conn           *websocket.Conn
	timeoutSeconds int
}

type AbstractConfig struct {
	APIKey         string
	Token          string
	URL            string
	TimeoutSeconds int
}

// NewClient is a function that returns a AbstractServiceClient
func NewClient(config AbstractConfig) *AbstractServiceClient {

	if config.TimeoutSeconds == 0 {
		config.TimeoutSeconds = 30
	}

	return &AbstractServiceClient{
		url:            config.URL,
		token:          config.Token,
		key:            config.APIKey,
		timeoutSeconds: config.TimeoutSeconds,
		restyClient:    resty.New(),
		headers:        make(map[string]string),
	}
}
