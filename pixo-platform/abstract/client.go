package abstract

import (
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/gorilla/websocket"
	"sync"
)

// ServiceClient is a struct that handles generic http client operations
type ServiceClient struct {
	serviceConfig urlfinder.ServiceConfig
	headers       map[string]string
	token         string
	key           string

	websocketConn  *websocket.Conn
	timeoutSeconds int

	lock sync.Mutex
}

// Config is a struct that holds the configuration for the ServiceClient
type Config struct {
	ServiceConfig  urlfinder.ServiceConfig
	APIKey         string
	Token          string
	TimeoutSeconds int
}

// NewClient creates a new ServiceClient given a config struct
func NewClient(config Config) *ServiceClient {

	if config.TimeoutSeconds == 0 {
		config.TimeoutSeconds = 30
	}

	abstractClient := ServiceClient{
		serviceConfig:  config.ServiceConfig,
		key:            config.APIKey,
		token:          config.Token,
		timeoutSeconds: config.TimeoutSeconds,
		headers:        make(map[string]string),
	}

	if config.Token != "" {
		abstractClient.SetToken(config.Token)
	} else if config.APIKey != "" {
		abstractClient.SetAPIKey(config.APIKey)
	}

	return &abstractClient
}
