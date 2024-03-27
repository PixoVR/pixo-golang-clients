package graphql_api

import (
	"context"
	"encoding/json"
)

type Platform struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	ShortName string `json:"shortName,omitempty"`
}

type GetPlatformsResponse struct {
	Platforms []*Platform `json:"platforms,omitempty"`
}

func (g *GraphQLAPIClient) GetPlatforms(ctx context.Context) ([]*Platform, error) {
	query := `query platforms { platforms { id name shortName } }`

	res, err := g.Client.ExecRaw(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	var gqlRes GetPlatformsResponse
	if err = json.Unmarshal(res, &gqlRes); err != nil {
		return nil, err
	}

	return gqlRes.Platforms, nil
}
