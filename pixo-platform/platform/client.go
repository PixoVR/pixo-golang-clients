package platform

import (
	"context"
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/middleware/auth"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"

	abstract "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	"github.com/hasura/go-graphql-client"
)

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

func (p *PlatformClient) SetToken(token string) {
	p.AbstractServiceClient.SetToken(token)
	httpClient := http.Client{Transport: &transport{underlyingTransport: http.DefaultTransport, token: token}}
	p.Client = graphql.NewClient(fmt.Sprintf("%s/query", p.GetURL()), &httpClient)
}

func (p *PlatformClient) SetAPIKey(key string) {
	p.AbstractServiceClient.SetAPIKey(key)
}

func (p *PlatformClient) ActiveUserID() int {

	if !p.IsAuthenticated() {
		return 0
	}

	token := p.GetToken()

	rawToken, err := auth.ParseJWT(token)
	if err != nil {
		return 0
	}

	return rawToken.UserID
}

func (p *PlatformClient) ActiveOrgID() int {

	if !p.IsAuthenticated() {
		return 0
	}

	token := p.GetToken()

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
