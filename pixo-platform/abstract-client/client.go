package abstract_client

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
)

// AbstractServiceClient is a struct that handles generic http client operations
type AbstractServiceClient struct {
	serviceConfig urlfinder.ServiceConfig
	token         string
	key           string
	headers       map[string]string

	restyClient    *resty.Client
	websocketConn  *websocket.Conn
	timeoutSeconds int
}

// AbstractConfig is a struct that holds the configuration for the AbstractServiceClient
type AbstractConfig struct {
	ServiceConfig  urlfinder.ServiceConfig
	APIKey         string
	Token          string
	TimeoutSeconds int
}

// NewClient creates a new AbstractServiceClient given a config struct
func NewClient(config AbstractConfig) *AbstractServiceClient {

	if config.TimeoutSeconds == 0 {
		config.TimeoutSeconds = 30
	}

	return &AbstractServiceClient{
		serviceConfig:  config.ServiceConfig,
		token:          config.Token,
		key:            config.APIKey,
		timeoutSeconds: config.TimeoutSeconds,
		restyClient:    resty.New(),
		headers:        make(map[string]string),
	}
}
