package graphql_api

import (
	"context"
	"fmt"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/middleware/auth"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"

	abstract_client "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	"github.com/hasura/go-graphql-client"
)

type PlatformClient interface {
	abstract_client.AbstractClient

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

	GetMultiplayerServerConfigs(ctx context.Context, params *MultiplayerServerConfigParams) ([]*MultiplayerServerConfigQueryParams, error)
	GetMultiplayerServerVersions(ctx context.Context, params *MultiplayerServerVersionQueryParams) ([]*MultiplayerServerVersion, error)
	GetMultiplayerServerVersion(ctx context.Context, id int) (*MultiplayerServerVersion, error)
	CreateMultiplayerServerVersion(ctx context.Context, input MultiplayerServerVersion) (*MultiplayerServerVersion, error)
}

var _ PlatformClient = (*GraphQLAPIClient)(nil)

// GraphQLAPIClient is a struct for the graphql API that contains an abstract client
type GraphQLAPIClient struct {
	*abstract_client.AbstractServiceClient
	*graphql.Client
	underlyingTransport http.RoundTripper
	defaultContext      context.Context
}

// NewClient is a function that returns a GraphQLAPIClient
func NewClient(config urlfinder.ClientConfig) *GraphQLAPIClient {

	if config.APIKey == "" {
		config.APIKey = os.Getenv("PIXO_API_KEY")
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
		AbstractServiceClient: abstract_client.NewClient(abstractConfig),
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
		AbstractServiceClient: abstract_client.NewClient(abstractConfig),
		defaultContext:        context.Background(),
	}

	if err := client.Login(username, password); err != nil {
		log.Error().Err(err).Msg("Failed to login to the pixo platform")
		return nil, err
	}

	httpClient := http.Client{Transport: &transport{underlyingTransport: http.DefaultTransport, token: client.GetToken()}}

	client.Client = graphql.NewClient(fmt.Sprintf("%s/query", url), &httpClient)

	return client, nil
}

func (g *GraphQLAPIClient) SetToken(token string) {
	g.AbstractServiceClient.SetToken(token)
	httpClient := http.Client{Transport: &transport{underlyingTransport: http.DefaultTransport, token: token}}
	g.Client = graphql.NewClient(fmt.Sprintf("%s/query", g.GetURL()), &httpClient)
}

func (g *GraphQLAPIClient) SetAPIKey(key string) {
	g.AbstractServiceClient.SetAPIKey(key)
}

func (g *GraphQLAPIClient) ActiveUserID() int {

	if !g.IsAuthenticated() {
		return 0
	}

	token := g.GetToken()

	rawToken, err := auth.ParseJWT(token)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse JWT")
		return 0
	}

	return rawToken.UserID
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
