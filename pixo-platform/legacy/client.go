package legacy

import (
	abstract "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
)

type LegacyClient interface {
	Login(username, password string) error
	GetOrgs() ([]Org, error)
	CreateWebhook(webhook Webhook) error
	GetWebhooks(orgID int) ([]Webhook, error)
	DeleteWebhook(webhookID int) error
}

// Client is a struct for the primary API that contains an abstract client
type Client struct {
	abstract.AbstractServiceClient
}

// NewClient is a function that returns a AbstractServiceClient
func NewClient(config urlfinder.ClientConfig) *Client {

	serviceConfig := newServiceConfig(config.Lifecycle, config.Region)

	abstractConfig := abstract.AbstractConfig{
		ServiceConfig: serviceConfig,
		Token:         config.Token,
	}

	return &Client{
		AbstractServiceClient: *abstract.NewClient(abstractConfig),
	}
}

// NewClientWithBasicAuth is a function that returns a AbstractServiceClient with basic auth performed
func NewClientWithBasicAuth(username, password string, config urlfinder.ClientConfig) (*Client, error) {

	serviceConfig := newServiceConfig(config.Lifecycle, config.Region)

	abstractConfig := abstract.AbstractConfig{
		ServiceConfig: serviceConfig,
	}

	primaryClient := &Client{
		AbstractServiceClient: *abstract.NewClient(abstractConfig),
	}

	if err := primaryClient.Login(username, password); err != nil {
		return nil, err
	}

	return primaryClient, nil
}

func (p *Client) ActiveUserID() int {
	return 0
}

func newServiceConfig(lifecycle, region string) urlfinder.ServiceConfig {
	return urlfinder.ServiceConfig{
		Service:   "api",
		Lifecycle: lifecycle,
		Region:    region,
		Port:      3001,
	}
}
