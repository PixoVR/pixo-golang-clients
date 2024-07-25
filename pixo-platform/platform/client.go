package platform

import (
	"context"
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/legacy"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/middleware/auth"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"

	abstract "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	"github.com/hasura/go-graphql-client"
)

type Client interface {
	abstract.AbstractClient

	ActiveUserID() int
	ActiveOrgID() int

	GetUserByUsername(ctx context.Context, username string) (*legacy.User, error)
	CreateUser(ctx context.Context, user legacy.User) (*legacy.User, error)
	UpdateUser(ctx context.Context, user legacy.User) (*legacy.User, error)
	DeleteUser(ctx context.Context, id int) error

	GetAPIKeys(ctx context.Context, params *APIKeyQueryParams) ([]APIKey, error)
	CreateAPIKey(ctx context.Context, input APIKey) (*APIKey, error)
	DeleteAPIKey(ctx context.Context, id int) error

	GetWebhooks(ctx context.Context, params *WebhookParams) ([]Webhook, error)
	GetWebhook(ctx context.Context, id int) (*Webhook, error)
	CreateWebhook(ctx context.Context, input Webhook) (*Webhook, error)
	UpdateWebhook(ctx context.Context, input Webhook) (*Webhook, error)
	DeleteWebhook(ctx context.Context, id int) error

	GetPlatforms(ctx context.Context) ([]*Platform, error)
	GetControlTypes(ctx context.Context) ([]*ControlType, error)

	CreateModuleVersion(ctx context.Context, input ModuleVersion) (*ModuleVersion, error)

	GetOrg(ctx context.Context, id int) (*Org, error)
	CreateOrg(ctx context.Context, org Org) (*Org, error)
	UpdateOrg(ctx context.Context, org Org) (*Org, error)
	DeleteOrg(ctx context.Context, id int) error

	GetSession(ctx context.Context, id int) (*Session, error)
	CreateSession(ctx context.Context, moduleID int, ipAddress, deviceId string) (*Session, error)
	UpdateSession(ctx context.Context, session Session) (*Session, error)
	CreateEvent(ctx context.Context, event Event) (*Event, error)

	GetMultiplayerServerConfigs(ctx context.Context, params *MultiplayerServerConfigParams) ([]*MultiplayerServerConfigQueryParams, error)
	GetMultiplayerServerVersions(ctx context.Context, params *MultiplayerServerVersionQueryParams) ([]*MultiplayerServerVersion, error)
	GetMultiplayerServerVersion(ctx context.Context, id int) (*MultiplayerServerVersion, error)
	CreateMultiplayerServerVersion(ctx context.Context, input MultiplayerServerVersion) (*MultiplayerServerVersion, error)
}

var _ Client = (*PlatformClient)(nil)

// PlatformClient is a struct for the graphql API that contains an abstract client
type PlatformClient struct {
	*abstract.AbstractServiceClient
	*graphql.Client
	defaultContext context.Context
}

// NewClient is a function that returns a PlatformClient
func NewClient(config urlfinder.ClientConfig) *PlatformClient {

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

	abstractConfig := abstract.AbstractConfig{
		ServiceConfig: serviceConfig,
		Token:         config.Token,
		APIKey:        config.APIKey,
	}

	return &PlatformClient{
		AbstractServiceClient: abstract.NewClient(abstractConfig),
		Client:                graphql.NewClient(fmt.Sprintf("%s/query", url), &c),
		defaultContext:        context.Background(),
	}
}

// NewClientWithBasicAuth is a function that returns a PlatformClient with basic auth performed
func NewClientWithBasicAuth(username, password string, config urlfinder.ClientConfig) (*PlatformClient, error) {

	serviceConfig := newServiceConfig(config)

	client := &PlatformClient{
		AbstractServiceClient: abstract.NewClient(abstract.AbstractConfig{ServiceConfig: serviceConfig}),
		defaultContext:        context.Background(),
	}

	if err := client.Login(username, password); err != nil {
		log.Error().Err(err).Msg("Failed to login to the pixo platform")
		return nil, err
	}

	httpClient := http.Client{Transport: &transport{underlyingTransport: http.DefaultTransport, token: client.GetToken()}}

	client.Client = graphql.NewClient(fmt.Sprintf("%s/query", serviceConfig.FormatURL()), &httpClient)

	return client, nil
}

func (g *PlatformClient) SetToken(token string) {
	g.AbstractServiceClient.SetToken(token)
	httpClient := http.Client{Transport: &transport{underlyingTransport: http.DefaultTransport, token: token}}
	g.Client = graphql.NewClient(fmt.Sprintf("%s/query", g.GetURL()), &httpClient)
}

func (g *PlatformClient) SetAPIKey(key string) {
	g.AbstractServiceClient.SetAPIKey(key)
}

func (g *PlatformClient) ActiveUserID() int {

	if !g.IsAuthenticated() {
		return 0
	}

	token := g.GetToken()

	rawToken, err := auth.ParseJWT(token)
	if err != nil {
		return 0
	}

	return rawToken.UserID
}

func (g *PlatformClient) ActiveOrgID() int {

	if !g.IsAuthenticated() {
		return 0
	}

	token := g.GetToken()

	rawToken, err := auth.ParseJWT(token)
	if err != nil {
		return 0
	}

	return rawToken.OrgID
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
