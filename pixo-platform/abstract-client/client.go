package abstract_client

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
	"sync"
)

// AbstractServiceClient is a struct that handles generic http client operations
type AbstractServiceClient struct {
	serviceConfig urlfinder.ServiceConfig
	token         string
	key           string

	client         *resty.Client
	websocketConn  *websocket.Conn
	timeoutSeconds int
	lock           sync.Mutex
	headers        sync.Map
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

	client := resty.New().SetHeader("Content-Type", "application/json")

	abstractClient := AbstractServiceClient{
		serviceConfig:  config.ServiceConfig,
		key:            config.APIKey,
		token:          config.Token,
		timeoutSeconds: config.TimeoutSeconds,
		client:         client,
	}

	if config.Token != "" {
		abstractClient.SetToken(config.Token)
	} else if config.APIKey != "" {
		abstractClient.SetAPIKey(config.APIKey)
	}

	return &abstractClient
}
