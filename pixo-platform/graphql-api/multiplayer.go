package graphql_api

import (
	"context"
	"encoding/json"
	"errors"
)

func (g *GraphQLAPIClient) GetMultiplayerServerConfigs(ctx context.Context, params *MultiplayerServerConfigParams) ([]*MultiplayerServerConfigQueryParams, error) {

	variables := map[string]interface{}{
		"params": params,
	}

	var query MultiplayerServerConfigQuery
	if err := g.Client.Query(ctx, &query, variables); err != nil {
		return nil, err
	}

	return query.MultiplayerServerConfigs, nil
}

func (g *GraphQLAPIClient) CreateMultiplayerServerVersion(ctx context.Context, moduleID int, image, semanticVersion, engine string) (*MultiplayerServerVersion, error) {
	query := `mutation createMultiplayerServerVersion($input: MultiplayerServerVersionInput!) { createMultiplayerServerVersion(input: $input) { id imageRegistry semanticVersion engine module { name } } }`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"moduleId":        moduleID,
			"imageRegistry":   image,
			"semanticVersion": semanticVersion,
			"engine":          engine,
			"status":          "enabled",
		},
	}

	res, err := g.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var response struct {
		ServerVersion *MultiplayerServerVersion `json:"createMultiplayerServerVersion"`
	}
	if err = json.Unmarshal(res, &response); err != nil {
		return nil, err
	}

	return response.ServerVersion, nil
}

func (g *GraphQLAPIClient) GetMultiplayerServerVersions(ctx context.Context, params *MultiplayerServerVersionQueryParams) ([]*MultiplayerServerVersion, error) {

	configs, err := g.GetMultiplayerServerConfigs(ctx, &MultiplayerServerConfigParams{
		ModuleID:      params.ModuleID,
		ServerVersion: params.SemanticVersion,
	})
	if err != nil {
		return nil, err
	}

	if len(configs) == 0 {
		return nil, errors.New("no multiplayer server configurations found")
	}

	res := make([]*MultiplayerServerVersion, len(configs[0].ServerVersions))

	for i, _ := range configs[0].ServerVersions {
		res[i] = &MultiplayerServerVersion{
			ModuleID:        configs[0].ModuleID,
			ImageRegistry:   configs[0].ServerVersions[i].ImageRegistry,
			SemanticVersion: configs[0].ServerVersions[i].SemanticVersion,
		}
	}

	return res, nil
}

func (g *GraphQLAPIClient) GetMultiplayerServerVersion(ctx context.Context, versionID int) (*MultiplayerServerVersion, error) {
	query := `query multiplayerServerVersion($id: ID!) { multiplayerServerVersion(id: $id) { id moduleId imageRegistry engine semanticVersion module { name } } }`

	variables := map[string]interface{}{
		"id": versionID,
	}

	res, err := g.Client.ExecRaw(ctx, query, variables)
	if err != nil {
		return nil, err
	}

	var response struct {
		MultiplayerServerVersion *MultiplayerServerVersion `json:"multiplayerServerVersion"`
	}
	if err = json.Unmarshal(res, &response); err != nil {
		return nil, err
	}

	return response.MultiplayerServerVersion, nil
}
