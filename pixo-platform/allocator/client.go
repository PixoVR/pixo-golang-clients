package allocator

import (
	abstractClient "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
)

// Client is a struct for the primary API that contains an abstract client
type Client struct {
	abstractClient.ServiceClient

	platformClient platform.Client
}

// NewClientWithBasicAuth is a function that returns a ServiceClient with basic auth
func NewClientWithBasicAuth(username, password string, config urlfinder.ClientConfig) (*Client, error) {
	platformClient, err := platform.NewClientWithBasicAuth(username, password, config)
	if err != nil {
		return nil, err
	}

	config.Token = platformClient.GetToken()

	return NewClient(config), nil
}

// NewClient is a function that returns a ServiceClient
func NewClient(config urlfinder.ClientConfig) *Client {

	serviceConfig := newServiceConfig(config.Lifecycle, config.Region)

	abstractConfig := abstractClient.AbstractConfig{
		ServiceConfig: serviceConfig,
		APIKey:        config.APIKey,
		Token:         config.Token,
	}

	return &Client{
		ServiceClient:  *abstractClient.NewClient(abstractConfig),
		platformClient: platform.NewClient(config),
	}
}

// Login send a login request to the platform
func (a *Client) Login(username, password string) error {
	return a.platformClient.Login(username, password)
}

func newServiceConfig(lifecycle, region string) urlfinder.ServiceConfig {
	return urlfinder.ServiceConfig{
		Service:   "allocator",
		Tenant:    "multiplayer",
		Lifecycle: lifecycle,
		Region:    region,
		Port:      8003,
	}
}
