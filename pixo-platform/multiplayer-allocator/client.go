package multiplayer_allocator

import (
	abstractClient "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	primary_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
)

// AllocatorClient is a struct for the primary API that contains an abstract client
type AllocatorClient struct {
	abstractClient.PixoAbstractAPIClient
}

// NewClient is a function that returns a PixoAbstractAPIClient
func NewClient(token, lifecycle, region string) *AllocatorClient {

	config := newServiceConfig(lifecycle, region)

	abstractConfig := abstractClient.AbstractConfig{
		URL:   config.FormatURL(),
		Token: token,
	}

	return &AllocatorClient{
		PixoAbstractAPIClient: *abstractClient.NewClient(abstractConfig),
	}
}

func NewClientWithBasicAuth(username, password, lifecycle, region string) (*AllocatorClient, error) {
	primaryClient, err := primary_api.NewClientWithBasicAuth(username, password, lifecycle, region)
	if err != nil {
		return nil, err
	}

	config := newServiceConfig(lifecycle, region)

	abstractConfig := abstractClient.AbstractConfig{
		URL:   config.FormatURL(),
		Token: primaryClient.GetToken(),
	}

	return &AllocatorClient{
		PixoAbstractAPIClient: *abstractClient.NewClient(abstractConfig),
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
