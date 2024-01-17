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
	abstract_client.PixoAbstractAPIClient
	underlyingTransport http.RoundTripper
	gqlClient           *graphql.Client
	defaultContext      context.Context
}

// NewClient is a function that returns a GraphQLAPIClient
func NewClient(token, lifecycle, region string) *GraphQLAPIClient {

	if token == "" {
		token = os.Getenv("SECRET_KEY")
	}

	config := newServiceConfig(lifecycle, region)

	url := config.FormatURL()

	c := http.Client{Transport: &transport{underlyingTransport: http.DefaultTransport, token: token}}

	return &GraphQLAPIClient{
		PixoAbstractAPIClient: *abstract_client.NewClient(token, url),
		gqlClient:             graphql.NewClient(fmt.Sprintf("%s/query", url), &c),
		defaultContext:        context.Background(),
	}
}

// NewClientWithBasicAuth is a function that returns a GraphQLAPIClient with basic auth performed
func NewClientWithBasicAuth(username, password, lifecycle, region string) (*GraphQLAPIClient, error) {

	config := newServiceConfig(lifecycle, region)

	url := config.FormatURL()

	client := &GraphQLAPIClient{
		PixoAbstractAPIClient: *abstract_client.NewClient("", config.FormatURL()),
		defaultContext:        context.Background(),
	}

	if err := client.Login(username, password); err != nil {
		log.Error().Err(err).Msg("Failed to login to the pixo platform")
		return nil, err
	}

	c := http.Client{Transport: &transport{underlyingTransport: http.DefaultTransport, token: client.GetToken()}}

	client.gqlClient = graphql.NewClient(fmt.Sprintf("%s/query", url), &c)

	return client, nil
}

func newServiceConfig(lifecycle, region string) urlfinder.ServiceConfig {
	return urlfinder.ServiceConfig{
		Service:   "v2",
		Lifecycle: lifecycle,
		Region:    region,
		Port:      8000,
	}
}
