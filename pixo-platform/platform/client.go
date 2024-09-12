package platform

import (
	"context"
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/middleware/auth"
	"github.com/rs/zerolog/log"
	"os"

	abstract "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract-client"
	"github.com/hasura/go-graphql-client"
)

var _ Client = (*clientImpl)(nil)

// clientImpl is a struct for the graphql API that contains an abstract client
type clientImpl struct {
	*abstract.AbstractServiceClient
	*graphql.Client
	defaultContext context.Context
}

// NewClient is a function that returns a clientImpl
func NewClient(config urlfinder.ClientConfig) Client {

	if config.Token == "" && config.APIKey == "" {
		config.APIKey = os.Getenv("PIXO_API_KEY")
	}

	serviceConfig := newServiceConfig(config)

	url := serviceConfig.FormatURL()

	abstractConfig := abstract.AbstractConfig{
		ServiceConfig: serviceConfig,
		Token:         config.Token,
		APIKey:        config.APIKey,
	}
	abstractClient := abstract.NewClient(abstractConfig)

	return &clientImpl{
		AbstractServiceClient: abstractClient,
		Client:                graphql.NewClient(fmt.Sprintf("%s/query", url), abstractClient.Client()),
		defaultContext:        context.Background(),
	}
}

// NewClientWithBasicAuth is a function that returns a clientImpl with basic auth performed
func NewClientWithBasicAuth(username, password string, config urlfinder.ClientConfig) (Client, error) {

	client := NewClient(config)

	if err := client.Login(username, password); err != nil {
		log.Error().Err(err).Msg("Failed to login to the pixo platform")
		return nil, err
	}

	return client, nil
}

func (p *clientImpl) ActiveUserID() int {
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

func (p *clientImpl) ActiveOrgID() int {
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
