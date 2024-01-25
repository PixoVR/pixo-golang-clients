package primary_api

import (
	abstractClient "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/rs/zerolog/log"
)

// PrimaryAPIClient is a struct for the primary API that contains an abstract client
type PrimaryAPIClient struct {
	abstractClient.PixoAbstractAPIClient
}

// NewClient is a function that returns a PixoAbstractAPIClient
func NewClient(token, lifecycle, region string) *PrimaryAPIClient {

	config := newServiceConfig(lifecycle, region)

	abstractConfig := abstractClient.AbstractConfig{
		Token: token,
		URL:   config.FormatURL(),
	}

	return &PrimaryAPIClient{
		PixoAbstractAPIClient: *abstractClient.NewClient(abstractConfig),
	}
}

// NewClientWithBasicAuth is a function that returns a PixoAbstractAPIClient with basic auth performed
func NewClientWithBasicAuth(username, password, lifecycle, region string) (*PrimaryAPIClient, error) {

	config := newServiceConfig(lifecycle, region)

	abstractConfig := abstractClient.AbstractConfig{
		URL: config.FormatURL(),
	}

	primaryClient := &PrimaryAPIClient{
		PixoAbstractAPIClient: *abstractClient.NewClient(abstractConfig),
	}

	if err := primaryClient.Login(username, password); err != nil {
		log.Error().Err(err).Msg("Failed to login to the pixo platform")
		return nil, err
	}

	return primaryClient, nil
}

func newServiceConfig(lifecycle, region string) urlfinder.ServiceConfig {
	return urlfinder.ServiceConfig{
		Service:   "api",
		Lifecycle: lifecycle,
		Region:    region,
		Port:      3001,
	}
}
