package primary_api

import (
	abstractClient "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	"github.com/rs/zerolog/log"
)

// PrimaryAPIClient is a struct for the primary API that contains an abstract client
type PrimaryAPIClient struct {
	abstractClient.PixoAbstractAPIClient
}

// NewClient is a function that returns a PixoAbstractAPIClient
func NewClient(token, apiURL string) *PrimaryAPIClient {

	if apiURL == "" {
		apiURL = getURL()
	}

	return &PrimaryAPIClient{
		PixoAbstractAPIClient: *abstractClient.NewClient(token, apiURL),
	}
}

// NewClientWithBasicAuth is a function that returns a PixoAbstractAPIClient with basic auth performed
func NewClientWithBasicAuth(username, password, apiURL string) *PrimaryAPIClient {
	if apiURL == "" {
		apiURL = getURL()
	}

	primaryClient := &PrimaryAPIClient{
		PixoAbstractAPIClient: *abstractClient.NewClient("", apiURL),
	}

	if err := primaryClient.Login(username, password); err != nil {
		log.Error().Err(err).Msg("Failed to login")
		return nil
	}

	return primaryClient
}
