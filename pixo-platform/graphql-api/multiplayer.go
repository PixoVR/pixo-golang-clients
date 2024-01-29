package graphql_api

import (
	"context"
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

func (g *GraphQLAPIClient) CreateMultiplayerServerVersion(ctx context.Context, moduleID int, image, semanticVersion string) error {
	query := `mutation createMultiplayerServerVersion($input: MultiplayerServerVersionInput!) { createMultiplayerServerVersion(input: $input) { id imageRegistry semanticVersion module { name } } }`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"moduleId":        moduleID,
			"imageRegistry":   image,
			"semanticVersion": semanticVersion,
			"status":          "enabled",
			"engine":          "unreal",
		},
	}

	if _, err := g.Client.ExecRaw(ctx, query, variables); err != nil {
		return err
	}

	return nil
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
