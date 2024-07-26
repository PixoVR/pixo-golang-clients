package heartbeat

import (
	abstract "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
)

type Client interface {
	abstract.AbstractClient
	SendPulse(sessionID int) error
}

// Client is a struct that contains an AbstractServiceClient
type client struct {
	abstract.AbstractServiceClient
	platformClient platform.Client
}

// NewClient is a function that returns a new Client
func NewClient(config urlfinder.ClientConfig) Client {

	abstractConfig := abstract.AbstractConfig{
		ServiceConfig: newServiceConfig(config.Lifecycle, config.Region),
		Token:         config.Token,
	}

	return &client{
		AbstractServiceClient: *abstract.NewClient(abstractConfig),
		platformClient:        platform.NewClient(config),
	}
}

func NewClientWithBasicAuth(username, password string, config urlfinder.ClientConfig) (Client, error) {
	platformClient, err := platform.NewClientWithBasicAuth(username, password, config)
	if err != nil {
		return nil, err
	}

	abstractConfig := abstract.AbstractConfig{
		ServiceConfig: newServiceConfig(config.Lifecycle, config.Region),
		Token:         platformClient.GetToken(),
	}

	return &client{
		AbstractServiceClient: *abstract.NewClient(abstractConfig),
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
		Service:   "heartbeat",
		Lifecycle: lifecycle,
		Region:    region,
		Port:      8002,
	}
}
