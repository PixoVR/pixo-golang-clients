package heartbeat

import (
	"context"
	abstract "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
)

// Client is an interface that contains an AbstractClient and a SendPulse method
type Client interface {
	abstract.AbstractClient
	SendPulse(ctx context.Context, sessionID int) error
	SendPulsesWithCancel(ctx context.Context, sessionID int, periodSeconds float64) (chan error, context.CancelFunc)
}

var _ Client = &client{}

// Client is a struct that contains an ServiceClient
type client struct {
	abstract.ServiceClient
	platformClient platform.Client
}

// NewClient is a function that returns a new Client
func NewClient(config urlfinder.ClientConfig) Client {

	abstractConfig := abstract.Config{
		ServiceConfig: newServiceConfig(config.Lifecycle, config.Region),
		Token:         config.Token,
	}

	return &client{
		ServiceClient:  *abstract.NewClient(abstractConfig),
		platformClient: platform.NewClient(config),
	}
}

// NewClientWithBasicAuth is a function that returns a new Client using basic auth for the platform client
func NewClientWithBasicAuth(username, password string, config urlfinder.ClientConfig) (Client, error) {
	platformClient, err := platform.NewClientWithBasicAuth(username, password, config)
	if err != nil {
		return nil, err
	}

	abstractConfig := abstract.Config{
		ServiceConfig: newServiceConfig(config.Lifecycle, config.Region),
		Token:         platformClient.GetToken(),
	}

	return &client{
		ServiceClient: *abstract.NewClient(abstractConfig),
	}, nil
}

// Login is a method that logs in the client using the platform client
func (c *client) Login(username, password string) error {
	if err := c.platformClient.Login(username, password); err != nil {
		return err
	}

	c.SetToken(c.platformClient.GetToken())
	return nil
}

// ActiveUserID returns the active user id
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
