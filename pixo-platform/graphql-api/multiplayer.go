package graphql_api

import (
	"context"
	primary_api "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/go-resty/resty/v2"
)

func (g *GraphQLAPIClient) DeployMultiplayerServerVersion(moduleID int, image, semanticVersion string) (*resty.Response, error) {
	return nil, nil
}

func (g *GraphQLAPIClient) GetMultiplayerServerVersions() ([]*primary_api.MultiplayerServerVersion, error) {
	var res struct {
		MultiplayerServerVersions []*primary_api.MultiplayerServerVersion `graphql:"multiplayerServerVersions"`
	}

	if err := g.gqlClient.Query(context.Background(), &res, nil); err != nil {
		return nil, err
	}

	return res.MultiplayerServerVersions, nil
}
