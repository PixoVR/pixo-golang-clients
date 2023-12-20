package graphql_api

import (
	"context"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
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

	url := getURL(config.FormatURL())

	c := http.Client{Transport: &transport{underlyingTransport: http.DefaultTransport, token: token}}

	return &GraphQLAPIClient{
		PixoAbstractAPIClient: *abstract_client.NewClient(token, url),
		gqlClient:             graphql.NewClient(url, &c),
		defaultContext:        context.Background(),
	}
}

func newServiceConfig(lifecycle, region string) urlfinder.ServiceConfig {
	return urlfinder.ServiceConfig{
		Service:   "primary",
		Lifecycle: lifecycle,
		Region:    region,
		Port:      8000,
	}
}
