package graphql_api

import (
	"context"
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"

	abstract_client "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	"github.com/hasura/go-graphql-client"
)

// GraphQLAPIClient is a struct for the graphql API that contains an abstract client
type GraphQLAPIClient struct {
	*abstract_client.PixoAbstractAPIClient
	*graphql.Client
	underlyingTransport http.RoundTripper
	defaultContext      context.Context
}

// NewClient is a function that returns a GraphQLAPIClient
func NewClient(config urlfinder.ClientConfig) *GraphQLAPIClient {

	if config.Token == "" {
		config.Token = os.Getenv("SECRET_KEY")
	}

	serviceConfig := newServiceConfig(config)

	url := serviceConfig.FormatURL()

	c := http.Client{Transport: &transport{underlyingTransport: http.DefaultTransport, token: config.Token}}

	return &GraphQLAPIClient{
		PixoAbstractAPIClient: abstract_client.NewClient(config.Token, url),
		Client:                graphql.NewClient(fmt.Sprintf("%s/query", url), &c),
		defaultContext:        context.Background(),
	}
}

// NewClientWithBasicAuth is a function that returns a GraphQLAPIClient with basic auth performed
func NewClientWithBasicAuth(username, password string, config urlfinder.ClientConfig) (*GraphQLAPIClient, error) {

	serviceConfig := newServiceConfig(config)

	url := serviceConfig.FormatURL()

	client := &GraphQLAPIClient{
		PixoAbstractAPIClient: abstract_client.NewClient("", serviceConfig.FormatURL()),
		defaultContext:        context.Background(),
	}

	if err := client.Login(username, password); err != nil {
		log.Error().Err(err).Msg("Failed to login to the pixo platform")
		return nil, err
	}

	c := http.Client{Transport: &transport{underlyingTransport: http.DefaultTransport, token: client.GetToken()}}

	client.Client = graphql.NewClient(fmt.Sprintf("%s/query", url), &c)

	return client, nil
}

func newServiceConfig(config urlfinder.ClientConfig) urlfinder.ServiceConfig {
	return urlfinder.ServiceConfig{
		Service:     "v2",
		ServiceName: "primary-api",
		Lifecycle:   config.Lifecycle,
		Region:      config.Region,
		Namespace:   fmt.Sprintf("%s-apex", config.Lifecycle),
		Port:        8000,
	}
}
