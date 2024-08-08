package headset

import (
	"context"
	abstractClient "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
)

type Client interface {
	abstractClient.AbstractClient
	StartSession(ctx context.Context, request EventRequest) (*EventResponse, error)
	SendEvent(ctx context.Context, request EventRequest) (*EventResponse, error)
	EndSession(ctx context.Context, request EventRequest) (*EventResponse, error)
}

// Client is a struct that contains an AbstractServiceClient
type client struct {
	abstractClient.AbstractServiceClient
	platformClient platform.Client
}

// NewClient is a function that returns a new Client
func NewClient(config urlfinder.ClientConfig) Client {

	abstractConfig := abstractClient.AbstractConfig{
		ServiceConfig: newServiceConfig(config.Lifecycle, config.Region),
		Token:         config.Token,
	}

	return &client{
		AbstractServiceClient: *abstractClient.NewClient(abstractConfig),
		platformClient:        platform.NewClient(config),
	}
}

func NewClientWithBasicAuth(username, password string, config urlfinder.ClientConfig) (Client, error) {
	abstractConfig := abstractClient.AbstractConfig{
		ServiceConfig: newServiceConfig(config.Lifecycle, config.Region),
	}

	c := &client{
		AbstractServiceClient: *abstractClient.NewClient(abstractConfig),
	}

	if err := c.Login(username, password); err != nil {
		return nil, err
	}

	return c, nil
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
