package graphql_api

import (
	"context"
	"fmt"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"

	abstract_client "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	"github.com/hasura/go-graphql-client"
)

type PlatformClient interface {
	GetUserByUsername(ctx context.Context, username string) (*platform.User, error)
	CreateUser(ctx context.Context, user platform.User) (*platform.User, error)
	UpdateUser(ctx context.Context, user platform.User) (*platform.User, error)
	DeleteUser(ctx context.Context, id int) error

	GetAPIKeys(ctx context.Context, params *APIKeyQueryParams) ([]*platform.APIKey, error)
	CreateAPIKey(ctx context.Context, input platform.APIKey) (*platform.APIKey, error)
	DeleteAPIKey(ctx context.Context, id int) error

	GetSession(ctx context.Context, id int) (*Session, error)
	CreateSession(ctx context.Context, moduleID int, ipAddress, deviceId string) (*Session, error)
	UpdateSession(ctx context.Context, id int, status string, completed bool) (*Session, error)
	CreateEvent(ctx context.Context, sessionID int, uuid string, eventType string, data string) (*platform.Event, error)
}

var _ PlatformClient = (*GraphQLAPIClient)(nil)

// GraphQLAPIClient is a struct for the graphql API that contains an abstract client
type GraphQLAPIClient struct {
	*abstract_client.PixoAbstractAPIClient
	*graphql.Client
	underlyingTransport http.RoundTripper
	defaultContext      context.Context
}

// NewClient is a function that returns a GraphQLAPIClient
func NewClient(config urlfinder.ClientConfig) *GraphQLAPIClient {

	if config.Token == "" && config.APIKey == "" {
		config.Token = os.Getenv("SECRET_KEY")
	}

	serviceConfig := newServiceConfig(config)

	url := serviceConfig.FormatURL()

	t := &transport{
		underlyingTransport: http.DefaultTransport,
		token:               config.Token,
		key:                 config.APIKey,
	}
	c := http.Client{Transport: t}

	abstractConfig := abstract_client.AbstractConfig{
		Token:  config.Token,
		APIKey: config.APIKey,
		URL:    url,
	}

	return &GraphQLAPIClient{
		PixoAbstractAPIClient: abstract_client.NewClient(abstractConfig),
		Client:                graphql.NewClient(fmt.Sprintf("%s/query", url), &c),
		defaultContext:        context.Background(),
	}
}

// NewClientWithBasicAuth is a function that returns a GraphQLAPIClient with basic auth performed
func NewClientWithBasicAuth(username, password string, config urlfinder.ClientConfig) (*GraphQLAPIClient, error) {

	serviceConfig := newServiceConfig(config)

	url := serviceConfig.FormatURL()

	abstractConfig := abstract_client.AbstractConfig{URL: url}

	client := &GraphQLAPIClient{
		PixoAbstractAPIClient: abstract_client.NewClient(abstractConfig),
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
