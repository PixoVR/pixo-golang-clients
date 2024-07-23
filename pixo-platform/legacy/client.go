package legacy

import (
	abstractClient "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
)

type OldAPIClient interface {
	Login(username, password string) error
	CreateWebhook(webhook Webhook) error
	GetWebhooks(orgID int) ([]Webhook, error)
	DeleteWebhook(webhookID int) error
}

// LegacyAPIClient is a struct for the primary API that contains an abstract client
type LegacyAPIClient struct {
	abstractClient.AbstractServiceClient
}

// NewClient is a function that returns a AbstractServiceClient
func NewClient(config urlfinder.ClientConfig) *LegacyAPIClient {

	serviceConfig := newServiceConfig(config.Lifecycle, config.Region)

	abstractConfig := abstractClient.AbstractConfig{
		ServiceConfig: serviceConfig,
		Token:         config.Token,
	}

	return &LegacyAPIClient{
		AbstractServiceClient: *abstractClient.NewClient(abstractConfig),
	}
}

// NewClientWithBasicAuth is a function that returns a AbstractServiceClient with basic auth performed
func NewClientWithBasicAuth(username, password string, config urlfinder.ClientConfig) (*LegacyAPIClient, error) {

	serviceConfig := newServiceConfig(config.Lifecycle, config.Region)

	abstractConfig := abstractClient.AbstractConfig{
		ServiceConfig: serviceConfig,
	}

	primaryClient := &LegacyAPIClient{
		AbstractServiceClient: *abstractClient.NewClient(abstractConfig),
	}

	if err := primaryClient.Login(username, password); err != nil {
		return nil, err
	}

	return primaryClient, nil
}

func (p *LegacyAPIClient) ActiveUserID() int {
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
