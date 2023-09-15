package multiplayer_allocator

import (
	abstractClient "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
)

// AllocatorClient is a struct for the primary API that contains an abstract client
type AllocatorClient struct {
	abstractClient.PixoAbstractAPIClient
}

// NewClient is a function that returns a PixoAbstractAPIClient
func NewClient(token, apiURL string) *AllocatorClient {

	if apiURL == "" {
		apiURL = getURL()
	}

	return &AllocatorClient{
		PixoAbstractAPIClient: *abstractClient.NewClient(token, apiURL),
	}
}
