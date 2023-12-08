package multiplayer_allocator

import (
	abstractClient "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
)

// AllocatorClient is a struct for the primary API that contains an abstract client
type AllocatorClient struct {
	abstractClient.PixoAbstractAPIClient
}

// NewClient is a function that returns a PixoAbstractAPIClient
func NewClient(token, lifecycle, region string) *AllocatorClient {

	config := newServiceConfig(lifecycle, region)

	return &AllocatorClient{
		PixoAbstractAPIClient: *abstractClient.NewClient(token, config.FormatURL()),
	}
}

func newServiceConfig(lifecycle, region string) urlfinder.ServiceConfig {
	return urlfinder.ServiceConfig{
		Service:   "allocator",
		Tenant:    "multiplayer",
		Lifecycle: lifecycle,
		Region:    region,
	}
}
