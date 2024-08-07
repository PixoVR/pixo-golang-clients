package headset

import (
	abstractClient "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	primary_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/legacy"
	graphql_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
)

type Client interface {
	abstractClient.AbstractClient
}

// Client is a struct that contains an AbstractServiceClient
type client struct {
	abstractClient.AbstractServiceClient
	platformClient graphql_api.Client
}

// NewClient is a function that returns a new Client
func NewClient(config urlfinder.ClientConfig) Client {

	abstractConfig := abstractClient.AbstractConfig{
		ServiceConfig: newServiceConfig(config.Lifecycle, config.Region),
		Token:         config.Token,
	}

	return &client{
		AbstractServiceClient: *abstractClient.NewClient(abstractConfig),
		platformClient:        graphql_api.NewClient(config),
	}
}

func NewClientWithBasicAuth(username, password string, config urlfinder.ClientConfig) (Client, error) {
	primaryClient, err := primary_api.NewClientWithBasicAuth(username, password, config)
	if err != nil {
		return nil, err
	}

	abstractConfig := abstractClient.AbstractConfig{
		ServiceConfig: newServiceConfig(config.Lifecycle, config.Region),
		Token:         primaryClient.GetToken(),
	}

	return &client{
		AbstractServiceClient: *abstractClient.NewClient(abstractConfig),
	}, nil
}

func (c *client) Login(username, password string) error {
	if err := c.platformClient.Login(username, password); err != nil {
		return err
	}

	c.SetToken(c.platformClient.GetToken())
	return nil
}

func (c *client) ActiveUserID() int {
	return c.platformClient.ActiveUserID()
}

func newServiceConfig(lifecycle, region string) urlfinder.ServiceConfig {
	return urlfinder.ServiceConfig{
		Service:   "modules",
		Lifecycle: lifecycle,
		Region:    region,
		Port:      8001,
	}
}
