package allocator

import (
	abstractClient "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	primary_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
)

// AllocatorClient is a struct for the primary API that contains an abstract client
type AllocatorClient struct {
	abstractClient.AbstractServiceClient
}

// NewClient is a function that returns a AbstractServiceClient
func NewClient(config urlfinder.ClientConfig) *AllocatorClient {

	serviceConfig := newServiceConfig(config.Lifecycle, config.Region)

	abstractConfig := abstractClient.AbstractConfig{
		URL:   serviceConfig.FormatURL(),
		Token: config.Token,
	}

	return &AllocatorClient{
		AbstractServiceClient: *abstractClient.NewClient(abstractConfig),
	}
}

func NewClientWithBasicAuth(username, password string, config urlfinder.ClientConfig) (*AllocatorClient, error) {
	primaryClient, err := primary_api.NewClientWithBasicAuth(username, password, config)
	if err != nil {
		return nil, err
	}

	serviceConfig := newServiceConfig(config.Lifecycle, config.Region)

	abstractConfig := abstractClient.AbstractConfig{
		Path:  serviceConfig.Service,
		URL:   serviceConfig.FormatURL(),
		Token: primaryClient.GetToken(),
	}

	return &AllocatorClient{
		AbstractServiceClient: *abstractClient.NewClient(abstractConfig),
	}, nil
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
