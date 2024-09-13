package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/middleware/auth"
	"github.com/rs/zerolog/log"
	"io"
	"os"

	abstract "github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract"
)

var _ Client = (*clientImpl)(nil)

// clientImpl is a struct for the graphql API that contains an abstract client
type clientImpl struct {
	*abstract.ServiceClient
	defaultContext context.Context
}

// NewClient is a function that returns a clientImpl
func NewClient(config urlfinder.ClientConfig) Client {

	if config.Token == "" && config.APIKey == "" {
		config.APIKey = os.Getenv("PIXO_API_KEY")
	}

	abstractConfig := abstract.AbstractConfig{
		ServiceConfig: newServiceConfig(config),
		Token:         config.Token,
		APIKey:        config.APIKey,
	}
	abstractClient := abstract.NewClient(abstractConfig)

	return &clientImpl{
		ServiceClient:  abstractClient,
		defaultContext: context.Background(),
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

func (p *clientImpl) Exec(ctx context.Context, query string, v any, variables map[string]interface{}) error {
	req := GraphQLRequestPayload{
		Query:     query,
		Variables: variables,
	}
	reqBody, _ := json.Marshal(req)
	p.SetHeader("Content-Type", "application/json")
	res, err := p.Post(ctx, "query", reqBody)
	if err != nil {
		return err
	}

	var gqlRes struct {
		Messages []string        `json:"messages"`
		Errors   []string        `json:"errors"`
		Data     json.RawMessage `json:"data"`
	}

	resBody, _ := io.ReadAll(res.Body)
	if err = json.Unmarshal(resBody, &gqlRes); err != nil {
		return err
	}

	if len(gqlRes.Errors) > 0 {
		return fmt.Errorf("graphql error: %v", gqlRes.Errors)
	}

	return json.Unmarshal(gqlRes.Data, v)
}
