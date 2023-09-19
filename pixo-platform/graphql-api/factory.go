package graphql_api

import (
	"fmt"
	abstract_client "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	"github.com/hasura/go-graphql-client"
	"net/http"
)

// GraphQLAPIClient is a struct for the graphql API that contains an abstract client
type GraphQLAPIClient struct {
	abstract_client.PixoAbstractAPIClient
	underlyingTransport http.RoundTripper
	gqlClient           *graphql.Client
}

// NewClient is a function that returns a PixoAbstractAPIClient
func NewClient(token, apiURL string) *GraphQLAPIClient {

	if apiURL == "" {
		apiURL = getURL()
	}

	url := fmt.Sprintf("%s/v2/query", apiURL)

	c := http.Client{Transport: &transport{underlyingTransport: http.DefaultTransport, token: token}}

	return &GraphQLAPIClient{
		PixoAbstractAPIClient: *abstract_client.NewClient(token, apiURL),
		gqlClient:             graphql.NewClient(url, &c),
	}
}
